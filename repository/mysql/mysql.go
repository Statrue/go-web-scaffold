package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.Port"),
		viper.GetString("mysql.schema"),
		viper.GetString("mysql.params"),
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}

	db.SetMaxOpenConns(viper.GetInt("mysql.maxConn"))
	db.SetMaxIdleConns(viper.GetInt("mysql.maxIdle"))

	return
}

func Close() {
	if err := db.Close(); err != nil {
		zap.L().Error("MySQL close failed: ", zap.Error(err))
	}
}
