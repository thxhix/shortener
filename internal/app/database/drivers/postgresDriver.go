package drivers

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	customErrors "github.com/thxhix/shortener/internal/app/errors"
	"github.com/thxhix/shortener/internal/app/models"
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
        ON CONFLICT (original) DO NOTHING
        RETURNING shorten
    `

	var insertedShorten string
	err := db.driver.QueryRowContext(ctx, query, original, shorten).Scan(&insertedShorten)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = db.driver.QueryRowContext(ctx,
				"SELECT shorten FROM shortener WHERE original = $1", original,
			).Scan(&insertedShorten)
			if err != nil {
				return "", err
			}
			return insertedShorten, customErrors.NewDuplicateError()
		}
		return "", err
	}

	return insertedShorten, nil
}

func (db *PostgresQLDatabase) AddLinks(ctx context.Context, list models.DatabaseRowList) error {
	tx, err := db.driver.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO shortener (original, shorten) VALUES($1, $2)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range list {
		_, err := stmt.ExecContext(ctx, row.URL, row.Hash)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresQLDatabase) GetFullLink(hash string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM shortener WHERE (shorten) LIKE ($1)"

	row := db.driver.QueryRowContext(ctx, query, hash)

	var data models.DatabaseRow
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
