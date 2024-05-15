package service

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Protocol string `json:"protocol"`
}

var _mysqlClient *sql.DB
var _mysqlConfig DbConfig

func InitMysql(dbConfig DbConfig) error {
	_mysqlConfig = dbConfig

	user := dbConfig.User
	password := dbConfig.Password
	protocol := dbConfig.Protocol
	host := dbConfig.Host
	port := dbConfig.Port
	database := dbConfig.Database

	var err error
	addr := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", user, password, protocol, host, port, database)
	_mysqlClient, err = sql.Open("mysql", addr)
	if err != nil {
		return err
	}

	return nil
}

func LoadDNS(dn string) string {
	stmt, err := _mysqlClient.Prepare("select ip from upperDns where domainname = ?")
	if err != nil {
		return err.Error()
	}

	defer stmt.Close()
	row := stmt.QueryRow(dn)

	var ip string
	err = row.Scan(&ip)

	if err != nil && err == sql.ErrNoRows {
		return "Error2"
	}

	return ip
}

func RegsterDNS(ip string, dn string) string {
	stmt, err := _mysqlClient.Prepare("insert into upperDns (ip, domainname) values(?, ?)")
	if err != nil {
		return err.Error()
	}

	defer stmt.Close()

	result, err := stmt.Exec(string(ip), string(dn))
	if err != nil {
		return "이미 존재하는 값을 입력했는지 확인하십시오."
	}

	_, _ = result.LastInsertId()

	return "Insert Success"
}
