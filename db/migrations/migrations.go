package migrations

import (
	_ "database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func MigrateGameDB() error {
	m, err := migrate.New(
		"file://sql",
		"postgres://user:mypassword@localhost:5432/migratetest?sslmode=disable",
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
