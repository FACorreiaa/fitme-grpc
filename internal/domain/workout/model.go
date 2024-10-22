package workout

import (
	"database/sql"
	"time"
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
