package drivers

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	customErrors "github.com/thxhix/shortener/internal/errors"
	"github.com/thxhix/shortener/internal/models"
	"log"
	"time"
)

// PostgresQLDatabase implements the Database interface using PostgreSQL.
// It stores shortened URLs in the "shortener" table and supports
// migrations, transactions, and user-specific queries.
type PostgresQLDatabase struct {
	driver *sql.DB
}

// NewPQLDatabase creates a new PostgresQLDatabase instance with the given DSN.
// Returns an error if the connection cannot be opened.
func NewPQLDatabase(params string) (*PostgresQLDatabase, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	return &PostgresQLDatabase{
		driver: db,
	}, nil
}

// RunMigrations runs pending database migrations from the "migrations" folder.
// If no changes are required, it succeeds silently.
func (db *PostgresQLDatabase) RunMigrations() error {
	driver, err := postgres.WithInstance(db.driver, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// AddLink inserts a single link into the database.
// If the original URL already exists, it returns the existing shorten hash
// and ErrDuplicate.
func (db *PostgresQLDatabase) AddLink(ctx context.Context, original string, shorten string, userID string) (string, error) {
	var user interface{}
	if userID == "" {
		user = nil
	} else {
		user = userID
	}

	query := `
        INSERT INTO shortener (original, shorten, user_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (original) DO UPDATE
        SET original = EXCLUDED.original
        RETURNING shorten
    `
	var insertedShorten string
	err := db.driver.QueryRowContext(ctx, query, original, shorten, user).Scan(&insertedShorten)
	if err != nil {
		return "", err
	}

	if insertedShorten != shorten {
		return insertedShorten, customErrors.ErrDuplicate
	}

	return insertedShorten, nil
}

// AddLinks inserts multiple links into the database in a transaction.
// If any insert fails, the transaction is rolled back.
func (db *PostgresQLDatabase) AddLinks(ctx context.Context, list models.DBShortenRowList, userID string) (err error) {
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

	var user interface{}
	if userID == "" {
		user = nil
	} else {
		user = userID
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO shortener (original, shorten, user_id) VALUES($1, $2, $3)")

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
		_, err = stmt.ExecContext(ctx, row.URL, row.Hash, user)
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

// GetFullLink retrieves a link by its hash.
// Returns an error if the hash is not found.
func (db *PostgresQLDatabase) GetFullLink(ctx context.Context, hash string) (models.DBShortenRow, error) {
	query := "SELECT id, original, shorten, is_deleted, created_at FROM shortener WHERE (shorten) LIKE ($1)"

	row := db.driver.QueryRowContext(ctx, query, hash)

	var data models.DBShortenRow
	err := row.Scan(
		&data.ID,
		&data.URL,
		&data.Hash,
		&data.IsDeleted,
		&data.Time,
	)
	if err != nil {
		return models.DBShortenRow{}, err
	}

	return data, nil
}

// Close closes the underlying PostgreSQL connection.
func (db *PostgresQLDatabase) Close() error {
	return db.driver.Close()
}

// PingConnection checks whether the PostgreSQL connection is alive.
// Uses a 1-second timeout context.
func (db *PostgresQLDatabase) PingConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return db.driver.PingContext(ctx)
}

// GetDriver returns the underlying *sql.DB driver instance.
func (db *PostgresQLDatabase) GetDriver() *sql.DB { return db.driver }

// GetUserFullLinks retrieves all links created by the given user.
// Returns an empty slice if the user has no links.
func (db *PostgresQLDatabase) GetUserFullLinks(ctx context.Context, userID string) (models.DBShortenRowList, error) {
	if userID == "" {
		return nil, nil
	}

	query := `SELECT id, original, shorten, created_at FROM shortener WHERE user_id = $1`

	rows, err := db.driver.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results models.DBShortenRowList

	for rows.Next() {
		var row models.DBShortenRow
		err := rows.Scan(&row.ID, &row.URL, &row.Hash, &row.Time)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	// проверка на ошибки сканирования после Next
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// RemoveUserLinks marks user links as deleted by setting is_deleted = true.
// Returns an error if the update fails.
func (db *PostgresQLDatabase) RemoveUserLinks(ctx context.Context, userID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE shortener SET is_deleted = true 
	          WHERE user_id = $1 AND shorten = ANY($2)`
	_, err := db.driver.ExecContext(ctx, query, userID, pq.Array(ids))
	return err
}
