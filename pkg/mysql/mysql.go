package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ImamTry257/Notify-Service/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(cfg config.MySQLConfig) (*sql.DB, error) {
	log.Println("MYSQL_HOST:", os.Getenv("MYSQL_HOST"), cfg.Host)
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
