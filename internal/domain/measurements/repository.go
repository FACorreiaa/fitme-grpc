package measurements

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	pbm "github.com/FACorreiaa/fitme-protos/modules/measurement/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type RepositoryMeasurement struct {
	pbm.UnimplementedUserMeasurementsServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewRepositoryMeasurement(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *RepositoryMeasurement {
	return &RepositoryMeasurement{pgpool: db, redis: redis, sessionManager: sessionManager}
}

// WEIGHTS

func (r *RepositoryMeasurement) CreateWeight(ctx context.Context, req *pbm.CreateWeightReq) (*pbm.XWeight, error) {
	query := `
		INSERT INTO weight_measure
		    (user_id, weight_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	currentTime := time.Now()
	var weightID string
	var updatedAt sql.NullTime
	if req.Weight.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: req.Weight.UpdatedAt.AsTime(), Valid: true}
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	err := r.pgpool.QueryRow(ctx, query, req.UserId, req.Weight.WeightValue, currentTime, updatedAt).Scan(&weightID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert weight: %v", err)
	}

	weightProto := &pbm.XWeight{
		WeightId:    weightID,
		UserId:      req.UserId,
		WeightValue: req.Weight.WeightValue,
		CreatedAt:   timestamppb.New(currentTime),
		UpdatedAt:   timestamppb.New(currentTime),
	}

	return weightProto, nil
}

func (r *RepositoryMeasurement) GetWeights(ctx context.Context) ([]*pbm.XWeight, error) {
	weightsProto := make([]*pbm.XWeight, 0)
	// test without WHERE user_id = $1
	query := `
		SELECT id, user_id, weight_value, created_at, updated_at FROM weight_measure`

	rows, err := r.pgpool.Query(ctx, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch weights: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "operation cancelled: %v", ctx.Err())
		default:
		}

		weightProto := &pbm.XWeight{}
		//var weightID, userID pgtype.UUID
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&weightProto.WeightId, &weightProto.UserId, &weightProto.WeightValue, &createdAt, &updatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch weights: %v", err)
		}

		weightProto.CreatedAt = timestamppb.New(createdAt)
		if updatedAt.Valid {
			weightProto.UpdatedAt = timestamppb.New(updatedAt.Time)
		} else {
			weightProto.UpdatedAt = nil
		}

		weightsProto = append(weightsProto, weightProto)
	}

	if len(weightsProto) == 0 {
		return nil, status.Errorf(codes.NotFound, "weight not found")
	}

	return weightsProto, nil
}

func (r *RepositoryMeasurement) GetWeight(ctx context.Context, req *pbm.GetWeightReq) (*pbm.XWeight, error) {
	weightProto := &pbm.XWeight{}
	query := `
		SELECT id, user_id, weight_value, created_at, updated_at FROM weight_measure
		WHERE id = $1 AND user_id = $2
	`
	var createdAt time.Time
	var updatedAt sql.NullTime

	err := r.pgpool.QueryRow(ctx, query, req.WeightId, req.UserId).Scan(
		&weightProto.WeightId, &weightProto.UserId, &weightProto.WeightValue, &createdAt, &updatedAt)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch weight: %v", err)
	}

	weightProto.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		weightProto.UpdatedAt = timestamppb.New(updatedAt.Time)
	} else {
		weightProto.UpdatedAt = nil
	}

	return weightProto, nil

}

func (r *RepositoryMeasurement) UpdateWeight(ctx context.Context, req *pbm.UpdateWeightReq) (*pbm.XWeight, error) {
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	weightProto := &pbm.XWeight{}

	// Define the base query
	query := `UPDATE weight_measure SET `
	var setClauses []string
	var args []interface{}
	argIndex := 1
	updatedFields := make(map[string]string)

	// Dynamically build the SET clauses based on updates
	for _, update := range req.Updates {
		switch update.Field {
		case "weight_value":
			newValue, err := strconv.ParseUint(update.NewValue, 10, 32)
			if err != nil {
				return nil, err
			}
			setClauses = append(setClauses, fmt.Sprintf("weight_value = $%d", argIndex))
			args = append(args, update.NewValue)
			updatedFields["weight_value"] = update.NewValue
			weightProto.WeightValue = int32(newValue)
			argIndex++
		case "UpdatedAt":
			setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
			args = append(args, update.NewValue)
			updatedFields["UpdatedAt"] = update.NewValue
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no updates provided")
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argIndex, argIndex+1)
	args = append(args, req.WeightId, req.UserId)

	// Execute the query
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update weight")
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	var updatedAt time.Time

	// Return a response with success
	return &pbm.XWeight{
		WeightId:    req.WeightId,
		UserId:      req.UserId,
		UpdatedAt:   timestamppb.New(updatedAt),
		WeightValue: weightProto.WeightValue,
	}, nil
}

func (r *RepositoryMeasurement) DeleteWeight(ctx context.Context, req *pbm.DeleteWeightReq) (*pbm.NilRes, error) {
	query := `
		DELETE FROM weight_measure
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.pgpool.Exec(ctx, query, req.WeightId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete exercise session: %w", err)
	}

	return &pbm.NilRes{}, nil
}

// WATER INTAKE

func (r *RepositoryMeasurement) CreateWaterMeasurement(ctx context.Context, req *pbm.CreateWaterIntakeReq) (*pbm.XWaterIntake, error) {
	query := `
		INSERT INTO water_intake
		    (user_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	currentTime := time.Now()
	var weightID string
	var updatedAt sql.NullTime
	if req.Water.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: req.Water.UpdatedAt.AsTime(), Valid: true}
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	err := r.pgpool.QueryRow(ctx, query, req.UserId, req.Water.Quantity, currentTime, updatedAt).Scan(&weightID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert weight: %v", err)
	}

	waterProto := &pbm.XWaterIntake{
		WaterIntakeId: weightID,
		UserId:        req.UserId,
		Quantity:      req.Water.Quantity,
		CreatedAt:     timestamppb.New(currentTime),
		UpdatedAt:     timestamppb.New(currentTime),
	}

	return waterProto, nil
}

func (r *RepositoryMeasurement) GetWaterMeasurements(ctx context.Context) ([]*pbm.XWaterIntake, error) {
	waterProtos := make([]*pbm.XWaterIntake, 0)
	// test without WHERE user_id = $1
	query := `
		SELECT id, user_id, quantity, created_at, updated_at FROM water_intake`

	rows, err := r.pgpool.Query(ctx, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch weights: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "operation cancelled: %v", ctx.Err())
		default:
		}

		waterProto := &pbm.XWaterIntake{}
		//var weightID, userID pgtype.UUID
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&waterProto.WaterIntakeId, &waterProto.UserId, &waterProto.WaterIntakeId, &createdAt, &updatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch weights: %v", err)
		}

		waterProto.CreatedAt = timestamppb.New(createdAt)
		if updatedAt.Valid {
			waterProto.UpdatedAt = timestamppb.New(updatedAt.Time)
		} else {
			waterProto.UpdatedAt = nil
		}

		waterProtos = append(waterProtos, waterProto)
	}

	if len(waterProtos) == 0 {
		return nil, status.Errorf(codes.NotFound, "weight not found")
	}

	return waterProtos, nil
}

func (r *RepositoryMeasurement) GetWaterMeasurement(ctx context.Context, req *pbm.GetWaterIntakeReq) (*pbm.XWaterIntake, error) {
	waterProto := &pbm.XWaterIntake{}
	query := `
		SELECT id, user_id, quantity, created_at, updated_at FROM water_intake
		WHERE id = $1 AND user_id = $2
	`
	var createdAt time.Time
	var updatedAt sql.NullTime

	err := r.pgpool.QueryRow(ctx, query, req.WaterIntakeId, req.UserId).Scan(
		&waterProto.WaterIntakeId, &waterProto.UserId, &waterProto.Quantity, &createdAt, &updatedAt)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch weight: %v", err)
	}

	waterProto.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		waterProto.UpdatedAt = timestamppb.New(updatedAt.Time)
	} else {
		waterProto.UpdatedAt = nil
	}

	return waterProto, nil
}

func (r *RepositoryMeasurement) UpdateWaterMeasurement(ctx context.Context, req *pbm.UpdateWaterIntakeReq) (*pbm.XWaterIntake, error) {
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	waterProto := &pbm.XWaterIntake{}

	// Define the base query
	query := `UPDATE water_intake SET `
	var setClauses []string
	var args []interface{}
	argIndex := 1
	updatedFields := make(map[string]string)

	// Dynamically build the SET clauses based on updates
	for _, update := range req.Updates {
		switch update.Field {
		case "quantity":
			newValue, err := strconv.ParseUint(update.NewValue, 10, 32)
			if err != nil {
				return nil, err
			}
			setClauses = append(setClauses, fmt.Sprintf("weight_value = $%d", argIndex))
			args = append(args, update.NewValue)
			updatedFields["weight_value"] = update.NewValue
			waterProto.Quantity = int32(newValue)
			argIndex++
		case "UpdatedAt":
			setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
			args = append(args, update.NewValue)
			updatedFields["UpdatedAt"] = update.NewValue
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no updates provided")
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argIndex, argIndex+1)
	args = append(args, req.WaterIntakeId, req.UserId)

	// Execute the query
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update weight")
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	var updatedAt time.Time

	// Return a response with success
	return &pbm.XWaterIntake{
		WaterIntakeId: req.WaterIntakeId,
		UserId:        req.UserId,
		UpdatedAt:     timestamppb.New(updatedAt),
		Quantity:      waterProto.Quantity,
	}, nil
}

func (r *RepositoryMeasurement) DeleteWaterMeasurement(ctx context.Context, req *pbm.DeleteWaterIntakeReq) (*pbm.NilRes, error) {
	query := `
		DELETE FROM water_intake
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.pgpool.Exec(ctx, query, req.WaterIntakeId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete exercise session: %w", err)
	}

	return &pbm.NilRes{}, nil
}

// WASTE LINE

func (r *RepositoryMeasurement) CreateWasteLineMeasurement(ctx context.Context, req *pbm.CreateWasteLineReq) (*pbm.XWasteLine, error) {
	query := `
		INSERT INTO waist_line
		    (user_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	currentTime := time.Now()
	var weightID string
	var updatedAt sql.NullTime
	if req.WasteLine.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: req.WasteLine.UpdatedAt.AsTime(), Valid: true}
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	err := r.pgpool.QueryRow(ctx, query, req.UserId, req.WasteLine.Measurement, currentTime, updatedAt).Scan(&weightID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert weight: %v", err)
	}

	waterProto := &pbm.XWasteLine{
		WasteLineId: weightID,
		UserId:      req.UserId,
		Measurement: req.WasteLine.Measurement,
		CreatedAt:   timestamppb.New(currentTime),
		UpdatedAt:   timestamppb.New(currentTime),
	}

	return waterProto, nil
}

func (r *RepositoryMeasurement) GetWasteLineMeasurements(ctx context.Context) ([]*pbm.XWasteLine, error) {
	wasteLineProtos := make([]*pbm.XWasteLine, 0)
	// test without WHERE user_id = $1
	query := `
		SELECT id, user_id, quantity, created_at, updated_at FROM waist_line`

	rows, err := r.pgpool.Query(ctx, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch weights: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "operation cancelled: %v", ctx.Err())
		default:
		}

		wastelineProto := &pbm.XWasteLine{}
		//var weightID, userID pgtype.UUID
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&wastelineProto.WasteLineId, &wastelineProto.UserId, &wastelineProto.Measurement, &createdAt, &updatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch weights: %v", err)
		}

		wastelineProto.CreatedAt = timestamppb.New(createdAt)
		if updatedAt.Valid {
			wastelineProto.UpdatedAt = timestamppb.New(updatedAt.Time)
		} else {
			wastelineProto.UpdatedAt = nil
		}

		wasteLineProtos = append(wasteLineProtos, wastelineProto)
	}

	if len(wasteLineProtos) == 0 {
		return nil, status.Errorf(codes.NotFound, "weight not found")
	}

	return wasteLineProtos, nil
}

func (r *RepositoryMeasurement) GetWasteLineMeasurement(ctx context.Context, req *pbm.GetWasteLineReq) (*pbm.XWasteLine, error) {
	waistlineProto := &pbm.XWasteLine{}
	query := `
		SELECT id, user_id, quantity, created_at, updated_at FROM waist_line
		WHERE id = $1 AND user_id = $2
	`
	var createdAt time.Time
	var updatedAt sql.NullTime

	err := r.pgpool.QueryRow(ctx, query, req.WasteLineId, req.UserId).Scan(
		&waistlineProto.WasteLineId, &waistlineProto.UserId, &waistlineProto.Measurement, &createdAt, &updatedAt)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch weight: %v", err)
	}

	waistlineProto.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		waistlineProto.UpdatedAt = timestamppb.New(updatedAt.Time)
	} else {
		waistlineProto.UpdatedAt = nil
	}

	return waistlineProto, nil
}

func (r *RepositoryMeasurement) UpdateWasteLineMeasurement(ctx context.Context, req *pbm.UpdateWasteLineReq) (*pbm.XWasteLine, error) {
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	waistlineProto := &pbm.XWasteLine{}

	query := `UPDATE waist_line SET `
	var setClauses []string
	var args []interface{}
	argIndex := 1
	updatedFields := make(map[string]string)

	for _, update := range req.Updates {
		switch update.Field {
		case "measurement":
			newValue, err := strconv.ParseUint(update.NewValue, 10, 32)
			if err != nil {
				return nil, err
			}
			setClauses = append(setClauses, fmt.Sprintf("measurement = $%d", argIndex))
			args = append(args, update.NewValue)
			updatedFields["measurement"] = update.NewValue
			waistlineProto.Measurement = int32(newValue)
			argIndex++
		case "UpdatedAt":
			setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
			args = append(args, update.NewValue)
			updatedFields["UpdatedAt"] = update.NewValue
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no updates provided")
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argIndex, argIndex+1)
	args = append(args, req.WasteLineId, req.UserId)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update weight")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	var updatedAt time.Time

	return &pbm.XWasteLine{
		WasteLineId: req.WasteLineId,
		UserId:      req.UserId,
		UpdatedAt:   timestamppb.New(updatedAt),
		Measurement: waistlineProto.Measurement,
	}, nil
}

func (r *RepositoryMeasurement) DeleteWasteLineMeasurement(ctx context.Context, req *pbm.DeleteWasteLineReq) (*pbm.NilRes, error) {
	query := `
		DELETE FROM waist_line
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.pgpool.Exec(ctx, query, req.WasteLineId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete exercise session: %w", err)
	}

	return &pbm.NilRes{}, nil
}
