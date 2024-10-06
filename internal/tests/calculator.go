package tests

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	pbc "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/calculator"
)

// Mock repository to simulate database interactions
type MockCalculatorRepository struct {
	mock.Mock
}

func (m *MockCalculatorRepository) CreateUserMacro(ctx context.Context, macro *pbc.UserMacroDistribution) (*pbc.UserMacroDistribution, error) {
	args := m.Called(ctx, macro)
	return args.Get(0).(*pbc.UserMacroDistribution), args.Error(1)
}

func (m *MockCalculatorRepository) GetUsersMacros(ctx context.Context, req *pbc.GetAllUserMacrosRequest) (*pbc.GetAllUserMacrosResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pbc.GetAllUserMacrosResponse), args.Error(1)
}

func (m *MockCalculatorRepository) GetUserMacros(ctx context.Context, req *pbc.GetUserMacroRequest) (*pbc.GetUserMacroResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pbc.GetUserMacroResponse), args.Error(1)
}

func (m *MockCalculatorRepository) DeleteUserMacro(ctx context.Context, req *pbc.DeleteUserMacroRequest) (*pbc.DeleteUserMacroResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pbc.DeleteUserMacroResponse), args.Error(1)
}

func TestCreateUserMacro(t *testing.T) {
	mockRepo := new(MockCalculatorRepository)
	calculatorService := calculator.NewCalculatorService(mockRepo)

	ctx := context.Background()

	req := &pbc.CreateUserMacroRequest{
		UserMacro: &pbc.UserMacroDistribution{
			Age:                  25,
			Height:               180,
			Weight:               75,
			Gender:               "Male",
			System:               "Metric",
			Activity:             "HIGH",
			Objective:            "BULKING",
			CaloriesDistribution: "BALANCED",
		},
	}

	expectedResponse := &pbc.UserMacroDistribution{
		Age:                  25,
		Height:               180,
		Weight:               75,
		Gender:               "Male",
		System:               "Metric",
		Activity:             "HIGH",
		Objective:            "BULKING",
		CaloriesDistribution: "BALANCED",
	}
	mockRepo.On("CreateUserMacro", ctx, mock.AnythingOfType("*pbc.UserMacroDistribution")).Return(expectedResponse, nil)

	res, err := calculatorService.CreateUserMacro(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedResponse.Age, res.UserMacro.Age)
	assert.Equal(t, expectedResponse.Height, res.UserMacro.Height)
	assert.Equal(t, expectedResponse.Weight, res.UserMacro.Weight)
}

func TestGetUsersMacros(t *testing.T) {
	mockRepo := new(MockCalculatorRepository)
	calculatorService := calculator.NewCalculatorService(mockRepo)
	var req *pbc.GetAllUserMacrosRequest
	ctx := context.Background()

	expectedResponse := &pbc.GetAllUserMacrosResponse{
		UserMacros: []*pbc.UserMacroDistribution{
			{
				Age:                  25,
				Height:               180,
				Weight:               75,
				Gender:               "Male",
				System:               "Metric",
				Activity:             "HIGH",
				Objective:            "BULKING",
				CaloriesDistribution: "BALANCED",
			},
		},
	}

	mockRepo.On("GetUsersMacros", ctx).Return(expectedResponse, nil)

	res, err := calculatorService.GetUsersMacros(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, len(expectedResponse.UserMacros), len(res.UserMacros))
}

func TestGetUserMacros(t *testing.T) {
	mockRepo := new(MockCalculatorRepository)
	calculatorService := calculator.NewCalculatorService(mockRepo)

	ctx := context.Background()

	// Create a mock request
	mockRequest := &pbc.GetUserMacroRequest{
		PlanId: "123",
	}

	// Expected response
	expectedResponse := &pbc.GetUserMacroResponse{
		UserMacro: &pbc.UserMacroDistribution{
			Id:                   "123",
			UserId:               "12345",
			Age:                  25,
			Height:               180,
			Weight:               75,
			Gender:               "Male",
			System:               "Metric",
			Activity:             "HIGH",
			Objective:            "BULKING",
			CaloriesDistribution: "BALANCED",
			CreatedAt:            timestamppb.New(time.Now()), // Set during execution
		},
	}

	// Set up mock response for the repository
	mockRepo.On("GetUserMacros", ctx, mockRequest).Return(expectedResponse, nil)

	// Call the service method
	res, err := calculatorService.GetUserMacros(ctx, mockRequest)

	// Ensure no errors
	assert.NoError(t, err)

	// Ensure the response is not nil
	assert.NotNil(t, res)

	// Compare the expected response and actual response
	assert.Equal(t, expectedResponse.UserMacro.UserId, res.UserMacro.UserId)
	assert.Equal(t, expectedResponse.UserMacro.Age, res.UserMacro.Age)
	assert.Equal(t, expectedResponse.UserMacro.Height, res.UserMacro.Height)
	assert.Equal(t, expectedResponse.UserMacro.Weight, res.UserMacro.Weight)
	assert.Equal(t, expectedResponse.UserMacro.Gender, res.UserMacro.Gender)
	assert.Equal(t, expectedResponse.UserMacro.System, res.UserMacro.System)
	assert.Equal(t, expectedResponse.UserMacro.Activity, res.UserMacro.Activity)
	assert.Equal(t, expectedResponse.UserMacro.Objective, res.UserMacro.Objective)
	assert.Equal(t, expectedResponse.UserMacro.CaloriesDistribution, res.UserMacro.CaloriesDistribution)
	assert.WithinDuration(t, expectedResponse.UserMacro.CreatedAt.AsTime(), res.UserMacro.CreatedAt.AsTime(), time.Second)
}

func TestDeleteUserMacro(t *testing.T) {
	mockRepo := new(MockCalculatorRepository)
	calculatorService := calculator.NewCalculatorService(mockRepo)

	ctx := context.Background()

	req := &pbc.DeleteUserMacroRequest{
		MacroId: "test-user-id",
	}

	expectedResponse := &pbc.DeleteUserMacroResponse{}

	mockRepo.On("DeleteUserMacro", ctx, req).Return(expectedResponse, nil)

	res, err := calculatorService.DeleteUserMacro(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	// assert.True(t, res.Success)
}

func TestCalculateBMR(t *testing.T) {
	tests := []struct {
		name     string
		userData calculator.UserData
		system   calculator.System
		expected float64
		err      error
	}{
		{
			name: "Valid Metric Male",
			userData: calculator.UserData{
				Weight: 70,  // kg
				Height: 175, // cm
				Age:    25,
				Gender: "m",
			},
			system:   calculator.Metric,
			expected: 1715.00, // Adjust the expected value based on calculation
			err:      nil,
		},
		{
			name: "Valid Metric Female",
			userData: calculator.UserData{
				Weight: 60,  // kg
				Height: 165, // cm
				Age:    30,
				Gender: "f",
			},
			system:   calculator.Metric,
			expected: 1393.75, // Adjust the expected value based on calculation
			err:      nil,
		},
		{
			name: "Valid Imperial Male",
			userData: calculator.UserData{
				Weight: 154, // lbs
				Height: 68,  // inches
				Age:    25,
				Gender: "m",
			},
			system:   calculator.Imperial,
			expected: 1756.36, // Adjust the expected value based on calculation
			err:      nil,
		},
		{
			name: "Valid Imperial Female",
			userData: calculator.UserData{
				Weight: 132, // lbs
				Height: 64,  // inches
				Age:    30,
				Gender: "f",
			},
			system:   calculator.Imperial,
			expected: 1262.64, // Adjust the expected value based on calculation
			err:      nil,
		},
		{
			name: "Invalid Weight",
			userData: calculator.UserData{
				Weight: -70,
				Height: 175,
				Age:    25,
				Gender: "m",
			},
			system:   calculator.Metric,
			expected: 0,
			err:      errors.New("weight, height, and age must be positive values"),
		},
		{
			name: "Invalid Gender",
			userData: calculator.UserData{
				Weight: 70,
				Height: 175,
				Age:    25,
				Gender: "x", // Invalid gender
			},
			system:   calculator.Metric,
			expected: 0,
			err:      errors.New("gender must be 'm' or 'f'"),
		},
		{
			name: "Zero Age",
			userData: calculator.UserData{
				Weight: 70,
				Height: 175,
				Age:    0, // Invalid age
				Gender: "m",
			},
			system:   calculator.Metric,
			expected: 0,
			err:      errors.New("weight, height, and age must be positive values"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bmr, err := calculator.CalculateBMR(tt.userData, tt.system)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
			if bmr != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, bmr)
			}
		})
	}
}
