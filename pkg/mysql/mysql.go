package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ImamTry257/Notify-Service/config"
)

func NewMySQLConnection(cfg config.MySQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
