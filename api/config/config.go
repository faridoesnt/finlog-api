package config

import (
	"github.com/joho/godotenv"
	zero "github.com/rs/zerolog/log"
)

func Init() map[string]string {
	var err error

	config, err := godotenv.Read()
	if err != nil {
		zero.Panic().
			Str("Context", "Load Env").
			Err(err)
	}

	return config
}
