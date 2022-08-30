package repository

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"os"

	"github.com/XBozorg/bookstore/config"

	"github.com/go-redis/redis/v9"
	"github.com/go-sql-driver/mysql"
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

	initFile, err := os.Open("db/init.sql")
	if err != nil {
		return err
	}
	initBytes, err := io.ReadAll(initFile)
	if err != nil {
		return err
	}
	if _, err = mysqldb.Exec(string(initBytes)); err != nil {
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
