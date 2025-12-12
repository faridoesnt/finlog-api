package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Init() map[string]string {
	config := make(map[string]string)

	_ = godotenv.Load()

	keys := []string{
		"SERVER_ENV",
		"SERVER_PORT",

		"DB_DIALEG",
		"DB_HOSTWRITER",
		"DB_HOSTREADER",
		"DB_PORT",
		"DB_NAME",

		"USERNAME",
		"PASSWORD",

		"JWT_SECRET",
		"REFRESH_SECRET",
		"JWT_TTL",
		"REFRESH_TTL",

		"RESEND_API_KEY",
		"EMAIL_FROM",
		"RESEND_WEBHOOK_SECRET",

		"MYSQL_ROOT_PASSWORD",
		"MYSQL_DATABASE",
		"MYSQL_USER",
		"MYSQL_PASSWORD",

		"NETDATA_MYSQL_USER",
		"NETDATA_MYSQL_PASSWORD",
	}

	for _, key := range keys {
		val := os.Getenv(key)
		if val == "" {
			panic(fmt.Sprintf("❌ ENV %s is required but not set", key))
		}
		config[key] = val
	}

	return config
}
