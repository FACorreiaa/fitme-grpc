package config

import (
	"time"

	"github.com/FACorreiaa/fitme-protos/modules/customer"
	"github.com/spf13/viper"
)

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
		Customer string `mapstructure:"customer"`
		Auth     string `mapstructure:"auth"`
	} `mapstructure:"upstreamServices"`
}

func InitConfig() (Config, error) {
	var config Config
	v := viper.New()
	v.AddConfigPath("config")
	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		return config, err
	}
	if err := v.Unmarshal(&config); err != nil {
		return config, err
	}
	return config, nil
}
