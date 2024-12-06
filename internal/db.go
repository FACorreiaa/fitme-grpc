package internal

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"fmt"
	"math"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"context"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
	"go.uber.org/zap/zapcore"

	"github.com/pkg/errors"

	config "github.com/FACorreiaa/fitme-grpc/config"
	"github.com/FACorreiaa/fitme-grpc/logger"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

const retries = 25

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

type DatabaseConfig struct {
	ConnectionURL string
}

var (
	// pg
	database = os.Getenv("POSTGRES_DB")
	password = os.Getenv("POSTGRES_PASSWORD")
	host     = os.Getenv("POSTGRES_HOST")
	schema   = os.Getenv("POSTGRES_SCHEMA")

	// redis
	redisPassword = os.Getenv("REDIS_PASSWORD")
)

func NewRedisConfig() (*redis.Client, error) {
	cfg, err := config.InitConfig()

	if err != nil {
		zap.Error(err)
	}
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Repositories.Redis.Host, cfg.Repositories.Redis.Port),
		Password: redisPassword,
		DB:       cfg.Repositories.Redis.DB,
	}), nil
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	if err := logger.Init(zapcore.InfoLevel); err != nil {
		fmt.Println("Error initializing logger:", err)
		os.Exit(1)
	}
	log := logger.Log
	cfg, err := config.InitConfig()
	if err != nil {
		log.Error("Failed loading Postgres config", zap.Error(err))
		log.Fatal("Error initializing config", zap.Error(err))
	}
	//err = godotenv.Load(".env")
	//if err != nil {
	//	log.Error("Error loading .env file", zap.Error(err))
	//	log.Fatal("Failed to load .env file")
	//}

	query := url.Values{
		"sslmode":  []string{"disable"},
		"timezone": []string{"utc"},
	}
	if schema != "" {
		query.Add("search_path", schema)
	}
	connURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.Repositories.Postgres.Username, password),
		Host:     host + ":" + cfg.Repositories.Postgres.Port,
		Path:     cfg.Repositories.Postgres.DB,
		RawQuery: query.Encode(),
	}
	log.Info("Database connection URL generated", zap.String("connectionURL", connURL.String()))

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
	log := logger.Log
	log.Info("Running migrations")
	ctx := context.Background()
	files, err := migrationFS.ReadDir("migrations")
	if err != nil {
		log.Error("Failed executing migrations", zap.Error(err))
		return err
	}

	log.Info("Creating migrations table")
	_, err = conn.Exec(ctx, `
		create table if not exists _migrations (
			name text primary key,
			hash text not null,
			created_at timestamp default now()
		);
	`)
	if err != nil {
		log.Error("Failed loading Postgres config", zap.Error(err))
	}

	log.Info("Checking applied migrations")
	rows, _ := conn.Query(ctx, `select name, hash from _migrations order by created_at desc`)
	var name, hash string
	appliedMigrations := make(map[string]string)
	_, err = pgx.ForEachRow(rows, []any{&name, &hash}, func() error {
		appliedMigrations[name] = hash
		return nil
	})

	if err != nil {
		log.Error("Failed loading Postgres config", zap.Error(err))
		return err
	}

	for _, file := range files {
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})

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

			log.Info(file.Name() + " already applied")
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
			log.Error("Failed loading Postgres config", zap.Error(err))
			return err
		}
		log.Info(file.Name() + " applied")
	}

	log.Info("Migrations finished")
	return nil
}

// PgService Service represents a service that interacts with a database.
type PgService interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type pgService struct {
	db *sql.DB
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *pgService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := logger.Init(zapcore.InfoLevel); err != nil {
		fmt.Println("Error initializing logger:", err)
		os.Exit(1)
	}
	log := logger.Log
	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatal(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *pgService) Close() error {
	if err := logger.Init(zapcore.InfoLevel); err != nil {
		fmt.Println("Error initializing logger:", err)
		os.Exit(1)
	}
	log := logger.Log
	cfg, err := config.InitConfig()
	if err != nil {
		return errors.Wrap(err, "error loading Postgres config")
	}
	database := host + ":" + cfg.Repositories.Postgres.Port
	log.Fatal(fmt.Sprintf("Disconnected from database: %s", database))
	return s.db.Close()
}

// redis check

type RedisService interface {
	Health() map[string]string
}

type redisService struct {
	db *redis.Client
}

// checkRedisHealth checks the health of the Redis server and adds the relevant statistics to the stats map.
func (s *redisService) checkRedisHealth(ctx context.Context, stats map[string]string) map[string]string {
	// Ping the Redis server to check its availability.
	pong, err := s.db.Ping(ctx).Result()
	log := logger.Log

	// Note: By extracting and simplifying like this, `log.Fatalf(fmt.Sprintf("db down: %v", err))`
	// can be changed into a standard error instead of a fatal error.
	if err != nil {
		log.Fatal(fmt.Sprintf("db down: %v", err))
	}

	// Redis is up
	stats["redis_status"] = "up"
	stats["redis_message"] = "It's healthy"
	stats["redis_ping_response"] = pong

	// Retrieve Redis server information.
	info, err := s.db.Info(ctx).Result()
	if err != nil {
		stats["redis_message"] = fmt.Sprintf("Failed to retrieve Redis info: %v", err)
		return stats
	}

	// Parse the Redis info response.
	redisInfo := parseRedisInfo(info)

	// Get the pool stats of the Redis client.
	poolStats := s.db.PoolStats()

	// Prepare the stats map with Redis server information and pool statistics.
	// Note: The "stats" map in the code uses string keys and values,
	// which is suitable for structuring and serializing the data for the frontend (e.g., JSON, XML, HTMX).
	// Using string types allows for easy conversion and compatibility with various data formats,
	// making it convenient to create health stats for monitoring or other purposes.
	// Also note that any raw "memory" (e.g., used_memory) value here is in bytes and can be converted to megabytes or gigabytes as a float64.
	stats["redis_version"] = redisInfo["redis_version"]
	stats["redis_mode"] = redisInfo["redis_mode"]
	stats["redis_connected_clients"] = redisInfo["connected_clients"]
	stats["redis_used_memory"] = redisInfo["used_memory"]
	stats["redis_used_memory_peak"] = redisInfo["used_memory_peak"]
	stats["redis_uptime_in_seconds"] = redisInfo["uptime_in_seconds"]
	stats["redis_hits_connections"] = strconv.FormatUint(uint64(poolStats.Hits), 10)
	stats["redis_misses_connections"] = strconv.FormatUint(uint64(poolStats.Misses), 10)
	stats["redis_timeouts_connections"] = strconv.FormatUint(uint64(poolStats.Timeouts), 10)
	stats["redis_total_connections"] = strconv.FormatUint(uint64(poolStats.TotalConns), 10)
	stats["redis_idle_connections"] = strconv.FormatUint(uint64(poolStats.IdleConns), 10)
	stats["redis_stale_connections"] = strconv.FormatUint(uint64(poolStats.StaleConns), 10)
	stats["redis_max_memory"] = redisInfo["maxmemory"]

	// Calculate the number of active connections.
	// Note: We use math.Max to ensure that activeConns is always non-negative,
	// avoiding the need for an explicit check for negative values.
	// This prevents a potential underflow situation.
	activeConns := uint64(math.Max(float64(poolStats.TotalConns-poolStats.IdleConns), 0))
	stats["redis_active_connections"] = strconv.FormatUint(activeConns, 10)

	// Calculate the pool size percentage.
	poolSize := s.db.Options().PoolSize
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
	poolSizePercentage := float64(connectedClients) / float64(poolSize) * 100
	stats["redis_pool_size_percentage"] = fmt.Sprintf("%.2f%%", poolSizePercentage)

	// Evaluate Redis stats and update the stats map with relevant messages.
	return s.evaluateRedisStats(redisInfo, stats)
}

// evaluateRedisStats evaluates the Redis server statistics and updates the stats map with relevant messages.
func (s *redisService) evaluateRedisStats(redisInfo, stats map[string]string) map[string]string {
	poolSize := s.db.Options().PoolSize
	poolStats := s.db.PoolStats()
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
	highConnectionThreshold := int(float64(poolSize) * 0.8)

	// Check if the number of connected clients is high.
	if connectedClients > highConnectionThreshold {
		stats["redis_message"] = "Redis has a high number of connected clients"
	}

	// Check if the number of stale connections exceeds a threshold.
	minStaleConnectionsThreshold := 500
	if int(poolStats.StaleConns) > minStaleConnectionsThreshold {
		stats["redis_message"] = fmt.Sprintf("Redis has %d stale connections.", poolStats.StaleConns)
	}

	// Check if Redis is using a significant amount of memory.
	usedMemory, _ := strconv.ParseInt(redisInfo["used_memory"], 10, 64)
	maxMemory, _ := strconv.ParseInt(redisInfo["maxmemory"], 10, 64)
	if maxMemory > 0 {
		usedMemoryPercentage := float64(usedMemory) / float64(maxMemory) * 100
		if usedMemoryPercentage >= 90 {
			stats["redis_message"] = "Redis is using a significant amount of memory"
		}
	}

	// Check if Redis has been recently restarted.
	uptimeInSeconds, _ := strconv.ParseInt(redisInfo["uptime_in_seconds"], 10, 64)
	if uptimeInSeconds < 3600 {
		stats["redis_message"] = "Redis has been recently restarted"
	}

	// Check if the number of idle connections is high.
	idleConns := int(poolStats.IdleConns)
	highIdleConnectionThreshold := int(float64(poolSize) * 0.7)
	if idleConns > highIdleConnectionThreshold {
		stats["redis_message"] = "Redis has a high number of idle connections"
	}

	// Check if the connection pool utilization is high.
	poolUtilization := float64(poolStats.TotalConns-poolStats.IdleConns) / float64(poolSize) * 100
	highPoolUtilizationThreshold := 90.0
	if poolUtilization > highPoolUtilizationThreshold {
		stats["redis_message"] = "Redis connection pool utilization is high"
	}

	return stats
}

// parseRedisInfo parses the Redis info response and returns a map of key-value pairs.
func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}
