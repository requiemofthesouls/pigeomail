package migrate

import (
	"fmt"
	"os"
	"text/template"

	"github.com/pressly/goose"
)

type (
	Migrate interface {
		// Up Migrate the DB to the most recent version available
		Up() (err error)
		// UpByOne Migrate the DB up by 1
		UpByOne() (err error)
		// UpTo Migrate the DB to a specific version
		UpTo(version int64) (err error)
		// Down Roll back all migrations
		Down() (err error)
		// DownByOne Roll back the version by 1
		DownByOne() (err error)
		// DownTo Roll back to a specific version
		DownTo(version int64) (err error)
		// SetPath set path to migrations
		SetPath(path string)
	}
)

func CreateMigrationForPostgres(name, migratePath string) error {
	if name == "" {
		return fmt.Errorf("no name provided")
	}

	var err error
	if err = os.MkdirAll(migratePath, 0755); err != nil {
		return err
	}

	var sqlMigrationTemplate = template.Must(
		template.New("goose.sql-migration").
			Parse(
				`-- +goose Up 
-- SQL in this section is executed when the migration is applied. 
	  
-- +goose Down 
-- SQL in this section is executed when the migration is rolled back. 
`))

	if err = goose.CreateWithTemplate(nil, migratePath, sqlMigrationTemplate, name, "sql"); err != nil {
		return err
	}

	return nil
}
