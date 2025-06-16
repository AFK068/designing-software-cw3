package migration

import (
	"database/sql"

	"github.com/pressly/goose"
	"go.uber.org/zap"
)

const (
	DialectType = "postgres"
	DriverName  = "postgres"
)

func RunMigrations(dbConn, migrationsPath string, log *zap.Logger) error {
	log.Info("Migrating database started")

	db, err := sql.Open(DriverName, dbConn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect(DialectType); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		return err
	}

	log.Info("Migrations applied successfully")

	return nil
}
