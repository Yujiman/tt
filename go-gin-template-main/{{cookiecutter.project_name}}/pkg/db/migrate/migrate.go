package migrate

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

func SetBaseFS(baseFS fs.FS) {
	goose.SetBaseFS(baseFS)
}

func Up(db *sql.DB, schemaName string, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	goose.SetTableName(fmt.Sprintf("%s.goose_db_version", schemaName))

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}
