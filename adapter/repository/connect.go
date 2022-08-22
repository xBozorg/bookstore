package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/XBozorg/bookstore/config"

	"github.com/go-redis/redis/v9"
	"github.com/go-sql-driver/mysql"
)

type Repo struct {
	MySQL *sql.DB
	Redis *redis.Client
}

func (r *Repo) Close() {
	r.MySQL.Close()
	r.Redis.Close()
}

func (r *Repo) mysqlConnect(conf *config.MySQLConfig) error {
	cfg := mysql.Config{
		User:   conf.User,
		Passwd: conf.Pass,
		Net:    conf.Net,
		Addr:   conf.Address,
		DBName: conf.Name,
	}
	mysqldb, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	pingErr := mysqldb.Ping()
	if pingErr != nil {
		return err
	}

	r.MySQL = mysqldb
	return nil
}

func (r *Repo) redisConnect(conf *config.RedisConfig) error {
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

	r.Redis = client
	return nil
}

func (r *Repo) Connect(conf *config.Config) error {

	err := r.mysqlConnect(conf.GetMySQlConfig())
	if err != nil {
		return err
	}

	err = r.redisConnect(conf.GetRedisConfig())
	if err != nil {
		return err
	}

	return nil
}
