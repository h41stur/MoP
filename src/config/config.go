package config

import "fmt"

type Connection struct {
	Port               int
	Server             string
	Hostname           string
	DbUser             string
	DbName             string
	DbConnectionString string
}

func Load() Connection {

	var conn Connection

	conn.Port = 8080
	conn.Server = "192.168.2.149"
	conn.Hostname = "http://192.168.2.149:8080"
	conn.DbUser = "root"
	conn.DbName = "MoP"
	conn.DbConnectionString = fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True&loc=Local",
		conn.DbUser,
		conn.DbName,
	)

	return conn
}
