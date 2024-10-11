package activity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	pba "github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type ActivityRepository struct {
	pba.UnimplementedActivityServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewActivityRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *ActivityRepository {
	return &ActivityRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func (a *ActivityRepository) GetActivity(ctx context.Context, req *pba.GetActivityReq) (*pba.GetActivityRes, error) {
	activities := make([]*pba.XActivity, 0)
	query := `SELECT id, user_id, name,
					duration_minutes, total_calories, calories_per_hour,
					created_at, updated_at
			FROM activity`

	rows, err := a.pgpool.Query(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pba.GetActivityRes{
				Success: false,
				Message: "No activities found",
				Response: &pba.BaseResponse{
					Upstream:  "activity-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, fmt.Errorf("activities not found %w", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to query activities: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a pba.XActivity
		var userId, name sql.NullString
		var durationInMinutes, totalCalories, caloriesPerHour sql.NullFloat64

		var createdAt time.Time
		var updatedAt sql.NullTime

		err := rows.Scan(
			&a.ActivityId, &userId, &name, &durationInMinutes, &totalCalories, &caloriesPerHour,
			&createdAt, &updatedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &pba.GetActivityRes{
					Success: false,
					Message: "No activities found",
					Response: &pba.BaseResponse{
						Upstream:  "activity-service",
						RequestId: domain.GenerateRequestID(ctx),
					},
				}, fmt.Errorf("activities not found: %w", err)
			}
			return nil, status.Error(codes.Internal, "Internal server error")
		}

		// Populate activity fields
		a.UserId = userId.String
		a.Name = name.String

		if durationInMinutes.Valid {
			a.DurationInMinutes = float32(durationInMinutes.Float64)
		}
		if totalCalories.Valid {
			a.TotalCalories = float32(totalCalories.Float64)
		}
		if caloriesPerHour.Valid {
			a.CaloriesPerHour = float32(caloriesPerHour.Float64)
		}

		// Use the time values directly to create protobuf timestamps
		a.CreatedAt = timestamppb.New(createdAt)

		if updatedAt.Valid {
			a.UpdatedAt = timestamppb.New(updatedAt.Time)
		} else {
			a.UpdatedAt = timestamppb.New(time.Now()) // default value if updated_at is NULL
		}

		activities = append(activities, &a)
	}

	if len(activities) == 0 {
		return &pba.GetActivityRes{
			Success: false,
			Message: "No activities found",
			Response: &pba.BaseResponse{
				Upstream:  "activity-service",
				RequestId: domain.GenerateRequestID(ctx),
			},
		}, nil
	}

	return &pba.GetActivityRes{
		Success:  true,
		Message:  "Activities retrieved successfully",
		Activity: activities,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (a *ActivityRepository) GetActivitiesByID(ctx context.Context, req *pba.GetActivityIDReq) (*pba.GetActivityIDRes, error) {
	activity := &pba.XActivity{}
	activityID := req.PublicId
	var createdAt time.Time
	var updatedAt sql.NullTime

	var userId, name sql.NullString
	var durationInMinutes, totalCalories, caloriesPerHour sql.NullFloat64

	if activityID == "" {
		return nil, status.Error(codes.InvalidArgument, "activity ID is required")
	}

	query := `SELECT 	id, user_id, name, duration_minutes,
       					total_calories, calories_per_hour, created_at,
       					updated_at
			   FROM activity
			   WHERE id = $1`

	err := a.pgpool.QueryRow(ctx, query, activityID).Scan(
		&activity.ActivityId, &userId, &name, &durationInMinutes, &totalCalories, &caloriesPerHour,
		&createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pba.GetActivityIDRes{
				Success: false,
				Message: "No activities found",
				Response: &pba.BaseResponse{
					Upstream:  "activity-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, fmt.Errorf("activities not found: %w", err)
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	activity.UserId = userId.String
	activity.Name = name.String
	activity.DurationInMinutes = float32(durationInMinutes.Float64)
	activity.TotalCalories = float32(totalCalories.Float64)
	activity.CaloriesPerHour = float32(caloriesPerHour.Float64)

	activity.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		activity.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &pba.GetActivityIDRes{
		Success:  true,
		Message:  "Activity retrieved successfully",
		Activity: activity,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (a *ActivityRepository) GetActivitiesByName(ctx context.Context, req *pba.GetActivityNameReq) (*pba.GetActivityNameRes, error) {
	activity := &pba.XActivity{}
	nameReq := req.PublicId
	var createdAt time.Time
	var updatedAt sql.NullTime

	var userId, name sql.NullString
	var durationInMinutes, totalCalories, caloriesPerHour sql.NullFloat64

	log.Printf("Searching for activity with name: '%s'", nameReq)

	if nameReq == "" {
		return nil, status.Error(codes.InvalidArgument, "activity ID is required")
	}

	query := `SELECT 	id, user_id, name, duration_minutes,
       					total_calories, calories_per_hour, created_at,
       					updated_at
			   FROM activity
			   WHERE name LIKE '%' || $1 || '%'`

	err := a.pgpool.QueryRow(ctx, query, nameReq).Scan(
		&activity.ActivityId, &userId, &name, &durationInMinutes, &totalCalories, &caloriesPerHour,
		&createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pba.GetActivityNameRes{
				Success: false,
				Message: "No activities found",
				Response: &pba.BaseResponse{
					Upstream:  "activity-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, fmt.Errorf("activities not found: %w", err)
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if userId.Valid {
		activity.UserId = userId.String
	}
	if name.Valid {
		activity.Name = name.String
	}
	if durationInMinutes.Valid {
		activity.DurationInMinutes = float32(durationInMinutes.Float64)
	}
	if totalCalories.Valid {
		activity.TotalCalories = float32(totalCalories.Float64)
	}
	if caloriesPerHour.Valid {
		activity.CaloriesPerHour = float32(caloriesPerHour.Float64)
	}

	activity.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		activity.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &pba.GetActivityNameRes{
		Success:  true,
		Message:  "Activity retrieved successfully",
		Activity: activity,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

// implement later
//func (a *ActivityRepository) SaveSession(ctx context.Context, req *pba.XExerciseSession) error {
//	query := `
//		INSERT INTO exercise_session
//		    (user_id, activity_id, session_name, start_time,
//		     end_time, duration_hours, duration_minutes, duration_seconds,
//		     calories_burned, created_at)
//		VALUES (:user_id, :activity_id, :session_name, :start_time,
//		        :end_time, :duration_hours, :duration_minutes, :duration_seconds,
//		        :calories_burned, :created_at)
//		RETURNING id;`
//
//	es := &pba.XExerciseSession{}
//	var createdAt time.Time
//
//	rows, err := a.pgpool.Query(ctx, query, req.UserId, req.ActivityId, req.SessionName, req.StartTime, req.EndTime,
//		req.DurationHours, req.DurationMinutes, req.DurationSeconds, req.CaloriesBurned, req.CreatedAt)
//
//	defer rows.Close()
//	if err != nil {
//		log.Printf("Query execution error: %v", err) // Log detailed error
//		return fmt.Errorf("failed to execute query: %w", err)
//	}
//
//	if rows.Next() {
//		err = rows.Scan(
//			&es.UserId, &es.ActivityId, &es.SessionName, &es.StartTime, &es.EndTime, &es.DurationHours,
//			&es.DurationMinutes, &es.DurationSeconds, &es.CaloriesBurned, &es.CreatedAt,
//		)
//
//		if err != nil {
//			return fmt.Errorf("failed to scan row: %w", err)
//		}
//
//		// Convert `createdAt` (Go time.Time) to Protobuf Timestamp
//		es.CreatedAt = timestamppb.New(createdAt)
//	}
//
//	defer rows.Close()
//
//	return nil
//}

func (a *ActivityRepository) SaveSession(ctx context.Context, req *pba.XExerciseSession) error {
	query := `
		INSERT INTO exercise_session
		    (user_id, activity_id, session_name, start_time,
		     end_time, duration_hours, duration_minutes, duration_seconds,
		     calories_burned, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING *;`

	var sessionID uuid.UUID

	// Execute the query and get the inserted session ID
	err := a.pgpool.QueryRow(ctx, query,
		req.UserId, req.ActivityId, req.SessionName, req.StartTime, req.EndTime,
		req.DurationHours, req.DurationMinutes, req.DurationSeconds, req.CaloriesBurned, req.CreatedAt.AsTime(),
	).Scan(&sessionID)

	if err != nil {
		log.Printf("Query execution error: %v", err) // Log detailed error
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (a *ActivityRepository) GetUserExerciseSession(ctx context.Context, req *pba.GetUserExerciseSessionReq) (*pba.GetUserExerciseSessionRes, error) {
	//sessionStats := make([]*pba.XExerciseSession, 0)
	userID := req.PublicId
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	query := `
				SELECT es.session_name, es.activity_id,
       				COUNT(*) as number_of_times,
       				SUM(es.duration_seconds) as total_duration_seconds,
					SUM(es.duration_minutes) as total_duration_minutes,
					SUM(es.duration_hours) as total_duration_hours,
					SUM(es.calories_burned) as total_calories_burned
              FROM exercise_session es
              WHERE user_id = $1
              GROUP BY session_name, activity_id
              ORDER BY number_of_times DESC
              LIMIT 1
			`

	var (
		sessionName          string
		activityID           string
		numberOfTimes        int64
		totalDurationSeconds int64
		totalDurationMinutes int64
		totalDurationHours   int64
		totalCaloriesBurned  int64
	)

	err := a.pgpool.QueryRow(ctx, query, userID).Scan(
		&sessionName, &activityID, &numberOfTimes, &totalDurationSeconds,
		&totalDurationMinutes, &totalDurationHours, &totalCaloriesBurned,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "no exercise sessions found for the user")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return &pba.GetUserExerciseSessionRes{
				Success: false,
				Message: "No exercise session found",
				Response: &pba.BaseResponse{
					Upstream:  "activity-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, fmt.Errorf("exercise session not found: %w", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise session: %v", err)
	}

	sessionStat := &pba.XExerciseSession{
		SessionName:     sessionName,
		ActivityId:      activityID,
		NumberOfTimes:   strconv.FormatInt(numberOfTimes, 10),
		DurationSeconds: uint32(totalDurationSeconds),
		DurationMinutes: uint32(totalDurationMinutes),
		DurationHours:   uint32(totalDurationHours),
		CaloriesBurned:  uint32(totalCaloriesBurned),
	}

	return &pba.GetUserExerciseSessionRes{
		Success: true,
		Message: "Activity retrieved successfully",
		Session: sessionStat,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (a *ActivityRepository) GetUserExerciseTotalData(ctx context.Context, req *pba.GetUserExerciseTotalDataReq) (*pba.GetUserExerciseTotalDataRes, error) {
	userID := req.PublicId
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	query := `SELECT duration_hours, duration_minutes, duration_seconds, calories_burned, session_name
			  FROM exercise_session WHERE user_id = $1`

	var exerciseSessions []ExerciseSession

	rows, err := a.pgpool.Query(ctx, query, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise sessions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var eSession ExerciseSession
		if err := rows.Scan(&eSession.DurationHours, &eSession.DurationMinutes, &eSession.DurationSeconds,
			&eSession.CaloriesBurned, &eSession.SessionName); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan exercise session: %v", err)
		}
		exerciseSessions = append(exerciseSessions, eSession)
	}

	if len(exerciseSessions) == 0 {
		return nil, status.Errorf(codes.NotFound, "no exercise sessions found for user")
	}

	totalDuration, totalCaloriesBurned := calculateTotal(exerciseSessions)

	if err := saveTotalExerciseSession(ctx, a.pgpool, userID, totalDuration, totalCaloriesBurned, exerciseSessions[0].SessionName); err != nil {
		return nil, err
	}

	return &pba.GetUserExerciseTotalDataRes{
		Success: true,
		Message: "Total exercise session saved",
		Session: &pba.XTotalExerciseSession{
			UserId:               userID,
			TotalDurationHours:   uint32(totalDuration.Hours),
			TotalDurationMinutes: uint32(totalDuration.Minutes),
			TotalDurationSeconds: uint32(totalDuration.Seconds),
			TotalCaloriesBurned:  uint32(totalCaloriesBurned),
			SessionName:          exerciseSessions[0].SessionName,
		},
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func calculateTotal(sessions []ExerciseSession) (Duration, int) {
	totalDuration := Duration{}
	totalCaloriesBurned := 0

	for _, session := range sessions {
		totalDuration.Hours += session.DurationHours
		totalDuration.Minutes += session.DurationMinutes
		totalDuration.Seconds += session.DurationSeconds
		totalCaloriesBurned += session.CaloriesBurned
	}

	totalDuration.Minutes += totalDuration.Seconds / 60
	totalDuration.Seconds = totalDuration.Seconds % 60
	totalDuration.Hours += totalDuration.Minutes / 60
	totalDuration.Minutes = totalDuration.Minutes % 60

	return totalDuration, totalCaloriesBurned
}

func saveTotalExerciseSession(ctx context.Context, pgpool *pgxpool.Pool, userID string, totalDuration Duration, totalCaloriesBurned int, sessionName string) error {
	_, err := pgpool.Exec(ctx, `
		INSERT INTO total_exercise_session (user_id, total_duration_hours, total_duration_minutes, total_duration_seconds, total_calories_burned, session_name)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id)
		DO UPDATE SET
			total_duration_hours = EXCLUDED.total_duration_hours,
			total_duration_minutes = EXCLUDED.total_duration_minutes,
			total_duration_seconds = EXCLUDED.total_duration_seconds,
			total_calories_burned = EXCLUDED.total_calories_burned,
			session_name = EXCLUDED.session_name,
			updated_at = NOW()
	`, userID, totalDuration.Hours, totalDuration.Minutes, totalDuration.Seconds, totalCaloriesBurned, sessionName)

	if err != nil {
		return status.Errorf(codes.Internal, "failed to save total exercise session: %v", err)
	}
	return nil
}

// GetUserExerciseSessionStats
func (a *ActivityRepository) GetUserExerciseSessionStats(ctx context.Context, req *pba.GetUserExerciseSessionStatsReq) (*pba.GetUserExerciseSessionStatsRes, error) {
	id := req.PublicId
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	sessionStats := make([]*ExerciseCountStats, 0)
	query := `SELECT es.session_name, es.activity_id,
       				COUNT(*) as number_of_times,
       				SUM(es.duration_seconds) as total_duration_seconds,
					SUM(es.duration_minutes) as total_duration_minutes,
					SUM(es.duration_hours) as total_duration_hours,
					SUM(es.calories_burned) as total_calories_burned
              FROM exercise_session es
              WHERE user_id = $1
              GROUP BY session_name, activity_id
              ORDER BY number_of_times DESC
              LIMIT 1`

	rows, err := a.pgpool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stat ExerciseCountStats
		if err := rows.Scan(&stat.SessionName, &stat.ActivityID, &stat.NumberOfTimes, &stat.TotalExerciseDurationSeconds, &stat.TotalExerciseDurationMinutes, &stat.TotalExerciseDurationHours, &stat.TotalExerciseCaloriesBurned); err != nil {
			return nil, fmt.Errorf("failed to scan exercise stats: %w", err)
		}
		sessionStats = append(sessionStats, &stat)
	}

	// Check if no rows were found
	if len(sessionStats) == 0 {
		return nil, fmt.Errorf("no exercise stats found for user")
	}

	pbSessionStats := make([]*pba.XExerciseCountStats, 0, len(sessionStats))
	for _, stat := range sessionStats {
		pbStat := &pba.XExerciseCountStats{
			SessionName:                  stat.SessionName,
			ActivityId:                   stat.ActivityID,
			NumberOfTimes:                uint32(stat.NumberOfTimes),
			TotalExerciseDurationSeconds: uint32(stat.TotalExerciseDurationSeconds),
			TotalExerciseDurationMinutes: uint32(stat.TotalExerciseDurationMinutes),
			TotalExerciseDurationHours:   uint32(stat.TotalExerciseDurationHours),
			TotalExerciseCaloriesBurned:  uint32(stat.TotalExerciseCaloriesBurned),
		}
		pbSessionStats = append(pbSessionStats, pbStat)
	}

	return &pba.GetUserExerciseSessionStatsRes{
		Success:       true,
		Message:       "Get user exercise session stats successful",
		ExerciseCount: pbSessionStats,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

// GetExerciseSessionStats this might not be needed
//func (a *ActivityRepository) GetExerciseSessionStats(ctx context.Context, req *pba.GetExerciseSessionStatsOccurrenceReq) (*pba.GetExerciseSessionStatsOccurrenceRes, error) {
//	sessionStats := make([]*pba.XTotalExerciseSession, 0)
//	id := req.PublicId
//	query := `SELECT es.session_name, es.activity_id,
//       				COUNT(*) as number_of_times,
//       				SUM(es.duration_seconds) as total_duration_seconds,
//					SUM(es.duration_minutes) as total_duration_minutes,
//					SUM(es.duration_hours) as total_duration_hours,
//					SUM(es.calories_burned) as total_calories_burned
//              FROM exercise_session es
//              WHERE user_id = $1
//              GROUP BY session_name, activity_id
//              ORDER BY number_of_times DESC
//              LIMIT 1`
//
//	rows, err := a.pgpool.Query(ctx, query, id)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var stat ExerciseCountStats
//		if err := rows.Scan(&stat.SessionName, &stat.ActivityID, &stat.NumberOfTimes, &stat.TotalExerciseDurationSeconds, &stat.TotalExerciseDurationMinutes, &stat.TotalExerciseDurationHours, &stat.TotalExerciseCaloriesBurned); err != nil {
//			return nil, fmt.Errorf("failed to scan exercise stats: %w", err)
//		}
//		sessionStats = append(sessionStats, &stat)
//	}
//
//	// Check if no rows were found
//	if len(sessionStats) == 0 {
//		return nil, fmt.Errorf("no exercise stats found for user")
//	}
//
//	pbSessionStats := make([]*pba.XExerciseCountStats, 0, len(sessionStats))
//	for _, stat := range sessionStats {
//		pbStat := &pba.XExerciseCountStats{
//			SessionName:                  stat.SessionName,
//			ActivityId:                   stat.ActivityId,
//			NumberOfTimes:                uint32(stat.NumberOfTimes),
//			TotalExerciseDurationSeconds: uint32(stat.TotalDurationSeconds),
//			TotalExerciseDurationMinutes: uint32(stat.TotalDurationMinutes),
//			TotalExerciseDurationHours:   uint32(stat.TotalDurationHours),
//			TotalExerciseCaloriesBurned:  uint32(stat.TotalCaloriesBurned),
//		}
//		pbSessionStats = append(pbSessionStats, pbStat)
//	}
//
//	return &pba.GetExerciseSessionStatsOccurrenceRes{
//		Success:       true,
//		Message:       "Get user exercise session stats successful",
//		ExerciseCount: pbSessionStats,
//		Response: &pba.BaseResponse{
//			Upstream:  "activity-service",
//			RequestId: domain.GenerateRequestID(ctx),
//		},
//	}, nil
//}
