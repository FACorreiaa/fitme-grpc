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
			Port              string `mapstructure:"port"`
			Username          string `mapstructure:"username"`
			DB                string `mapstructure:"db"`
			SSLMODE           string `mapstructure:"SSLMODE"`
			MAXCONWAITINGTIME int    `mapstructure:"MAXCONWAITINGTIME"`
		}
		Redis struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
			Pass string `mapstructure:"pass"`
			DB   int    `mapstructure:"db"`
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
	} `mapstructure:"upstreamServices"`
}

func InitConfig() (Config, error) {
	var config Config
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("config")
	v.AddConfigPath("/app/config")
	v.AddConfigPath("/usr/local/bin")
	v.AddConfigPath("/usr/local/bin/fitme")

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read embedded config: %s", err)
	}

	if err := v.ReadConfig(bytes.NewReader(embeddedConfig)); err != nil {
		return Config{}, fmt.Errorf("failed to read embedded config: %s", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %s", err)
	}
	println("Successfully loaded app configs..")
	return config, nil
}
