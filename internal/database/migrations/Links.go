package migrations

import (
	"github.com/thxhix/shortener/internal/database/interfaces"
)

func Migrate(db interfaces.Database) error {
	// Дропать не надо, понятное дело, но для простоты пусть будет так.
	// Надеюсь позднее будет время, разберсь с норм миграциями..
	query := `
	DROP TABLE IF EXISTS shortener;

	CREATE TABLE IF NOT EXISTS shortener (
		id SERIAL PRIMARY KEY,
		original VARCHAR(512) UNIQUE NOT NULL,
		shorten VARCHAR(10) UNIQUE NOT NULL,
		user_id UUID NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP 
	);

	CREATE INDEX IF NOT EXISTS idx_shortener_user_id ON shortener(user_id);`
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
