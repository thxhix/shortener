package migrations

import (
	"github.com/thxhix/shortener/internal/app/database/interfaces"
)

func Migrate(db interfaces.Database) error {
	query := `
	CREATE TABLE IF NOT EXISTS shortener (
		id SERIAL PRIMARY KEY,
		original VARCHAR(512) UNIQUE NOT NULL,
		shorten VARCHAR(10) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	driver := db.GetDriver()
	if driver == nil {
		return nil
	}
	_, err := driver.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
