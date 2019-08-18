package app

import (
	"database/sql"

	"github.com/triasmoro/inventory-api/storage"
)

type App struct {
	DB    *sql.DB
	Store *storage.Store
}

func NewApp() (*App, error) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, err
	}

	return &App{
		DB:    db,
		Store: storage.NewStore(db),
	}, nil
}
