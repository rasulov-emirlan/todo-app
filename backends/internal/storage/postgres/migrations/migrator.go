package migrations

import (
	"database/sql"
	"embed"

	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrations embed.FS

func Up(url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	goose.SetBaseFS(migrations)
	goose.SetDialect("postgres")

	return goose.Up(db, ".")
}
