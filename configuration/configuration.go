package configuration

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	envFileName string
}

type EntityConfiguration struct {
	PostgresPort             int
	PostgresHost             string
	PostgresUser             string
	PostgresPassword         string
	PostgresDb               string
	PostgresPathToMigrations string
}

func NewConfiguration(envFileName string) *Configuration {
	return &Configuration{
		envFileName: envFileName,
	}
}

func (cfg *Configuration) ReadConfiguration() *EntityConfiguration {
	err := godotenv.Load(cfg.envFileName)
	if err != nil {
		panic("Error loading the .env file")
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")
	postgresPathToMigrations := os.Getenv("POSTGRES_PATH_TO_MIGRATIONS")

	postgresPortInt, postgresPortIntErr := strconv.Atoi(postgresPort)
	if postgresPortIntErr != nil {
		panic(postgresPortIntErr)
	}

	return &EntityConfiguration{
		PostgresHost:             postgresHost,
		PostgresPort:             postgresPortInt,
		PostgresUser:             postgresUser,
		PostgresPassword:         postgresPassword,
		PostgresDb:               postgresDb,
		PostgresPathToMigrations: postgresPathToMigrations,
	}
}
