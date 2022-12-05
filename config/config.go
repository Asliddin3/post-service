package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment         string
	PostgresHost        string
	PostgresPort        int
	PostgresDatabase    string
	PostgresUser        string
	PostgresPassword    string
	LogLevel            string
	RPCPort             string
	ReviewServiceHost   string
	ReviewServicePort   int
	CustomerSerivceHost string
	CustomerSerivcePort int
	KafkaHost           string
	KafkaPort           int
}

func Load() Config {
	c := Config{}
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "database-1.c9lxq3r1itbt.us-east-1.rds.amazonaws.com"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DB", "post"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "asliddin2001"))
	c.ReviewServiceHost = cast.ToString(getOrReturnDefault("REVIEW_SERVISE_HOST", "review-servise"))
	c.ReviewServicePort = cast.ToInt(getOrReturnDefault("REVIEW_SERVISE_PORT", 8840))
	c.CustomerSerivceHost = cast.ToString(getOrReturnDefault("CUSTOMER_SERVISE_HOST", "customer-servis"))
	c.CustomerSerivcePort = cast.ToInt(getOrReturnDefault("CUSTOMER_SERVISE_PORT", 8810))
	c.KafkaHost=cast.ToString(getOrReturnDefault("KAFKA_HOST","kafka"))
	c.KafkaPort=cast.ToInt(getOrReturnDefault("KAFKA_PORT",9092))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))

	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":8820"))
	return c
}

func getOrReturnDefault(key string, defaulValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaulValue
}
