package tests

import (
	"github.com/stretchr/testify/mock"
)

type MockActivityRepository struct {
	mock.Mock
}

//func TestSaveExerciseSession(t *testing.T) {
//	repo := &MockActivityRepository{} // Assuming you've implemented a mock repository
//	service := NewActivityService(repo)
//
//	session := &ExerciseSession{UserID: 1, ActivityID: 1, SessionName: "Morning Run"}
//	err := service.SaveExerciseSession(context.Background(), session)
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	}
//}
