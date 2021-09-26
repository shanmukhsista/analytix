package db

import (
	"database/sql"
	"log"
)


type ClickhouseConnectionConfig struct {
	Host string `envconfig:"CH_DB_HOST"`
	TcpPort string `envconfig:"CH_DB_TCP_PORT"`
	PingOnConnect bool `envconfig:"CH_DB_PING_ON_CONNECT"`
	Database string `envconfig:"CH_DB_DATABASE_NAME"`
	User string `envconfig:"CH_DB_USERNAME"`
	Password string `envconfig:"CH_DB_PASSWORD"`
}

func NewClickhouseConnection() ( *sql.DB, error){
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {

		return nil,err
	}
	return connect, nil
}