package workout

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	pbw "github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type RepositoryWorkout struct {
	pbw.UnimplementedWorkoutServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewRepositoryWorkout(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *RepositoryWorkout {
	return &RepositoryWorkout{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func nullTimeToTimestamppb(nt sql.NullTime) *timestamppb.Timestamp {
	if nt.Valid {
		return timestamppb.New(nt.Time)
	}
	return nil
}

func (r *RepositoryWorkout) GetExercises(ctx context.Context, req *pbw.GetExercisesReq) (*pbw.GetExercisesRes, error) {
	exercisesProtoList := make([]*pbw.XExercises, 0)
	query := `SELECT DISTINCT
    			id, name, type, muscle, equipment, difficulty,
				instructions, video, custom_created, created_at, updated_at
				FROM exercise_list`

	rows, err := r.pgpool.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no exercises found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch workout exercises: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		exerciseProto := pbw.XExercises{}
		e := &Exercises{}

		err := rows.Scan(
			&e.ID, &e.Name, &e.ExerciseType, &e.MuscleGroup, &e.Equipment, &e.Difficulty,
			&e.Instructions, &e.Video, &e.CustomCreated, &e.CreatedAt, &e.UpdatedAt,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("no exercises found: %w", err)
			}
			return nil, status.Error(codes.Internal, "Internal server error")
		}

		createdAt := timestamppb.New(e.CreatedAt)
		var updatedAt sql.NullTime
		if e.UpdatedAt.Valid {
			updatedAt = e.UpdatedAt
		} else {
			updatedAt = sql.NullTime{Valid: false}
		}

		exerciseProto.ExerciseId = e.ID
		exerciseProto.Name = e.Name
		exerciseProto.ExerciseType = e.ExerciseType
		exerciseProto.MuscleGroup = e.MuscleGroup
		exerciseProto.Equipment = e.Equipment
		exerciseProto.Difficulty = e.Difficulty
		exerciseProto.Instruction = e.Instructions
		exerciseProto.Video = e.Video
		exerciseProto.CustomCreated = e.CustomCreated
		exerciseProto.CreatedAt = createdAt
		exerciseProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)

		newProtoExerciseList := &pbw.XExercises{
			ExerciseId:    exerciseProto.ExerciseId,
			Name:          exerciseProto.Name,
			ExerciseType:  exerciseProto.ExerciseType,
			MuscleGroup:   exerciseProto.MuscleGroup,
			Equipment:     exerciseProto.Equipment,
			Difficulty:    exerciseProto.Difficulty,
			Instruction:   exerciseProto.Instruction,
			Video:         exerciseProto.Video,
			CustomCreated: exerciseProto.CustomCreated,
			CreatedAt:     exerciseProto.CreatedAt,
			UpdatedAt:     exerciseProto.UpdatedAt,
		}

		exercisesProtoList = append(exercisesProtoList, newProtoExerciseList)
	}

	if len(exercisesProtoList) == 0 {
		return nil, fmt.Errorf("no exercises found")
	}

	return &pbw.GetExercisesRes{
		Success:  true,
		Message:  "Exercises retrieved successfully",
		Exercise: exercisesProtoList,
		Response: &pbw.BaseResponse{
			Upstream: "workout-service",
		},
	}, nil

}

func (r *RepositoryWorkout) GetExerciseID(ctx context.Context, req *pbw.GetExerciseIDReq) (*pbw.GetExerciseIDRes, error) {
	exerciseProto := &pbw.XExercises{}
	exercise := &Exercises{}
	id := req.ExerciseId

	if id == "" {
		return &pbw.GetExerciseIDRes{}, status.Error(codes.InvalidArgument, "workout ID is required")
	}

	query := `SELECT 	id, name, type, muscle, equipment, difficulty,
						instructions, video, custom_created, created_at, updated_at
			   FROM exercise_list
			   WHERE id = $1`

	err := r.pgpool.QueryRow(ctx, query, id).Scan(
		&exercise.ID, &exercise.Name, &exercise.ExerciseType, &exercise.MuscleGroup, &exercise.Equipment,
		&exercise.Difficulty, &exercise.Instructions, &exercise.Video, &exercise.CustomCreated, &exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no exercises found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch workout exercises: %w", err)
	}

	createdAt := timestamppb.New(exercise.CreatedAt)
	var updatedAt sql.NullTime

	if exercise.UpdatedAt.Valid {
		updatedAt = exercise.UpdatedAt
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	exerciseProto.ExerciseId = exercise.ID
	exerciseProto.Name = exercise.Name
	exerciseProto.MuscleGroup = exercise.MuscleGroup
	exerciseProto.Equipment = exercise.Equipment
	exerciseProto.Difficulty = exercise.Difficulty
	exerciseProto.Instruction = exercise.Instructions
	exerciseProto.Video = exercise.Video
	exerciseProto.CustomCreated = exercise.CustomCreated
	exerciseProto.CreatedAt = createdAt
	exerciseProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)

	return &pbw.GetExerciseIDRes{
		Success:  true,
		Message:  "Exercise retrieved successfully",
		Exercise: exerciseProto,
		Response: &pbw.BaseResponse{
			Upstream: "workout-service",
		},
	}, nil
}

func (r *RepositoryWorkout) CreateExercise(ctx context.Context, req *pbw.CreateExerciseReq) (*pbw.CreateExerciseRes, error) {
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	query := `
				INSERT INTO exercise_list (name, type, muscle, equipment, difficulty,
                                   instructions, video,
                                   created_at, updated_at)
        		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                RETURNING id`

	currentTime := time.Now()

	var exerciseID string
	err = tx.QueryRow(ctx, query,
		req.Exercise.Name,
		req.Exercise.ExerciseType,
		req.Exercise.MuscleGroup,
		req.Exercise.Equipment,
		req.Exercise.Difficulty,
		req.Exercise.Instruction,
		req.Exercise.Video,
		currentTime,
		currentTime,
	).Scan(&exerciseID)

	setExerciseToUserQuery := `
				INSERT INTO user_exercises (user_id, exercise_id)
				VALUES ($1, $2)
				RETURNING user_id, exercise_id`

	req.Exercise.ExerciseId = exerciseID

	var userID, associatedExerciseID string

	err = tx.QueryRow(ctx, setExerciseToUserQuery, req.UserId, req.Exercise.ExerciseId).Scan(&userID, &associatedExerciseID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to associate exercise with user: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	exerciseProto := &pbw.XExercises{
		ExerciseId:    exerciseID,
		Name:          req.Exercise.Name,
		ExerciseType:  req.Exercise.ExerciseType,
		MuscleGroup:   req.Exercise.MuscleGroup,
		Equipment:     req.Exercise.Equipment,
		Difficulty:    req.Exercise.Difficulty,
		Instruction:   req.Exercise.Instruction,
		Video:         req.Exercise.Video,
		CustomCreated: true,
		CreatedAt:     timestamppb.New(currentTime), // Proto timestamp for created_at fix later
		UpdatedAt:     timestamppb.New(currentTime), // Proto timestamp for updated_at fix later
	}

	return &pbw.CreateExerciseRes{
		Success:  true,
		Message:  "Exercise created and associated with user successfully",
		Exercise: exerciseProto,
		Response: &pbw.BaseResponse{
			Upstream: "workout-service",
		},
	}, nil
}

func (r *RepositoryWorkout) DeleteExercise(ctx context.Context, req *pbw.DeleteExerciseReq) (*pbw.NilRes, error) {
	query := `DELETE FROM exercise_list WHERE id = $1
 			  AND exercise_list.custom_created = true`
	_, err := r.pgpool.Exec(ctx, query, req.ExerciseId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete exercise: %w", err)
	}
	return &pbw.NilRes{}, nil
}

func (r *RepositoryWorkout) UpdateExercise(ctx context.Context, req *pbw.UpdateExerciseReq) (*pbw.UpdateExerciseRes, error) {
	query := `UPDATE exercise_list SET `
	var setClauses []string
	var args []interface{}
	argIndex := 1
	updatedExercise := &pbw.XExercises{}

	for _, update := range req.Updates {
		switch update.Field {
		case "name":
			setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedExercise.Name = update.NewValue
		case "muscle_group":
			setClauses = append(setClauses, fmt.Sprintf("muscle_group = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedExercise.MuscleGroup = update.NewValue
		case "equipment":
			setClauses = append(setClauses, fmt.Sprintf("equipment = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedExercise.Equipment = update.NewValue
		case "difficulty":
			setClauses = append(setClauses, fmt.Sprintf("difficulty = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedExercise.Difficulty = update.NewValue
		case "instruction":
			setClauses = append(setClauses, fmt.Sprintf("instruction = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedExercise.Instruction = update.NewValue
		case "video":
			setClauses = append(setClauses, fmt.Sprintf("video = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedExercise.Video = update.NewValue
		default:
			return nil, fmt.Errorf("unsupported update field: %s", update.Field)
		}
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no updates provided")
	}

	query += strings.Join(setClauses, ", ")
	query += ` WHERE id = $` + fmt.Sprintf("%d", argIndex)
	args = append(args, req.ExerciseId)

	_, err := r.pgpool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update exercise: %w", err)
	}

	updatedExercise.ExerciseId = req.ExerciseId
	getQuery := `SELECT id, name, muscle, equipment, difficulty, instructions, video FROM exercise_list WHERE id = $1`
	err = r.pgpool.QueryRow(ctx, getQuery, req.ExerciseId).Scan(
		&updatedExercise.ExerciseId,
		&updatedExercise.Name,
		&updatedExercise.MuscleGroup,
		&updatedExercise.Equipment,
		&updatedExercise.Difficulty,
		&updatedExercise.Instruction,
		&updatedExercise.Video,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated exercise: %w", err)
	}

	return &pbw.UpdateExerciseRes{
		Success:  true,
		Message:  "Exercise updated successfully",
		Exercise: updatedExercise,
	}, nil
}

//func (r *RepositoryWorkout) CreateWorkoutPlan(ctx context.Context, req *pbw.InsertWorkoutPlanReq) (*pbw.InsertWorkoutPlanRes, error) {
//	logger := zap.L() // Assuming you have initialized your logger
//	newWorkoutPlan := req.Workout
//	plan := req.PlanDay
//
//	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
//	if err != nil {
//		logger.Error("Failed to start transaction", zap.Error(err))
//		return nil, status.Error(codes.Internal, "failed to start transaction")
//	}
//	defer func() {
//		if err != nil {
//			_ = tx.Rollback(ctx)
//			logger.Warn("Transaction rolled back", zap.Error(err))
//		}
//	}()
//
//	// Insert workout plan
//	query := `
//        INSERT INTO workout_plan (id, user_id, description, notes, rating, created_at)
//        VALUES ($1, $2, $3, $4, $5, $6)
//    `
//	_, err = tx.Exec(ctx, query, newWorkoutPlan.WorkoutId, newWorkoutPlan.UserId, newWorkoutPlan.Description, newWorkoutPlan.Notes, newWorkoutPlan.Rating, time.Now())
//	if err != nil {
//		logger.Error("Failed to insert workout plan", zap.Error(err))
//		return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to insert workout plan")
//	}
//
//	// Insert workout days and details in batches
//	var workoutDayValues []interface{}
//	var workoutPlanDetailValues []interface{}
//
//	for _, day := range plan {
//		createdAt := timestamppb.New(time.Now())
//		workDayID := uuid.NewString()
//
//		workoutDayValues = append(workoutDayValues, workDayID, newWorkoutPlan.WorkoutId, day.Day, createdAt.AsTime())
//
//		workoutPlanDetailID := uuid.NewString()
//		workoutPlanDetailValues = append(workoutPlanDetailValues, workoutPlanDetailID, newWorkoutPlan.WorkoutId, day.Day, day.ExerciseId, createdAt.AsTime())
//	}
//
//	// Batch insert workout days
//	workoutDayQuery := `
//        INSERT INTO workout_day (id, workout_plan_id, day, created_at)
//        VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ...
//    `
//	_, err = tx.Exec(ctx, workoutDayQuery, workoutDayValues...)
//	if err != nil {
//		logger.Error("Failed to insert workout days", zap.Error(err))
//		return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to insert workout days")
//	}
//
//	// Batch insert workout plan details
//	workoutPlanDetailQuery := `
//        INSERT INTO workout_plan_detail (id, workout_plan_id, day, exercises, created_at)
//        VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10), ...
//    `
//	_, err = tx.Exec(ctx, workoutPlanDetailQuery, workoutPlanDetailValues...)
//	if err != nil {
//		logger.Error("Failed to insert workout plan details", zap.Error(err))
//		return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to insert workout plan details")
//	}
//
//	err = tx.Commit(ctx)
//	if err != nil {
//		logger.Error("Failed to commit transaction", zap.Error(err))
//		return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to commit transaction")
//	}
//
//	logger.Info("Workout plan created successfully", zap.String("workout_id", newWorkoutPlan.WorkoutId))
//
//	return &pbw.InsertWorkoutPlanRes{
//		Success: true,
//		Message: "Workout plan created successfully",
//		Workout: &insertedPlan,
//	}, nil
//}

func (r *RepositoryWorkout) CreateWorkoutPlan(ctx context.Context, req *pbw.InsertWorkoutPlanReq) (*pbw.InsertWorkoutPlanRes, error) {
	newWorkoutPlan := req.Workout
	plan := req.PlanDay

	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error("Failed to start transaction", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			logger.Warn("Transaction rolled back", zap.Error(err))
		}
	}()

	var workoutPlanID string
	insertWorkoutPlanQuery := `
       INSERT INTO workout_plan (user_id, description, notes, rating, created_at)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING id`

	err = tx.QueryRow(ctx, insertWorkoutPlanQuery, newWorkoutPlan.UserId,
		newWorkoutPlan.Description, newWorkoutPlan.Notes, newWorkoutPlan.Rating, time.Now()).Scan(&workoutPlanID)

	if err != nil {
		logger.Error("Failed to insert workout plan", zap.Error(err))
		return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to insert workout plan")
	}

	row := tx.QueryRow(ctx, "SELECT id, user_id, description, notes, rating, created_at, updated_at FROM workout_plan WHERE id = $1", workoutPlanID)
	insertedPlan := &pbw.XWorkoutPlan{}
	createdAt := time.Time{}
	updatedAt := sql.NullTime{}

	err = row.Scan(&insertedPlan.WorkoutId, &insertedPlan.UserId, &insertedPlan.Description, &insertedPlan.Notes, &insertedPlan.Rating,
		&createdAt, &updatedAt)
	if err != nil {
		logger.Error("no rows found on workout", zap.Error(err))
		return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to fetch inserted workout plan")
	}
	insertedPlan.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		insertedPlan.UpdatedAt = timestamppb.New(updatedAt.Time)
	} else {
		insertedPlan.UpdatedAt = nil
	}

	// Insert each workout day
	for _, day := range plan {
		createdAt = time.Now()
		workDayID := uuid.NewString()

		workoutDayQuery := `
       INSERT INTO workout_day (id, workout_plan_id, day, created_at)
       VALUES ($1, $2, $3, $4)
   `
		_, err := tx.Exec(ctx, workoutDayQuery, workDayID, workoutPlanID, day.Day, createdAt)

		if err != nil {
			logger.Error("Failed to insert workout plan details", zap.Error(err))
			return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to insert workout day")
		}

		workoutPlanDetail := &pbw.XWorkoutPlanDetail{
			WorkoutPlanId:       uuid.NewString(),
			WorkoutPlanDetailId: insertedPlan.WorkoutId,
			Day:                 day.Day,
			Exercises:           day.ExerciseId,
			CreatedAt:           timestamppb.New(createdAt),
		}

		// Insert workout plan details
		workoutPlanDetailQuery := `
       INSERT INTO workout_plan_detail (id, workout_plan_id, day, exercises, created_at)
       VALUES ($1, $2, $3, $4, $5)
   `
		_, err = tx.Exec(ctx, workoutPlanDetailQuery, workoutPlanDetail.WorkoutPlanDetailId,
			workoutPlanID,
			workoutPlanDetail.Day,
			workoutPlanDetail.Exercises,
			createdAt)

		if err != nil {
			logger.Error("Failed to insert workout plan detail", zap.Error(err))
			return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to insert workout plan detail")
		}
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		logger.Error("Failed to commit transaction", zap.Error(err))
		return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "failed to commit transaction")
	}

	//logger.Info("Workout plan created successfully", zap.String("workout_id", newWorkoutPlan.WorkoutId))

	// Ensure insertedPlan is not nil before returning
	//if insertedPlan == nil {
	//	logger.Error("Inserted plan is nil after successful commit")
	//	return &pbw.InsertWorkoutPlanRes{}, status.Error(codes.Internal, "inserted workout plan is nil")
	//}

	return &pbw.InsertWorkoutPlanRes{
		Success: true,
		Message: "Workout plan created successfully",
		Workout: insertedPlan,
	}, nil
}

func (r *RepositoryWorkout) GetWorkoutPlans(ctx context.Context, req *pbw.GetWorkoutPlansReq) (*pbw.GetWorkoutPlansRes, error) {
	query := `
		SELECT wp.id AS workout_plan_id, wp.user_id, wp.description,
			   wp.notes, wp.rating, wp.created_at, wp.updated_at, wd.day, wpd.exercises
		FROM workout_plan AS wp
		LEFT JOIN workout_plan_detail AS wpd ON wp.id = wpd.workout_plan_id
		LEFT JOIN workout_day AS wd ON wp.id = wd.workout_plan_id
		GROUP BY wp.id, wd.day, wpd.exercises
		ORDER BY wd.day;
	`

	rows, err := r.pgpool.Query(ctx, query)
	if err != nil {
		return &pbw.GetWorkoutPlansRes{}, status.Error(codes.Internal, "failed to query workout plans")
	}
	defer rows.Close()

	workouts := make(map[string]*pbw.XWorkoutPlanResponse)

	for rows.Next() {

		var row struct {
			WorkoutPlanID uuid.UUID    `db:"workout_plan_id"`
			UserID        uuid.UUID    `db:"user_id"`
			Description   string       `db:"description"`
			Notes         string       `db:"notes"`
			Rating        int          `db:"rating"`
			CreatedAt     time.Time    `db:"created_at"`
			UpdatedAt     sql.NullTime `db:"updated_at"`
			Day           string       `db:"day"`
			Exercises     []string     `db:"exercises"`
		}

		createdAt := timestamppb.New(row.CreatedAt)
		var updatedAt *timestamppb.Timestamp

		if row.UpdatedAt.Valid {
			updatedAt = timestamppb.New(row.UpdatedAt.Time)
		} else {
			updatedAt = nil
		}

		err = rows.Scan(&row.WorkoutPlanID, &row.UserID, &row.Description, &row.Notes, &row.Rating, &row.CreatedAt,
			&row.UpdatedAt, &row.Day, &row.Exercises)
		if err != nil {
			logger.Error("Failed to scan workout plan", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to scan workout plan")
		}

		if _, ok := workouts[row.WorkoutPlanID.String()]; !ok {
			workouts[row.WorkoutPlanID.String()] = &pbw.XWorkoutPlanResponse{
				WorkoutPlanId: row.WorkoutPlanID.String(),
				UserId:        row.UserID.String(),
				Description:   row.Description,
				WorkoutDay:    []*pbw.WorkoutDayResponse{},
				Notes:         row.Notes,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
				Rating:        uint32(row.Rating),
			}
		}

		if plan, ok := workouts[row.WorkoutPlanID.String()]; ok {
			day := &pbw.WorkoutDayResponse{
				Day:       row.Day,
				Exercises: row.Exercises,
			}
			plan.WorkoutDay = append(plan.WorkoutDay, day)
			workouts[row.WorkoutPlanID.String()] = plan
		}
	}

	result := &pbw.GetWorkoutPlansRes{
		Success:     true,
		Message:     "Workout plans retrieved successfully",
		WorkoutPlan: []*pbw.XWorkoutPlanResponse{},
	}
	for _, workout := range workouts {
		result.WorkoutPlan = append(result.WorkoutPlan, workout)
	}

	return result, nil
}

func (r *RepositoryWorkout) GetWorkoutPlan(ctx context.Context, req *pbw.GetWorkoutPlanReq) (*pbw.GetWorkoutPlanRes, error) {
	workoutPlanProto := &pbw.GetWorkoutPlanRes{
		WorkoutPlan: &pbw.XWorkoutPlanResponse{},
	}
	workoutPlanID := req.WorkoutPlanId

	if workoutPlanID == "" {
		return &pbw.GetWorkoutPlanRes{}, status.Error(codes.InvalidArgument, "workout ID is required")
	}

	var row struct {
		WorkoutPlanID uuid.UUID      `db:"workout_plan_id"`
		UserID        uuid.UUID      `db:"user_id"`
		Description   string         `db:"description"`
		Notes         string         `db:"notes"`
		Rating        int            `db:"rating"`
		CreatedAt     time.Time      `db:"created_at"`
		UpdatedAt     sql.NullTime   `db:"updated_at"`
		Day           string         `db:"day"`
		Exercises     pq.StringArray `db:"exercises"`
	}

	query := `SELECT
				wp.id AS workout_plan_id, wp.user_id, wp.description,
				wp.notes, wp.rating, wp.created_at, wp.updated_at, wd.day, wpd.exercises
				FROM workout_plan AS wp
				JOIN workout_plan_detail AS wpd ON wp.id = wpd.workout_plan_id
				JOIN workout_day AS wd ON wp.id = wd.workout_plan_id
				WHERE wp.id = $1
				GROUP BY wp.id, wd.day, wpd.exercises
				ORDER BY wd.day;`

	// Perform the query and scan the results
	err := r.pgpool.QueryRow(ctx, query, workoutPlanID).Scan(
		&row.WorkoutPlanID,
		&row.UserID,
		&row.Description,
		&row.Notes,
		&row.Rating,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.Day,
		&row.Exercises,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to query workout plan")
	}

	createdAt := timestamppb.New(row.CreatedAt)
	updatedAt := nullTimeToTimestamppb(row.UpdatedAt)

	workoutPlanProto.WorkoutPlan.WorkoutPlanId = row.WorkoutPlanID.String()
	workoutPlanProto.WorkoutPlan.UserId = row.UserID.String()
	workoutPlanProto.WorkoutPlan.Description = row.Description
	workoutPlanProto.WorkoutPlan.Notes = row.Notes
	workoutPlanProto.WorkoutPlan.Rating = uint32(row.Rating)
	workoutPlanProto.WorkoutPlan.CreatedAt = createdAt
	workoutPlanProto.WorkoutPlan.UpdatedAt = updatedAt
	workoutPlanProto.WorkoutPlan.Day = row.Day
	workoutPlanProto.WorkoutPlan.Exercises = row.Exercises

	return workoutPlanProto, nil
}

func (r *RepositoryWorkout) DeleteWorkoutPlan(ctx context.Context, req *pbw.DeleteWorkoutPlanReq) (*pbw.NilRes, error) {
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	workoutPlanID := req.WorkoutPlanId
	userID := req.UserId

	// TODO refactor
	// Delete from workout_plan
	result, err := tx.Exec(ctx, `
        DELETE FROM workout_day
	   	WHERE workout_plan_id = $1`,
		workoutPlanID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "workout plan not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete workout plan")
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM workout_plan_detail
	   	WHERE workout_plan_id = $1`,
		workoutPlanID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete workout plan detail")
	}

	// Delete from workout_plan_detail
	_, err = tx.Exec(ctx, `
		DELETE FROM workout_plan
		WHERE id = $1 AND user_id = $2`,
		workoutPlanID, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete workout plan detail")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "no rows deleted")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	return &pbw.NilRes{}, nil
}

// GetWorkoutPlanExercises verify later
func (r *RepositoryWorkout) GetWorkoutPlanExercises(ctx context.Context, req *pbw.GetWorkoutPlanExercisesReq) (*pbw.GetWorkoutPlanExercisesRes, error) {
	workoutProtoList := make([]*pbw.XWorkoutExerciseDay, 0)
	query := `SELECT el.id, el.name, el.type, el.muscle, el.equipment, el.difficulty, el.instructions,
       				el.video, el.custom_created, el.created_at, el.updated_at, wpd.day
					FROM workout_plan_detail wpd
					JOIN exercise_list el ON el.id = ANY(wpd.exercises)`
	rows, err := r.pgpool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workout exercises: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		workoutProto := pbw.XWorkoutExerciseDay{}
		workoutList := &WorkoutExerciseDay{}

		err = rows.Scan(
			&workoutList.ID, &workoutList.Name, &workoutList.ExerciseType, &workoutList.MuscleGroup, &workoutList.Equipment,
			&workoutList.Difficulty, &workoutList.Instructions, &workoutList.Video, &workoutList.CustomCreated,
			&workoutList.CreatedAt, &workoutList.UpdatedAt, &workoutList.Day)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("no exercises found: %w", err)
			}
			return nil, fmt.Errorf("failed to fetch workout exercises: %w", err)
		}

		createdAt := timestamppb.New(workoutList.CreatedAt)
		var updatedAt sql.NullTime
		if workoutList.UpdatedAt.Valid {
			updatedAt = workoutList.UpdatedAt
		} else {
			updatedAt = sql.NullTime{Valid: false}
		}
		workoutProto.WorkoutExerciseDay = workoutList.Day
		workoutProto.Name = workoutList.Name
		workoutProto.ExerciseType = workoutList.ExerciseType
		workoutProto.MuscleGroup = workoutList.MuscleGroup
		workoutProto.Equipment = workoutList.Equipment
		workoutProto.Difficulty = workoutList.Difficulty
		workoutProto.Instruction = workoutList.Instructions
		workoutProto.Video = workoutList.Video
		workoutProto.CustomCreated = workoutList.CustomCreated
		workoutProto.CreatedAt = createdAt
		workoutProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)
		workoutProto.Day = workoutList.Day

		newProtoWorkoutList := &pbw.XWorkoutExerciseDay{
			WorkoutExerciseDay: workoutProto.WorkoutExerciseDay,
			Name:               workoutProto.Name,
			ExerciseType:       workoutProto.ExerciseType,
			MuscleGroup:        workoutProto.MuscleGroup,
			Equipment:          workoutProto.Equipment,
			Difficulty:         workoutProto.Difficulty,
			Instruction:        workoutProto.Instruction,
			Video:              workoutProto.Video,
			CustomCreated:      workoutProto.CustomCreated,
			CreatedAt:          workoutProto.CreatedAt,
			UpdatedAt:          workoutProto.UpdatedAt,
			Day:                workoutProto.Day,
		}

		workoutProtoList = append(workoutProtoList, newProtoWorkoutList)
	}
	if len(workoutProtoList) == 0 {
		return nil, fmt.Errorf("no exercises found: %w", err)
	}

	return &pbw.GetWorkoutPlanExercisesRes{
		Success:            true,
		Message:            "workouts retrieved successfully",
		WorkoutExerciseDay: workoutProtoList,
		Response: &pbw.BaseResponse{
			Upstream: "workout-service",
		},
	}, nil
}

// GetWorkoutPlanExercisesByID verify later
func (r *RepositoryWorkout) GetWorkoutPlanExercisesByID(ctx context.Context, req *pbw.GetExerciseByIdWorkoutPlanReq) (*pbw.GetExerciseByIdWorkoutPlanRes, error) {
	workoutProto := &pbw.XWorkoutExerciseDay{}
	workout := &WorkoutExerciseDay{}
	exerciseID := req.ExerciseWorkoutPlan

	if exerciseID == "" {
		return &pbw.GetExerciseByIdWorkoutPlanRes{}, status.Error(codes.InvalidArgument, "workout ID is required")
	}

	query := `SELECT el.id, el.name, el.type, el.muscle, el.equipment, el.difficulty, el.instructions,
       				el.video, el.custom_created, el.created_at, el.updated_at, wpd.day
					FROM workout_plan_detail wpd
					JOIN exercise_list el ON el.id = ANY(wpd.exercises)
					WHERE wpd.workout_plan_id = $1`
	err := r.pgpool.QueryRow(ctx, query, &exerciseID).Scan(&exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workout exercises: %w", err)
	}

	createdAt := timestamppb.New(workout.CreatedAt)
	var updatedAt sql.NullTime
	if workout.UpdatedAt.Valid {
		updatedAt = workout.UpdatedAt
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}
	workoutProto.WorkoutExerciseDay = workout.Day
	workoutProto.Name = workout.Name
	workoutProto.ExerciseType = workout.ExerciseType
	workoutProto.MuscleGroup = workout.MuscleGroup
	workoutProto.Equipment = workout.Equipment
	workoutProto.Difficulty = workout.Difficulty
	workoutProto.Instruction = workout.Instructions
	workoutProto.Video = workout.Video
	workoutProto.CustomCreated = workout.CustomCreated
	workoutProto.CreatedAt = createdAt
	workoutProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)
	workoutProto.Day = workout.Day

	return &pbw.GetExerciseByIdWorkoutPlanRes{
		Success:            true,
		Message:            "workout successfully",
		WorkoutExerciseDay: workoutProto,
		Response: &pbw.BaseResponse{
			Upstream: "workout-service",
		},
	}, nil
}

func (r *RepositoryWorkout) InsertExerciseWorkoutPlan(ctx context.Context, req *pbw.InsertExerciseWorkoutPlanReq) (*pbw.NilRes, error) {
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	exerciseID := req.ExerciseId
	workoutPlanID := req.WorkoutPlanId
	workoutDay := req.WorkoutDay

	query := `
		UPDATE workout_plan_detail
		SET exercises = array_append(exercises, $1)
		WHERE workout_plan_id = $2 AND day = $3
	`

	_, err = tx.Exec(ctx, query, exerciseID, workoutPlanID, workoutDay)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to insert workout plan")
	}

	return &pbw.NilRes{}, nil
}

func (r *RepositoryWorkout) DeleteExerciseWorkoutPlan(ctx context.Context, req *pbw.DeleteExerciseByIdWorkoutPlanReq) (*pbw.NilRes, error) {
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	query := `UPDATE workout_plan_detail
		SET exercises = array_remove(exercises, $1)
		WHERE workout_plan_id = $2 AND day = $3`

	// TODO Check later if we delete one or multiple exercises. Or both
	//result, err := tx.Exec(ctx, query, req.ExerciseId, req.WorkoutPlanId, req.Day)

	var totalRowsAffected int64
	for _, exerciseId := range req.ExerciseId {
		result, err := tx.Exec(ctx, query, exerciseId, req.WorkoutPlanId, req.Day)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to delete exercise from workout plan")
		}
		totalRowsAffected += result.RowsAffected()
	}

	if totalRowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "no exercises found to delete")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	return &pbw.NilRes{}, nil
}

func (r *RepositoryWorkout) UpdateExerciseWorkoutPLan(ctx context.Context, req *pbw.UpdateExerciseByIdWorkoutPlanReq) (*pbw.UpdateExerciseByIdWorkoutPlanRes, error) {
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	query := `UPDATE workout_plan_detail SET `
	var setClauses []string
	var args []interface{}
	argIndex := 1
	updatedFields := make(map[string]string)

	// Build the SET clauses dynamically based on updates
	for _, update := range req.Updates {
		switch update.Field {
		case "exercise_id":
			setClauses = append(setClauses, fmt.Sprintf("exercises = array_replace(exercises, $%d, $%d)", argIndex, argIndex+1))
			args = append(args, update.OldValue, update.NewValue)
			updatedFields["exercise_id"] = update.NewValue
			argIndex += 2
		}
	}

	if len(setClauses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no updates provided")
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE workout_plan_id = $%d", argIndex)
	args = append(args, req.WorkoutPlanId)
	argIndex++

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update exercise in workout plan")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	return &pbw.UpdateExerciseByIdWorkoutPlanRes{
		Success: true,
		Message: "Exercise updated successfully in workout plan",
	}, nil

}

// UpdateWorkoutPlan TODO add proper logic to update all fields ok?
// add transaction to update all tables that are in a workout
func (r *RepositoryWorkout) UpdateWorkoutPlan(ctx context.Context, req *pbw.UpdateWorkoutPlanReq) (*pbw.UpdateWorkoutPlanRes, error) {
	//tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	//if err != nil {
	//	return nil, status.Error(codes.Internal, "failed to start transaction")
	//}
	//
	//defer func() {
	//	if err != nil {
	//		_ = tx.Rollback(ctx)
	//	}
	//}()
	query := `UPDATE workout_plan SET `

	var setClauses []string
	var args []interface{}
	argIndex := 1
	updatedWorkouts := &pbw.XWorkoutPlan{}
	for _, update := range req.Updates {
		switch update.Field {
		case "description":
			setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedWorkouts.Description = update.NewValue
		case "notes":
			setClauses = append(setClauses, fmt.Sprintf("notes = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedWorkouts.Notes = update.NewValue
		case "rating":
			newValue, err := strconv.ParseUint(update.NewValue, 10, 32)
			if err != nil {
				return nil, err
			}
			setClauses = append(setClauses, fmt.Sprintf("rating = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedWorkouts.Rating = uint32(newValue)
		}
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no updates provided")
	}

	query += strings.Join(setClauses, ", ")
	query += ` WHERE id = $` + fmt.Sprintf("%d", argIndex)
	args = append(args, req.WorkoutId)

	_, err := r.pgpool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update exercise: %w", err)
	}

	updatedWorkouts.WorkoutId = req.WorkoutId

	return &pbw.UpdateWorkoutPlanRes{
		Success: true,
		Message: "Workout updated successfully",
		Workout: updatedWorkouts,
	}, nil
}
