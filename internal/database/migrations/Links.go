package migrations

import (
	"github.com/thxhix/shortener/internal/database/interfaces"
)

func Migrate(db interfaces.Database) error {
	query := `
		CREATE TABLE IF NOT EXISTS shortener (
			id SERIAL PRIMARY KEY,
			original VARCHAR(512) UNIQUE NOT NULL,
			shorten VARCHAR(10) UNIQUE NOT NULL,
			user_id UUID,
			is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
		
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM information_schema.columns 
				WHERE table_name='shortener' AND column_name='user_id'
			) THEN
				ALTER TABLE shortener ADD COLUMN user_id UUID NULL;
			END IF;
		END$$;

		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM information_schema.columns 
				WHERE table_name='shortener' AND column_name='is_deleted'
			) THEN
				ALTER TABLE shortener ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;
			END IF;
		END$$;
		
		CREATE INDEX IF NOT EXISTS idx_shortener_user_id ON shortener(user_id);
		CREATE INDEX IF NOT EXISTS idx_shortener_is_deleted ON shortener(is_deleted);
	`
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
