package Utils

import (
	"database/sql"
	_"github.com/lib/pq"
	"log"
	"fmt"
)


var Connect *sql.DB

func ConnectToDB() {
	configs := MainConfig.DBConf
	var err error
	connectInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s", configs.User, configs.Password,
		configs.Dbname, configs.SSLmode, configs.Address)
	Connect, err = sql.Open("postgres", connectInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = Connect.Ping()
	if err != nil {
		log.Fatal(err)
	}
	Connect.SetMaxIdleConns(100)
	Connect.SetMaxOpenConns(100)
}

