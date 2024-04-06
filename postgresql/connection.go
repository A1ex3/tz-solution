package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func (pgsql *Postgresql) Connect(ctx context.Context) (*pgx.Conn, error) {
	// Connection string "postgres://<username>:<password>@<host>:<port>/<database_name>"
	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		pgsql.User,
		pgsql.Password,
		pgsql.Host,
		pgsql.Port,
		pgsql.Database,
	))
	if err != nil {
		logrus.Infof("Unable to connect to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}

func (pgsql *Postgresql) Close(ctx context.Context, conn *pgx.Conn) {
	conn.Close(ctx)
}

func (pgsql *Postgresql) Migrate(migrationDir string) error {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			pgsql.User,
			pgsql.Password,
			pgsql.Host,
			pgsql.Port,
			pgsql.Database,
		),
	)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationDir,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
