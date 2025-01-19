package meals

import (
	"context"
	"database/sql"
	"sync"
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
	MealID             uuid.UUID `protobuf:"bytes,1,opt,name=meal_id,json=mealId" db:"meal_id"`
	IngredientID       uuid.UUID `protobuf:"bytes,2,opt,name=ingredient_id,json=ingredientId" db:"ingredient_id"`
	Quantity           float64   `protobuf:"fixed64,3,opt,name=quantity" db:"quantity"`
	Calories           float64   `protobuf:"fixed64,4,opt,name=calories" db:"calories"`
	Protein            float64   `protobuf:"fixed64,5,opt,name=protein" db:"protein"`
	FatTotal           float64   `protobuf:"fixed64,6,opt,name=fat_total,json=fatTotal" db:"fat_total"`
	FatSaturated       float64   `protobuf:"fixed64,7,opt,name=fat_saturated,json=fatSaturated" db:"fat_saturated"`
	CarbohydratesTotal float64   `protobuf:"fixed64,8,opt,name=carbohydrates_total,json=carbohydratesTotal" db:"carbohydrates_total"`
	Fiber              float64   `protobuf:"fixed64,9,opt,name=fiber" db:"fiber"`
	Sugar              float64   `protobuf:"fixed64,10,opt,name=sugar" db:"sugar"`
	Sodium             float64   `protobuf:"fixed64,11,opt,name=sodium" db:"sodium"`
	Potassium          float64   `protobuf:"fixed64,12,opt,name=potassium" db:"potassium"`
	Cholesterol        float64   `protobuf:"fixed64,13,opt,name=cholesterol" db:"cholesterol"`
	CreatedAt          time.Time `protobuf:"bytes,4,opt,name=created_at,json=createdAt" db:"created_at"`
	UpdatedAt          time.Time `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt" db:"updated_at"`
}

type TotalNutrients struct {
	Calories           sql.NullFloat64 `protobuf:"fixed64,1,opt,name=calories" db:"calories"`
	Protein            sql.NullFloat64 `protobuf:"fixed64,2,opt,name=protein" db:"protein"`
	FatTotal           sql.NullFloat64 `protobuf:"fixed64,3,opt,name=fat_total,json=fatTotal" db:"fat_total"`
	FatSaturated       sql.NullFloat64 `protobuf:"fixed64,4,opt,name=fat_saturated,json=fatSaturated" db:"fat_saturated"`
	CarbohydratesTotal sql.NullFloat64 `protobuf:"fixed64,5,opt,name=carbohydrates_total,json=carbohydratesTotal" db:"carbohydrates_total"`
	Fiber              sql.NullFloat64 `protobuf:"fixed64,6,opt,name=fiber" db:"fiber"`
	Sugar              sql.NullFloat64 `protobuf:"fixed64,7,opt,name=sugar" db:"sugar"`
	Sodium             sql.NullFloat64 `protobuf:"fixed64,8,opt,name=sodium" db:"sodium"`
	Potassium          sql.NullFloat64 `protobuf:"fixed64,9,opt,name=potassium" db:"potassium"`
	Cholesterol        sql.NullFloat64 `protobuf:"fixed64,10,opt,name=cholesterol" db:"cholesterol"`
}

//type MealPlan struct {
//	ID              uuid.UUID `db:"id"`
//	UserID          uuid.UUID `db:"user_id"`
//	MealPlanID      uuid.UUID `db:"meal_plan_id"`
//	MealPlanNumber  int       `db:"meal_plan_number"`
//	MealDescription string    `db:"description"`
//	Meals           []Meal    `db:"meals"`
//	TotalNutrients  TotalNutrients
//	Notes           string       `db:"notes"`
//	Objective       string       `db:"objective"`
//	CreatedAt       time.Time    `db:"created_at"`
//	UpdatedAt       sql.NullTime `db:"updated_at"`
//	Name            string       `db:"name"`
//}

type MealPlan struct {
	ID           uuid.UUID       `protobuf:"bytes,1,opt,name=meal_plan_id,json=meal_plan_id" db:"id"`
	UserID       uuid.UUID       `db:"user_id"`
	Name         sql.NullString  `db:"name"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    sql.NullTime    `db:"updated_at"`
	Meals        []Meal          `db:"-"` // Exclude from DB scan, populate later
	TotalMacros  *TotalNutrients `protobuf:"bytes,8,rep,name=total_macros,json=totalMacros" db:"total_macros"`
	Description  sql.NullString  `db:"description"`
	Notes        sql.NullString  `db:"notes"`
	Rating       sql.NullFloat64 `db:"rating"`
	Objective    sql.NullString  `db:"objective"`
	Activity     sql.NullString  `db:"activity"`
	Gender       sql.NullString  `db:"gender"`
	QuantityUnit sql.NullString  `db:"quantity_unit"`
}

// to use later

type Broadcaster struct {
	mu        sync.Mutex
	consumers map[chan<- any]struct{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		consumers: make(map[chan<- any]struct{}),
	}
}

func (mb *Broadcaster) Send(msg any) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	for listener := range mb.consumers {
		listener <- msg
	}
}

func (mb *Broadcaster) Close() {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	for listener := range mb.consumers {
		close(listener)
		delete(mb.consumers, listener)
	}
}

func Listen[T any](ctx context.Context, mb *Broadcaster) <-chan T {
	all := make(chan any)
	mb.mu.Lock()
	mb.consumers[all] = struct{}{}
	mb.mu.Unlock()
	go func() {
		<-ctx.Done()
		mb.mu.Lock()
		delete(mb.consumers, all)
		mb.mu.Unlock()
	}()
	ch := make(chan T)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-all:
				if !ok {
					return
				}
				if m, ok := msg.(T); ok {
					select {
					case ch <- m:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return ch
}
