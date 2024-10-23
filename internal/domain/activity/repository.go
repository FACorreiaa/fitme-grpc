package activity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

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

type RepositoryActivity struct {
	pba.UnimplementedActivityServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewRepositoryActivity(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *RepositoryActivity {
	return &RepositoryActivity{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func float64OrZero(nullFloat sql.NullFloat64) float64 {
	if nullFloat.Valid {
		return nullFloat.Float64
	}
	return 0
}

func stringOrEmpty(nullString sql.NullString) string {
	if nullString.Valid {
		return nullString.String
	}
	return ""

}

func nullTimeToTimestamppb(nt sql.NullTime) *timestamppb.Timestamp {
	if nt.Valid {
		return timestamppb.New(nt.Time)
	}
	return nil
}

func (a *RepositoryActivity) GetActivity(ctx context.Context, req *pba.GetActivityReq) (*pba.GetActivityRes, error) {
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
		ac := &Activity{}
		a := pba.XActivity{}

		err := rows.Scan(
			&ac.ID, &ac.ID, &ac.Name, &ac.DurationMinutes, &ac.TotalCalories, &ac.CaloriesPerHour,
			&ac.CreatedAt, &ac.UpdatedAt,
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

		createdAt := timestamppb.New(ac.CreatedAt)
		var updatedAt sql.NullTime
		if ac.UpdatedAt.Valid {
			updatedAt = ac.UpdatedAt
		} else {
			updatedAt = sql.NullTime{Valid: false}
		}

		a.ActivityId = ac.ID.String
		a.UserId = ac.UserID.String
		a.Name = ac.Name.String
		a.DurationInMinutes = float32(float64OrZero(ac.DurationMinutes))
		a.TotalCalories = float32(float64OrZero(ac.TotalCalories))
		a.CaloriesPerHour = float32(float64OrZero(ac.CaloriesPerHour))
		a.CreatedAt = createdAt
		a.UpdatedAt = nullTimeToTimestamppb(updatedAt)

		activityProto := &pba.XActivity{
			ActivityId:        a.ActivityId,
			UserId:            a.UserId,
			Name:              a.Name,
			DurationInMinutes: a.DurationInMinutes,
			TotalCalories:     a.TotalCalories,
			CaloriesPerHour:   a.CaloriesPerHour,
			CreatedAt:         a.CreatedAt,
		}

		activities = append(activities, activityProto)
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

func (a *RepositoryActivity) GetActivitiesByName(ctx context.Context, req *pba.GetActivityNameReq) (*pba.GetActivityNameRes, error) {
	activity := &pba.XActivity{}
	nameReq := req.PublicId
	ac := &Activity{}

	if nameReq == "" {
		return nil, status.Error(codes.InvalidArgument, "activity ID is required")
	}

	query := `SELECT 	id, user_id, name, duration_minutes,
       					total_calories, calories_per_hour, created_at,
       					updated_at
			   FROM activity
			   WHERE name LIKE '%' || $1 || '%'`

	err := a.pgpool.QueryRow(ctx, query, nameReq).Scan(
		&ac.ID, &ac.UserID, &ac.Name, &ac.DurationMinutes, &ac.TotalCalories, &ac.CaloriesPerHour,
		&ac.CreatedAt, &ac.UpdatedAt,
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

	createdAt := timestamppb.New(ac.CreatedAt)
	var updatedAt sql.NullTime
	if ac.UpdatedAt.Valid {
		updatedAt = ac.UpdatedAt
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	activity.ActivityId = ac.ID.String
	activity.UserId = ac.UserID.String
	activity.Name = ac.Name.String
	activity.DurationInMinutes = float32(float64OrZero(ac.DurationMinutes))
	activity.TotalCalories = float32(float64OrZero(ac.TotalCalories))
	activity.CaloriesPerHour = float32(float64OrZero(ac.CaloriesPerHour))
	activity.CreatedAt = createdAt
	activity.UpdatedAt = nullTimeToTimestamppb(updatedAt)

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

func (a *RepositoryActivity) GetActivitiesByID(ctx context.Context, req *pba.GetActivityIDReq) (*pba.GetActivityIDRes, error) {
	activity := &pba.XActivity{}
	ac := &Activity{}
	activityID := req.PublicId

	if activityID == "" {
		return nil, status.Error(codes.InvalidArgument, "activity ID is required")
	}

	query := `SELECT 	id, user_id, name, duration_minutes,
       					total_calories, calories_per_hour, created_at,
       					updated_at
			   FROM activity
			   WHERE id = $1`

	err := a.pgpool.QueryRow(ctx, query, activityID).Scan(
		&ac.ID, &ac.UserID, &ac.Name, &ac.DurationMinutes, &ac.TotalCalories, &ac.CaloriesPerHour,
		&ac.CreatedAt, &ac.UpdatedAt,
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

	createdAt := timestamppb.New(ac.CreatedAt)
	var updatedAt sql.NullTime
	if ac.UpdatedAt.Valid {
		updatedAt = ac.UpdatedAt
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	activity.ActivityId = ac.ID.String
	activity.UserId = ac.UserID.String
	activity.Name = ac.Name.String
	activity.DurationInMinutes = float32(float64OrZero(ac.DurationMinutes))
	activity.TotalCalories = float32(float64OrZero(ac.TotalCalories))
	activity.CaloriesPerHour = float32(float64OrZero(ac.CaloriesPerHour))
	activity.CreatedAt = createdAt
	activity.UpdatedAt = nullTimeToTimestamppb(updatedAt)

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

// implement later
//func (a *RepositoryActivity) SaveSession(ctx context.Context, req *pba.XExerciseSession) error {
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

func (a *RepositoryActivity) SaveSession(ctx context.Context, req *pba.XExerciseSession) error {
	query := `
		INSERT INTO exercise_session
		    (user_id, activity_id, session_name, start_time,
		     end_time, duration_hours, duration_minutes, duration_seconds,
		     calories_burned, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id;`

	var sessionID uuid.UUID

	startTime := req.StartTime.AsTime()
	fmt.Printf("%#v", startTime)
	endTime := req.EndTime.AsTime()
	fmt.Printf("%#v", endTime)
	createdAt := req.CreatedAt.AsTime()
	fmt.Printf("%#v", createdAt)

	// Execute the query and get the inserted session ID
	err := a.pgpool.QueryRow(ctx, query,
		req.UserId, req.ActivityId, req.SessionName, startTime, endTime,
		req.DurationHours, req.DurationMinutes, req.DurationSeconds, req.CaloriesBurned, createdAt,
	).Scan(&sessionID)

	if err != nil {
		log.Printf("Query execution error: %v", err) // Log detailed error
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (a *RepositoryActivity) GetUserExerciseSession(ctx context.Context, req *pba.GetUserExerciseSessionReq) (*pba.GetUserExerciseSessionRes, error) {
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

	//var (
	//	sessionName          string
	//	activityID           string
	//	numberOfTimes        int64
	//	totalDurationSeconds int64
	//	totalDurationMinutes int64
	//	totalDurationHours   int64
	//	totalCaloriesBurned  int64
	//)

	e := &ExerciseCountStats{}

	err := a.pgpool.QueryRow(ctx, query, userID).Scan(
		&e.SessionName, &e.ActivityID, &e.NumberOfTimes, &e.TotalExerciseDurationHours,
		&e.TotalExerciseDurationMinutes, &e.TotalExerciseDurationSeconds, &e.TotalExerciseCaloriesBurned,
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
		SessionName:     e.SessionName,
		ActivityId:      e.ActivityID,
		NumberOfTimes:   strconv.Itoa(e.NumberOfTimes),
		DurationSeconds: uint32(e.TotalExerciseDurationSeconds),
		DurationMinutes: uint32(e.TotalExerciseDurationMinutes),
		DurationHours:   uint32(e.TotalExerciseDurationHours),
		CaloriesBurned:  uint32(e.TotalExerciseCaloriesBurned),
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

func (a *RepositoryActivity) GetUserExerciseTotalData(ctx context.Context, req *pba.GetUserExerciseTotalDataReq) (*pba.GetUserExerciseTotalDataRes, error) {
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

// GetUserExerciseSessionStats review
func (a *RepositoryActivity) GetUserExerciseSessionStats(ctx context.Context, req *pba.GetUserExerciseSessionStatsReq) (*pba.GetUserExerciseSessionStatsRes, error) {
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
//func (a *RepositoryActivity) GetExerciseSessionStats(ctx context.Context, req *pba.GetExerciseSessionStatsOccurrenceReq) (*pba.GetExerciseSessionStatsOccurrenceRes, error) {
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

func (a *RepositoryActivity) DeleteExerciseSession(ctx context.Context, req *pba.DeleteExerciseSessionReq) (*pba.NilRes, error) {
	if req.PublicId == "" {
		return nil, status.Error(codes.InvalidArgument, "public_id is required")
	}

	query := `DELETE FROM exercise_session WHERE id = $1`

	_, err := a.pgpool.Exec(ctx, query, req.PublicId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete exercise session: %w", err)
	}

	return &pba.NilRes{}, nil
}

func (a *RepositoryActivity) DeleteAllExercisesSession(ctx context.Context, req *pba.DeleteAllExercisesSessionReq) (*pba.NilRes, error) {
	query := `DELETE FROM exercise_session WHERE user_id = $1`

	_, err := a.pgpool.Exec(ctx, query, req.PublicId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete all exercise sessions for user: %w", err)
	}

	return &pba.NilRes{}, nil
}
