package repository

import (
	"database/sql"

	"github.com/XBozorg/bookstore/config"
	"github.com/go-sql-driver/mysql"
)

type MySQLRepo struct {
	db *sql.DB
}

func (m *MySQLRepo) SetDB(db *sql.DB) {
	m.db = db
}

func Connect(conf *config.MySQLConfig) (MySQLRepo, error) {

	cfg := mysql.Config{
		User:   conf.User,
		Passwd: conf.Pass,
		Net:    conf.Net,
		Addr:   conf.Address,
		DBName: conf.Name,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return MySQLRepo{}, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return MySQLRepo{}, err
	}

	var repo MySQLRepo
	repo.SetDB(db)

	return repo, nil
}
