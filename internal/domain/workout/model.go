package workout

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Exercises struct {
	ID            string       `json:"id,string" db:"id" pg:"default:gen_random_uuid()"`
	Name          string       `json:"name" db:"name"`
	ExerciseType  string       `json:"type" db:"type"`
	MuscleGroup   string       `json:"muscle" db:"muscle"`
	Equipment     string       `json:"equipment" db:"equipment"`
	Difficulty    string       `json:"difficulty" db:"difficulty"`
	Instructions  string       `json:"instructions" db:"instructions"`
	Video         string       `json:"video" db:"video"`
	CustomCreated bool         `json:"custom_created" db:"custom_created"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     sql.NullTime `json:"updated_at" db:"updated_at"`
}

type WorkoutPlan struct {
	ID          string           `json:"id,string" db:"id" pg:"default:gen_random_uuid()"`
	UserID      int              `json:"user_id" db:"user_id"`
	Description string           `json:"description" db:"description"`
	Notes       string           `json:"notes" db:"notes"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time       `json:"updated_at" db:"updated_at"`
	Rating      int              `json:"rating" db:"rating"`
	WorkoutDays []WorkoutPlanDay `json:"workoutDays" db:"-"`
}

type WorkoutDay struct {
	ID            string      `json:"id,string" db:"id" pg:"default:gen_random_uuid()"`
	WorkoutPlanID string      `json:"workout_plan_id" db:"workout_plan_id"`
	Day           string      `json:"day" db:"day"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time  `json:"updated_at" db:"updated_at"`
	Exercises     []Exercises `json:"exercises" db:"exercises"`
}

type PlanDay struct {
	Day         string   `json:"day"`
	ExerciseIDs []string `json:"exercises"`
}

type CreateWorkoutPlanRequest struct {
	WorkoutPlan WorkoutPlan `json:"workoutPlan"`
	Plan        []PlanDay   `json:"plan"`
}

type WorkoutPlanDetail struct {
	ID            string    `db:"id"`
	WorkoutPlanID string    `db:"workout_plan_id"`
	Day           string    `db:"day"`
	Exercises     []string  `db:"exercises"`
	CreatedAt     time.Time `db:"created_at"`
}

type WorkoutPlanDay struct {
	Day       string      `json:"day" db:"day"`
	Exercises []Exercises `json:"exercises" db:"exercises"`
}

type WorkoutPlanResponse struct {
	WorkoutPlanID string                   `json:"workout_plan_id" db:"workout_plan_id"`
	UserID        int                      `json:"user_id" db:"user_id"`
	Description   string                   `json:"description" db:"description"`
	WorkoutDays   []WorkoutDayResponse     `json:"workoutDays" db:"-"`
	Notes         string                   `json:"notes" db:"notes"`
	CreatedAt     time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time               `json:"updated_at" db:"updated_at"`
	Rating        int                      `json:"rating" db:"rating"`
	Day           string                   `json:"day" db:"day"`
	Exercises     pgtype.FlatArray[string] `json:"exercises" db:"exercises"`
}

type WorkoutDayResponse struct {
	Day       string                   `json:"day" db:"day"`
	Exercises pgtype.FlatArray[string] `json:"exercises" db:"exercises" `
}

type WorkoutExerciseDay struct {
	ID            string       `json:"id,string" db:"id" pg:"default:gen_random_uuid()"`
	Name          string       `json:"name" db:"name"`
	ExerciseType  string       `json:"type" db:"type"`
	MuscleGroup   string       `json:"muscle" db:"muscle"`
	Equipment     string       `json:"equipment" db:"equipment"`
	Difficulty    string       `json:"difficulty" db:"difficulty"`
	Instructions  string       `json:"instructions" db:"instructions"`
	Video         string       `json:"video" db:"video"`
	CustomCreated bool         `json:"custom_created" db:"custom_created"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     sql.NullTime `json:"updated_at" db:"updated_at"`
	Day           string       `json:"day" db:"day"`
}
