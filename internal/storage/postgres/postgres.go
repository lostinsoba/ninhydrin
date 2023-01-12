package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"lostinsoba/ninhydrin/internal/model"
)

const Kind = "postgres"

type Storage struct {
	db *sqlx.DB
}

const (
	settingConnStr         = "conn_str"
	settingMaxOpenConns    = "max_open_conns"
	settingConnMaxLifetime = "conn_max_lifetime"
)

func NewPostgres(settings model.Settings) (*Storage, error) {
	connStr, err := settings.ReadStr(settingConnStr)
	if err != nil {
		return nil, err
	}
	maxOpenConns, err := settings.ReadInt(settingMaxOpenConns)
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := settings.ReadDuration(settingConnMaxLifetime)
	if err != nil {
		return nil, err
	}
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	return &Storage{db: db}, nil
}
