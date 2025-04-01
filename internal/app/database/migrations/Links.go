package migrations

import (
	"fmt"
	"github.com/thxhix/shortener/internal/app/database/interfaces"
)

func Migrate(db interfaces.Database) {
	query := `
	CREATE TABLE IF NOT EXISTS shortener (
		id SERIAL PRIMARY KEY,
		original VARCHAR(512) NOT NULL,
		shorten VARCHAR(10) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	driver := db.GetDriver()
	if driver == nil {
		return
	}
	_, err := driver.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		panic("Не удалось сделать миграцию")
	}
}
