package drivers

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	customErrors "github.com/thxhix/shortener/internal/errors"
	"github.com/thxhix/shortener/internal/models"
	"log"
	"time"
)

type PostgresQLDatabase struct {
	driver *sql.DB
}

func (db *PostgresQLDatabase) AddLink(original string, shorten string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO shortener (original, shorten)
        VALUES ($1, $2)
        ON CONFLICT (original) DO UPDATE
        SET original = EXCLUDED.original
        RETURNING shorten
    `
	var insertedShorten string
	err := db.driver.QueryRowContext(ctx, query, original, shorten).Scan(&insertedShorten)
	if err != nil {
		return "", err
	}

	// если insert не произошёл — значит был конфликт, возвращаем ошибку и уже существующий хэш
	if insertedShorten != shorten {
		return insertedShorten, customErrors.ErrDuplicate
	}

	return insertedShorten, nil
}

func (db *PostgresQLDatabase) AddLinks(ctx context.Context, list models.DBShortenRowList) (err error) {
	tx, err := db.driver.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		// Только если была ошибка вне defer
		if err != nil {
			if RBError := tx.Rollback(); RBError != nil {
				log.Printf("ошибка при rollback: %v", RBError)
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO shortener (original, shorten) VALUES($1, $2)")

	if err != nil {
		return err
	}
	defer func() {
		// Подменяем только если основной ошибки не было
		if CErr := stmt.Close(); CErr != nil && err == nil {
			err = CErr
		}
	}()

	for _, row := range list {
		_, err = stmt.ExecContext(ctx, row.URL, row.Hash)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return
}

func (db *PostgresQLDatabase) GetFullLink(hash string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM shortener WHERE (shorten) LIKE ($1)"

	row := db.driver.QueryRowContext(ctx, query, hash)

	var data models.DBShortenRow
	err := row.Scan(
		&data.ID,
		&data.URL,
		&data.Hash,
		&data.Time,
	)
	if err != nil {
		return "", err
	}

	return data.URL, nil
}

func (db *PostgresQLDatabase) Close() error {
	return db.driver.Close()
}

func (db *PostgresQLDatabase) PingConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return db.driver.PingContext(ctx)
}

func (db *PostgresQLDatabase) GetDriver() *sql.DB { return db.driver }

func NewPQLDatabase(params string) (*PostgresQLDatabase, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	return &PostgresQLDatabase{
		driver: db,
	}, nil
}
