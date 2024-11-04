package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type SessionManager struct {
	PgPool *pgxpool.Pool
	Redis  *redis.Client
}

func NewSessionManager(pgpool *pgxpool.Pool, redis *redis.Client) *SessionManager {
	return &SessionManager{PgPool: pgpool, Redis: redis}
}

func (s *SessionManager) GenerateSession(userSession UserSession) (string, error) {
	sessionId := uuid.NewString()
	jsonData, _ := json.Marshal(userSession)
	err := s.Redis.Set(context.Background(), sessionId, string(jsonData), 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *SessionManager) SignIn(ctx context.Context, email, password string) (string, error) {
	// check if the user exists
	var user User
	err := s.PgPool.QueryRow(ctx, "select id, username, email, password from users where email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return "", err
	}

	// check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	// create the session
	sessionId := uuid.NewString()
	jsonData, _ := json.Marshal(UserSession{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	err = s.Redis.Set(context.Background(), sessionId, string(jsonData), 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *SessionManager) SignOut(sessionId string) error {
	return s.Redis.Del(context.Background(), sessionId).Err()
}

func (s *SessionManager) GetSession(session string) (*UserSession, error) {
	data, err := s.Redis.Get(context.Background(), session).Result()
	if err != nil {
		return nil, err
	}

	// unmarshal the data
	var userSession UserSession
	err = json.Unmarshal([]byte(data), &userSession)
	if err != nil {
		return nil, err
	}

	return &userSession, nil
}

func SessionMiddleware(sessionManager *SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sessionHeader := r.Header.Get("Authorization")
			if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
				http.Error(w, "Invalid session header", http.StatusUnauthorized)
				return
			}

			sessionID := sessionHeader[7:]

			userSession, err := sessionManager.GetSession(sessionID)
			if err != nil {
				http.Error(w, "Invalid session", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), SessionManagerKey{}, userSession)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
