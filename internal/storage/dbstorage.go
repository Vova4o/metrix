package storage

import (
	"database/sql"
	"fmt"

	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
)

type DbStorage struct {
	DB *sql.DB
	handlers.Storager
}

func NewDBConnection(flag string) (*DbStorage, error) {
	var err error

	dbStorage := &DbStorage{}

	dbStorage.DB, err = sql.Open("postgres", flag)
	if err != nil {
		err = fmt.Errorf("failed to open database: %v", err)
		logger.Log.WithError(err).Error("Failed to open database")
		return nil, err
	}

	err = dbStorage.DB.Ping()
	if err != nil {
		err = fmt.Errorf("failed to ping database: %v", err)
		logger.Log.WithError(err).Error("Failed to ping database")
		return nil, err
	}

	defer func() {
		err := dbStorage.DB.Close()
		if err != nil {
			logger.Log.Fatalf("Failed to close the database connection: %v", err)
		}
	}()

	err = CreateTables(dbStorage.DB)
	if err != nil {
		return nil, err
	}

	return dbStorage, nil
}

func CreateTables(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS metrics (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        value DOUBLE PRECISION,
		type VARCHAR(30),
        timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
    `

	_, err := db.Exec(query)
	if err != nil {
		err = fmt.Errorf("failed to create table: %v", err)
		logger.Log.WithError(err).Error("Failed to create table")
		return err
	}

	return nil
}
