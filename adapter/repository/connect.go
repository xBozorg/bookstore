package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/XBozorg/bookstore/config"
	"github.com/golang-migrate/migrate/v4"

	"github.com/go-redis/redis/v9"
	"github.com/go-sql-driver/mysql"
	mdb "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Storage struct {
	MySQL *sql.DB
	Redis *redis.Client
}

func (s *Storage) Close() {
	s.MySQL.Close()
	s.Redis.Close()
}

func (s *Storage) mysqlConnect(conf *config.MySQLConfig) error {
	cfg := mysql.Config{
		User:                 conf.User,
		Passwd:               conf.Pass,
		Net:                  conf.Net,
		Addr:                 conf.Address,
		DBName:               conf.Name,
		MultiStatements:      true,
		AllowNativePasswords: true,
	}
	mysqldb, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	pingErr := mysqldb.Ping()
	if pingErr != nil {
		return err
	}

	mysqlDriver, err := mdb.WithInstance(mysqldb, &mdb.Config{})
	if err != nil {
		return err
	}

	migrate, err := migrate.NewWithDatabaseInstance(
		"file:///app/db/migrations",
		"mysql",
		mysqlDriver,
	)
	if err != nil {
		return err
	}
	if err := migrate.Up(); err != nil {
		return err
	}

	s.MySQL = mysqldb
	return nil
}

func (s *Storage) redisConnect(conf *config.RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Pass,
		DB:       conf.DB,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	if pong != "PONG" {
		return errors.New("redis connection - pong error")
	}

	s.Redis = client
	return nil
}

func (s *Storage) Connect(conf *config.Config) error {

	err := s.mysqlConnect(conf.GetMySQlConfig())
	if err != nil {
		return err
	}

	err = s.redisConnect(conf.GetRedisConfig())
	if err != nil {
		return err
	}

	return nil
}
