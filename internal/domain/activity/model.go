package activity

import (
	"database/sql"
	"time"
)

type Activity struct {
	ID              sql.NullString  `json:"id,string" db:"id" pg:"default:gen_random_uuid()" protobuf:"bytes,1,opt,name=activity_id,proto3"`
	UserID          sql.NullString  `json:"user_id,string" db:"user_id" protobuf:"bytes,1,opt,name=user_id,proto3"`
	Name            sql.NullString  `json:"name" db:"name" protobuf:"bytes,1,opt,name=name,proto3"`
	CaloriesPerHour sql.NullFloat64 `json:"calories_per_hour" db:"calories_per_hour" protobuf:"bytes,1,opt,name=calories_per_hour,proto3"`
	DurationMinutes sql.NullFloat64 `json:"duration_minutes" db:"duration_minutes" protobuf:"bytes,1,opt,name=duration_in_minutes,proto3"`
	TotalCalories   sql.NullFloat64 `json:"total_calories" db:"total_calories" protobuf:"bytes,1,opt,name=total_calories,proto3"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at" protobuf:"bytes,1,opt,name=created_at,proto3"`
	UpdatedAt       sql.NullTime    `json:"updated_at" db:"updated_at" protobuf:"bytes,1,opt,name=updated_at,proto3"`
}

type ExerciseSession struct {
	ID              string     `json:"id,string" db:"id" pg:"default:gen_random_uuid()" `
	UserID          string     `json:"user_id" db:"user_id"`
	ActivityID      string     `json:"activity_id" db:"activity_id"`
	SessionName     string     `json:"session_name" db:"session_name"`
	StartTime       time.Time  `json:"start_time" db:"start_time"`
	EndTime         time.Time  `json:"end_time" db:"end_time"`
	DurationHours   int        `json:"duration_hours" db:"duration_hours"`
	DurationMinutes int        `json:"duration_minutes" db:"duration_minutes"`
	DurationSeconds int        `json:"duration_seconds" db:"duration_seconds"`
	CaloriesBurned  int        `json:"calories_burned" db:"calories_burned"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
}

type Duration struct {
	Hours       int
	Minutes     int
	Seconds     int
	SessionName string
}

type TotalExerciseSession struct {
	ID                   string    `json:"id,string" db:"id" pg:"default:gen_random_uuid()"`
	UserID               string    `json:"user_id" db:"user_id"`
	ActivityID           string    `json:"activity_id" db:"activity_id"`
	TotalDurationHours   int       `json:"duration_hours" db:"total_duration_hours"`
	TotalDurationMinutes int       `json:"duration_minutes" db:"total_duration_minutes"`
	TotalDurationSeconds int       `json:"duration_seconds" db:"total_duration_seconds"`
	TotalCaloriesBurned  int       `json:"calories_burned" db:"total_calories_burned"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

//type SessionStats struct {
//	ID                     string `json:"id,string" db:"id" pg:"default:gen_random_uuid()"`
//	TotalExerciseSessionID string `json:"total_exercise_session_id,string" db:"total_exercise_session_id"`
//	ActivityID             int       `json:"activity_id" db:"activity_id"`
//	UserID                 int       `json:"user_id" db:"user_id"`
//	SessionName            string    `json:"session_name" db:"session_name"`
//	NumberOfTimes          int       `json:"number_of_times" db:"number_of_times"`
//	TotalDurationHours     int       `json:"duration_hours" db:"total_duration_hours"`
//	TotalDurationMinutes   int       `json:"duration_minutes" db:"total_duration_minutes"`
//	TotalDurationSeconds   int       `json:"duration_seconds" db:"total_duration_seconds"`
//	TotalCaloriesBurned    int       `json:"calories_burned" db:"total_calories_burned"`
//	CreatedAt              time.Time `json:"created_at" db:"created_at"`
//	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
//}

type ExerciseCountStats struct {
	ID                           string    `json:"id,string" db:"id" pg:"default:gen_random_uuid()"`
	ActivityID                   string    `json:"activity_id,string" db:"activity_id"`
	UserID                       string    `json:"user_id,string" db:"user_id"`
	SessionName                  string    `json:"session_name" db:"session_name"`
	NumberOfTimes                int       `json:"number_of_times" db:"number_of_times"`
	TotalExerciseDurationHours   int       `json:"total_duration_hours" db:"total_duration_hours"`
	TotalExerciseDurationMinutes int       `json:"total_duration_minutes" db:"total_duration_minutes"`
	TotalExerciseDurationSeconds int       `json:"total_duration_seconds" db:"total_duration_seconds"`
	TotalExerciseCaloriesBurned  int       `json:"total_calories_burned" db:"total_calories_burned"`
	CreatedAt                    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at" db:"updated_at"`
}

type Status int

const (
	StatusPending Status = iota + 1
	StatusInProgress
	StatusDone
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending:
		return true
	case StatusInProgress:
		return true
	case StatusDone:
		return true
	}
	return false
}
