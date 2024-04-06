// @title Car API
// @version 1.0

package main

import (
	"flag"
	"os"
	"time"
	"tzsolution/configuration"
	"tzsolution/http"
	"tzsolution/postgresql"

	"github.com/sirupsen/logrus"
)

func main() {
	envFile := flag.String("ENV_FILE", "", "")
	flag.Parse()
	if *envFile == "" {
		logrus.Infoln("The path to the .env file is not specified.")
	}

	newConfig := configuration.NewConfiguration(*envFile)
	config := newConfig.ReadConfiguration()

	logrus.SetLevel(logrus.InfoLevel | logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	db := postgresql.NewPostgresql(
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresDb,
		config.PostgresUser,
		config.PostgresPassword,
	)
	if err := db.Migrate(config.PostgresPathToMigrations); err != nil {
		panic(err)
	}

	server := http.NewHttpServer(newConfig, db)
	server.Run(":8080")
}
