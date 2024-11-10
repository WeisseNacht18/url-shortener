package databasestorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseStorage struct {
	database *sql.DB
}

func NewDatabaseStorage(dsn string) *DatabaseStorage {
	database, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	const query = `CREATE TABLE IF NOT EXISTS public.url
			(
				id integer NOT NULL GENERATED BY DEFAULT AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
				short_url character varying COLLATE pg_catalog."default" NOT NULL,
				original_url character varying COLLATE pg_catalog."default" NOT NULL,
				user_id character varying COLLATE pg_catalog."default",
				CONSTRAINT url_pk PRIMARY KEY (id)
			)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = database.ExecContext(ctx, query)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	storage := DatabaseStorage{
		database: database,
	}

	return &storage
}

func (storage *DatabaseStorage) AddURL(userID string, originalURL string, shortURL string) (ok bool) {
	err := storage.SaveURLToDatabase(userID, shortURL, originalURL)
	if err == nil {
		ok = true
	} else {
		ok = false
	}
	return
}

func (storage *DatabaseStorage) GetURL(userID string, shortURL string) (originalURL string, ok bool) {
	originalURL, ok = storage.GetURLFromDatabase(userID, shortURL)
	return
}

func (storage *DatabaseStorage) GetAllURLs(userID string) map[string]string {
	return storage.GetAllURLsFromDatabase(userID)
}

func (storage *DatabaseStorage) CheckStorage() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := storage.database.PingContext(ctx)
	logger.Logger.Infoln(err)

	return err
}

func (storage *DatabaseStorage) CheckURL(userID string, originalURL string) (shortLink string, hasURL bool) {
	hasURL = false

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var row *sql.Row

	if userID != "" {
		row = storage.database.QueryRowContext(ctx, "SELECT short_url FROM url WHERE original_url = $1 AND user_id = $2 LIMIT 1", originalURL, userID)
	} else {
		row = storage.database.QueryRowContext(ctx, "SELECT short_url FROM url WHERE original_url = $1 LIMIT 1", originalURL)
	}

	err := row.Scan(&shortLink)

	if err != nil {
		logger.Logger.Infoln(err)
		return
	}

	hasURL = true

	return
}

func (storage *DatabaseStorage) GetUsers() map[string]string {

	return map[string]string{}
}

func (storage *DatabaseStorage) Close() {
	storage.database.Close()
}

func (storage *DatabaseStorage) SaveURLToDatabase(userID string, shortURL string, originalURL string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	logger.Logger.Infoln(userID)

	if userID != "" {
		const query = "INSERT INTO url (short_url, original_url, user_id) VALUES ($1, $2, $3)"
		_, err = storage.database.ExecContext(ctx, query, shortURL, originalURL, userID)
	} else {
		const query = "INSERT INTO url (short_url, original_url) VALUES ($1, $2)"
		_, err = storage.database.ExecContext(ctx, query, shortURL, originalURL)
	}

	return
}

func (storage *DatabaseStorage) GetURLFromDatabase(userID string, shortURL string) (result string, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var originalURL string

	var row *sql.Row

	if userID != "" {
		const query = "SELECT original_url FROM url WHERE short_url = $1 AND user_id = $2 LIMIT 1"
		row = storage.database.QueryRowContext(ctx, query, shortURL, userID)
	} else {
		const query = "SELECT original_url FROM url WHERE short_url = $1 LIMIT 1"
		row = storage.database.QueryRowContext(ctx, query, shortURL)
	}

	err := row.Scan(&originalURL)
	if err != nil {
		return originalURL, false
	}

	return originalURL, true
}

func (storage *DatabaseStorage) GetAllURLsFromDatabase(userID string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result := map[string]string{}

	type Row struct {
		originalURL string `db:"original_url"`
		shortURL    string `db:"short_url"`
	}

	rows, err := storage.database.QueryContext(ctx, "SELECT original_url, short_url FROM url WHERE user_id = $1", userID)

	if err != nil {
		return result
	}

	//тут перебрать каждую строку запроса и засунуть внутрь результата как показано ниже
	for rows.Next() {
		var row Row
		err = rows.Scan(&row.originalURL, &row.shortURL)
		if err != nil {
			return result
		}

		result[row.shortURL] = row.originalURL
	}

	return result
}
