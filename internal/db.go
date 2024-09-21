package internal

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	config "github.com/FACorreiaa/fitme-grpc/.config"
	"github.com/FACorreiaa/fitme-grpc/logger"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

const retries = 25

func GetEnv(key, defaultVal string) string {
	key = strings.ToUpper(key)
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

type DatabaseConfig struct {
	ConnectionURL string
}

func NewRedisConfig() (*redis.Client, error) {
	cfg, err := config.InitConfig()
	pass := GetEnv("REDIS_PASSWORD", "")

	if err != nil {
		zap.Error(err)
	}
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Repositories.Redis.Host,
		Password: pass,
		DB:       cfg.Repositories.Redis.DB,
	}), nil
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Println(err)
		log.Fatal("Failed loading Postgres config")
	}
	err = godotenv.Load(".env")
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}

	pass := GetEnv("DB_PASS", "")
	schema := GetEnv("", "")

	query := url.Values{
		"sslmode":  []string{"disable"},
		"timezone": []string{"utc"},
	}
	if schema != "" {
		query.Add("search_path", schema)
	}
	connURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.Repositories.Postgres.Username, pass),
		Host:     cfg.Repositories.Postgres.Host + ":" + cfg.Repositories.Postgres.Port,
		Path:     cfg.Repositories.Postgres.DB,
		RawQuery: query.Encode(),
	}
	fmt.Printf("#%v", connURL)
	return &DatabaseConfig{
		ConnectionURL: connURL.String(),
	}, nil
}

func Init(connectionURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connectionURL)
	if err != nil {
		return nil, err
	}
	cfg.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		uuid.Register(conn.TypeMap())
		return nil
	}

	return pgxpool.NewWithConfig(context.Background(), cfg)
}

func WaitForDB(pgpool *pgxpool.Pool) {
	ctx := context.Background()

	for attempts := 1; ; attempts++ {
		if attempts > retries {
			break
		}

		if err := pgpool.Ping(ctx); err == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}
}

func hashVal(contents []byte) string {
	hashed := sha256.New()
	hashed.Write(contents)
	return fmt.Sprintf("%x", hashed.Sum(nil))
}
func Migrate(conn *pgxpool.Pool) error {

	// migrate db
	l := logger.Log
	l.Info("Running migrations")
	ctx := context.Background()
	files, err := migrationFS.ReadDir("migrations")
	if err != nil {
		zap.Error(err)
		return err
	}

	l.Info("Creating migrations table")
	_, err = conn.Exec(ctx, `
		create table if not exists _migrations (
			name text primary key,
			hash text not null,
			created_at timestamp default now()
		);
	`)
	if err != nil {
		zap.Error(err)
		return err
	}

	l.Info("Checking applied migrations")
	rows, _ := conn.Query(ctx, `select name, hash from _migrations order by created_at desc`)
	var name, hash string
	appliedMigrations := make(map[string]string)
	_, err = pgx.ForEachRow(rows, []any{&name, &hash}, func() error {
		appliedMigrations[name] = hash
		return nil
	})

	if err != nil {
		zap.Error(err)
		return err
	}

	for _, file := range files {
		contents, err := migrationFS.ReadFile("migrations/" + file.Name())
		if err != nil {
			return err
		}
		val := hashVal(contents)
		contentHash := fmt.Sprintf("%x", val)

		if prevHash, ok := appliedMigrations[file.Name()]; ok {
			if prevHash != contentHash {
				return errors.New("hash mismatch for")
			}

			l.Info(file.Name() + " already applied")
			continue
		}

		err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
			if _, err = tx.Exec(ctx, string(contents)); err != nil {
				return err
			}

			if _, err := tx.Exec(ctx, `insert into _migrations (name, hash) values ($1, $2)`,
				file.Name(), contentHash); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			zap.Error(err)
			return err
		}
		l.Info(file.Name() + " applied")
	}

	l.Info("Migrations finished")
	return nil
}
