package postgres

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

func NewPostgres(settings map[string]string) (*Storage, error) {
	connStr, ok := settings[settingConnStr]
	if !ok {
		return nil, fmt.Errorf("%s setting not present", settingMaxOpenConns)
	}
	maxOpenConnsRaw, ok := settings[settingMaxOpenConns]
	if !ok {
		return nil, fmt.Errorf("%s setting not present", settingMaxOpenConns)
	}
	maxOpenConns, err := strconv.Atoi(maxOpenConnsRaw)
	if err != nil {
		return nil, fmt.Errorf("%s value %s parsing failed: %s", settingMaxOpenConns, maxOpenConnsRaw, err)
	}
	connMaxLifetimeRaw, ok := settings[settingConnMaxLifetime]
	if !ok {
		return nil, fmt.Errorf("%s setting not present", settingConnMaxLifetime)
	}
	connMaxLifetime, err := time.ParseDuration(connMaxLifetimeRaw)
	if err != nil {
		return nil, fmt.Errorf("%s value %s parsing failed: %s", settingConnMaxLifetime, connMaxLifetimeRaw, err)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	return &Storage{db: db}, nil
}
