package tests

import (
	"github.com/stretchr/testify/mock"
)

type MockRepositoryActivity struct {
	mock.Mock
}

//func TestSaveExerciseSession(t *testing.T) {
//	repo := &MockRepositoryActivity{} // Assuming you've implemented a mock repository
//	service := NewServiceActivity(repo)
//
//	session := &ExerciseSession{UserID: 1, ActivityID: 1, SessionName: "Morning Run"}
//	err := service.SaveExerciseSession(context.Background(), session)
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	}
//}
