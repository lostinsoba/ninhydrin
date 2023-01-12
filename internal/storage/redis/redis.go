package redis

import (
	r "github.com/go-redis/redis/v9"

	"lostinsoba/ninhydrin/internal/model"
)

const (
	Kind         = "redis"
	KindSentinel = "redis.sentinel"
)

type Storage struct {
	client r.UniversalClient
}

const (
	settingMasterName = "master"
	settingAddrs      = "addrs"
	settingUsername   = "user"
	settingPassword   = "password"
	settingDatabase   = "database"
)

func NewRedis(settings model.Settings) (*Storage, error) {
	addr, err := settings.ReadStr(settingAddrs)
	if err != nil {
		return nil, err
	}
	username, err := settings.ReadStr(settingUsername)
	if err != nil {
		return nil, err
	}
	password, err := settings.ReadStr(settingPassword)
	if err != nil {
		return nil, err
	}
	database, err := settings.ReadInt(settingDatabase)
	if err != nil {
		return nil, err
	}

	options := &r.UniversalOptions{
		Addrs:    []string{addr},
		Username: username,
		Password: password,
		DB:       database,
	}

	return &Storage{client: r.NewUniversalClient(options)}, nil
}

func NewRedisSentinel(settings model.Settings) (*Storage, error) {
	masterName, err := settings.ReadStr(settingMasterName)
	if err != nil {
		return nil, err
	}
	addrs, err := settings.ReadStrArr(settingAddrs)
	if err != nil {
		return nil, err
	}
	username, err := settings.ReadStr(settingUsername)
	if err != nil {
		return nil, err
	}
	password, err := settings.ReadStr(settingPassword)
	if err != nil {
		return nil, err
	}
	database, err := settings.ReadInt(settingDatabase)
	if err != nil {
		return nil, err
	}

	options := &r.UniversalOptions{
		MasterName: masterName,
		Addrs:      addrs,
		Username:   username,
		Password:   password,
		DB:         database,
	}

	return &Storage{client: r.NewUniversalClient(options)}, nil
}
