package meals

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	ID                 uuid.UUID    `json:"id" db:"id"`                                   // Primary Key
	Name               string       `json:"name" db:"name"`                               // Ingredient name
	Calories           float64      `json:"calories" db:"calories"`                       // Calories
	ServingSize        float64      `json:"serving_size" db:"serving_size"`               // Serving size
	Protein            float64      `json:"protein" db:"protein"`                         // Protein content
	FatTotal           float64      `json:"fat_total" db:"fat_total"`                     // Total fat
	FatSaturated       float64      `json:"fat_saturated" db:"fat_saturated"`             // Saturated fat
	CarbohydratesTotal float64      `json:"carbohydrates_total" db:"carbohydrates_total"` // Total carbohydrates
	Fiber              float64      `json:"fiber" db:"fiber"`                             // Fiber content
	Sugar              float64      `json:"sugar" db:"sugar"`                             // Sugar content
	Sodium             float64      `json:"sodium" db:"sodium"`                           // Sodium content
	Potassium          float64      `json:"potassium" db:"potassium"`                     // Potassium content
	Cholesterol        float64      `json:"cholesterol" db:"cholesterol"`                 // Cholesterol content
	CreatedAt          time.Time    `json:"created_at" db:"created_at"`                   // Timestamp when created
	UpdatedAt          sql.NullTime `json:"updated_at,omitempty" db:"updated_at"`         // Timestamp when updated (nullable)
	UserID             *uuid.UUID   `json:"user_id,omitempty" db:"user_id"`               // Foreign key for user (nullable)
}
