package postgres

import (
	"github.com/pressly/goose"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/migrate"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/postgres"
)

const (
	defaultMigrationsPath = "migrations/pgsql"
)

type (
	wrapper struct {
		path string
		db   *postgres.SqlDB
	}
)

func NewWrapper(db *postgres.SqlDB) migrate.Migrate {
	var p = &wrapper{
		path: defaultMigrationsPath,
		db:   db,
	}

	return p
}

func (w *wrapper) SetPath(path string) {
	if path != "" {
		w.path = path
	}
}

func (w *wrapper) Up() (err error) {
	return goose.Up(w.db, w.path)
}

func (w *wrapper) UpByOne() (err error) {
	return goose.UpByOne(w.db, w.path)
}
func (w *wrapper) UpTo(version int64) (err error) {
	return goose.UpTo(w.db, w.path, version)
}

func (w *wrapper) Down() (err error) {
	return goose.Reset(w.db, w.path)
}

func (w *wrapper) DownByOne() (err error) {
	return goose.Down(w.db, w.path)
}

func (w *wrapper) DownTo(version int64) (err error) {
	return goose.DownTo(w.db, w.path, version)
}
