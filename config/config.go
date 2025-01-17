package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"time"

	"github.com/FACorreiaa/fitme-protos/modules/customer"
	_ "github.com/joho/godotenv"
	"github.com/spf13/viper"
	_ "go.uber.org/zap"
)

//go:embed config.yml
var embeddedConfig []byte

type Config struct {
	Mode     string `mapstructure:"mode"`
	Dotenv   string `mapstructure:"dotenv"`
	Handlers struct {
		ExternalAPI struct {
			Port      string `mapstrucutre:"port"`
			CertFile  string `mapstructure:"certFile"`
			KeyFile   string `mapstructure:"keyFile"`
			EnableTLS bool   `mapstracture:"enableTLS"`
		} `mapstructure:"externalAPI"`
		Pprof struct {
			Port      string `mapstructure:"port"`
			CertFile  string `mapstructure:"certFile"`
			KeyFile   string `mapstructure:"keyFile"`
			EnableTLS bool   `mapstructure:"enableTLS"`
		}
		Prometheus struct {
			Port      string `mapstructure:"port"`
			CertFile  string `mapstructure:"certFile"`
			KeyFile   string `mapstructure:"keyFile"`
			EnableTLS bool   `mapstructure:"enableTLS"`
		}
	} `mapstructure:"handlers"`
	Repositories struct {
		Postgres struct {
			Host              string `mapstructure:"host"`
			Password          string `mapstructure:"password"`
			Port              string `mapstructure:"port"`
			Username          string `mapstructure:"username"`
			DB                string `mapstructure:"db"`
			SSLMODE           string `mapstructure:"SSLMODE"`
			MAXCONWAITINGTIME int    `mapstructure:"MAXCONWAITINGTIME"`
		}
		Redis struct {
			Host string        `mapstructure:"host"`
			Port string        `mapstructure:"port"`
			Pass string        `mapstructure:"pass"`
			DB   int           `mapstructure:"db"`
			TTL  time.Duration `mapstructure:"ttl"`
		}
	}
	Server struct {
		HTTPPort       string        `mapstructure:"HTTPPort"`
		GrpcPort       string        `mapstructure:"GRPCPort"`
		Timeout        time.Duration `mapstructure:"HTTPTimeout"`
		CustomerBroker *customer.Broker
	} `mapstructure:"server"`
	UpstreamServices struct {
		Customer    string `mapstructure:"customer"`
		Auth        string `mapstructure:"auth"`
		Calculator  string `mapstructure:"calculator"`
		Activity    string `mapstructure:"activity"`
		Workout     string `mapstructure:"workout"`
		Measurement string `mapstructure:"measurement"`
		Ingredients string `mapstructure:"ingredients"`
		Meals       string `mapstructure:"meals"`
	} `mapstructure:"upstreamServices"`
}

func InitConfig() (Config, error) {
	var config Config
	v := viper.New()

	// Add file-based config paths
	v.AddConfigPath(".")
	v.AddConfigPath("config")
	v.AddConfigPath("/app/config")
	v.AddConfigPath("/usr/local/bin")
	v.AddConfigPath("/usr/local/bin/fitme")

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Try to load file-based config
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to find file-based config: %s. Falling back to embedded config.\n", err)

		// Fallback to embedded config
		if err = v.ReadConfig(bytes.NewReader(embeddedConfig)); err != nil {
			return Config{}, fmt.Errorf("failed to read embedded config: %s", err)
		}
	}

	// Unmarshal the config into the Config struct
	if err = v.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %s", err)
	}
	fmt.Println("Successfully loaded app configs...")
	return config, nil
}
