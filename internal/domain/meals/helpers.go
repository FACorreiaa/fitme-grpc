package meals

import (
	"database/sql"
	"strconv"

	pbml "github.com/FACorreiaa/fitme-protos/modules/meal/generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func safeString(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	}
	return ""
}

func safeFloat64ToInt32(value interface{}) int32 {
	if v, ok := value.(float64); ok {
		return int32(v)
	}
	return 0
}

func safeFloat64(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	}
	return 0
}

func parseStringToFloat(s, errorMessage string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "%s: %v", errorMessage, err)
	}
	return val, nil
}

func convertFloat(n sql.NullFloat64) float64 {
	if n.Valid {
		return n.Float64
	}
	return 0
}

// Helper function to calculate nutrient totals
func calculateTotals(ingredients []*pbml.XMealIngredient) *pbml.XTotalMealNutrients {
	totals := &pbml.XTotalMealNutrients{}
	for _, ing := range ingredients {
		totals.Calories += ing.Calories
		totals.Protein += ing.Protein
		totals.CarbohydratesTotal += ing.CarbohydratesTotal
		totals.FatTotal += ing.FatTotal
		totals.FatSaturated += ing.FatSaturated
		totals.Fiber += ing.Fiber
		totals.Sugar += ing.Sugar
		totals.Sodium += ing.Sodium
		totals.Potassium += ing.Potassium
		totals.Cholesterol += ing.Cholesterol
	}
	return totals
}
