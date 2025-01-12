package meals

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	ID                 uuid.UUID    `protobuf:"bytes,1,opt,name=ingredient_id,json=ingredientId" db:"id"`                                // Primary Key
	Name               string       `protobuf:"bytes,2,opt,name=name" db:"name"`                                                         // Ingredient name
	Calories           float64      `protobuf:"fixed64,3,opt,name=calories" db:"calories"`                                               // Calories
	ServingSize        float64      `protobuf:"fixed64,4,opt,name=serving_size,json=servingSize" db:"serving_size"`                      // Serving size
	Protein            float64      `protobuf:"fixed64,5,opt,name=protein" db:"protein"`                                                 // Protein content
	FatTotal           float64      `protobuf:"fixed64,6,opt,name=fat_total,json=fatTotal" db:"fat_total"`                               // Total fat
	FatSaturated       float64      `protobuf:"fixed64,7,opt,name=fat_saturated,json=fatSaturated" db:"fat_saturated"`                   // Saturated fat
	CarbohydratesTotal float64      `protobuf:"fixed64,8,opt,name=carbohydrates_total,json=carbohydratesTotal" db:"carbohydrates_total"` // Total carbohydrates
	Fiber              float64      `protobuf:"fixed64,9,opt,name=fiber" db:"fiber"`                                                     // Fiber content
	Sugar              float64      `protobuf:"fixed64,10,opt,name=sugar" db:"sugar"`                                                    // Sugar content
	Sodium             float64      `protobuf:"fixed64,11,opt,name=sodium" db:"sodium"`                                                  // Sodium content
	Potassium          float64      `protobuf:"fixed64,12,opt,name=potassium" db:"potassium"`                                            // Potassium content
	Cholesterol        float64      `protobuf:"fixed64,13,opt,name=cholesterol" db:"cholesterol"`                                        // Cholesterol content
	CreatedAt          time.Time    `protobuf:"bytes,14,opt,name=created_at,json=createdAt" db:"created_at"`                             // Timestamp when created
	UpdatedAt          sql.NullTime `protobuf:"bytes,15,opt,name=updated_at,json=updatedAt" db:"updated_at"`                             // Timestamp when updated (nullable)
	UserID             *uuid.UUID   `protobuf:"bytes,16,opt,name=user_id,json=userId" db:"user_id"`                                      // Foreign key for user (nullable)
}
type Meal struct {
	ID              uuid.UUID       `protobuf:"bytes,1,opt,name=meal_id,json=meal_id" db:"id"` // Primary Key
	MealNumber      int             `protobuf:"fixed32,3,opt,name=meal_number,json=mealNumber" db:"number"`
	MealDescription string          `protobuf:"bytes,4,opt,name=meal_description,json=mealDescription" db:"description"` // Description
	Ingredients     []Ingredient    `protobuf:"bytes,4,rep,name=ingredients" db:"ingredients"`                           // Ingredients
	CreatedAt       time.Time       `protobuf:"bytes,5,opt,name=created_at,json=createdAt" db:"created_at"`              // Timestamp when created
	UpdatedAt       sql.NullTime    `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt" db:"updated_at"`              // Timestamp when updated (nullable)
	UserID          *uuid.UUID      `protobuf:"bytes,7,opt,name=user_id,json=userId" db:"user_id"`
	TotalMacros     *TotalNutrients `protobuf:"bytes,8,rep,name=total_macros,json=totalMacros" db:"total_macros"`
}

type MealIngredient struct {
	MealID             uuid.UUID   `protobuf:"bytes,1,opt,name=meal_id,json=mealId" db:"meal_id"`
	IngredientID       []uuid.UUID `protobuf:"bytes,2,opt,name=ingredient_id,json=ingredientId" db:"ingredient_id"`
	Quantity           float64     `protobuf:"fixed64,3,opt,name=quantity" db:"quantity"`
	Calories           float64     `protobuf:"fixed64,4,opt,name=calories" db:"calories"`
	Protein            float64     `protobuf:"fixed64,5,opt,name=protein" db:"protein"`
	FatTotal           float64     `protobuf:"fixed64,6,opt,name=fat_total,json=fatTotal" db:"fat_total"`
	FatSaturated       float64     `protobuf:"fixed64,7,opt,name=fat_saturated,json=fatSaturated" db:"fat_saturated"`
	CarbohydratesTotal float64     `protobuf:"fixed64,8,opt,name=carbohydrates_total,json=carbohydratesTotal" db:"carbohydrates_total"`
	Fiber              float64     `protobuf:"fixed64,9,opt,name=fiber" db:"fiber"`
	Sugar              float64     `protobuf:"fixed64,10,opt,name=sugar" db:"sugar"`
	Sodium             float64     `protobuf:"fixed64,11,opt,name=sodium" db:"sodium"`
	Potassium          float64     `protobuf:"fixed64,12,opt,name=potassium" db:"potassium"`
	Cholesterol        float64     `protobuf:"fixed64,13,opt,name=cholesterol" db:"cholesterol"`
	CreatedAt          time.Time   `protobuf:"bytes,4,opt,name=created_at,json=createdAt" db:"created_at"`
	UpdatedAt          time.Time   `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt" db:"updated_at"`
}

type TotalNutrients struct {
	Calories           float64 `protobuf:"fixed64,1,opt,name=calories" db:"calories"`
	Protein            float64 `protobuf:"fixed64,2,opt,name=protein" db:"protein"`
	FatTotal           float64 `protobuf:"fixed64,3,opt,name=fat_total,json=fatTotal" db:"fat_total"`
	FatSaturated       float64 `protobuf:"fixed64,4,opt,name=fat_saturated,json=fatSaturated" db:"fat_saturated"`
	CarbohydratesTotal float64 `protobuf:"fixed64,5,opt,name=carbohydrates_total,json=carbohydratesTotal" db:"carbohydrates_total"`
	Fiber              float64 `protobuf:"fixed64,6,opt,name=fiber" db:"fiber"`
	Sugar              float64 `protobuf:"fixed64,7,opt,name=sugar" db:"sugar"`
	Sodium             float64 `protobuf:"fixed64,8,opt,name=sodium" db:"sodium"`
	Potassium          float64 `protobuf:"fixed64,9,opt,name=potassium" db:"potassium"`
	Cholesterol        float64 `protobuf:"fixed64,10,opt,name=cholesterol" db:"cholesterol"`
}
