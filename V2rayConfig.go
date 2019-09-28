package main

import (
	"database/sql"
	"fmt"
	"log"
)

type rayConfig struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  string `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
}

func (s *manageServer) rayConfigs(userId, email string) []*rayConfig {
	stmt, err := s.mysqlDb.Prepare("SELECT server_name FROM v2ray.t_user_server WHERE userid = ?")
	if len(userId) == 0 {
		stmt, err = s.mysqlDb.Prepare("SELECT server_name FROM v2ray.t_user_server JOIN v2ray.t_user ON t_user.userid = t_user_server.userid WHERE useremail = ?")
	}
	if err != nil {
		log.Println(err)
	}
	var rows *sql.Rows
	if len(userId) == 0 {
		rows, err = stmt.Query(email)
	} else {
		rows, err = stmt.Query(userId)
	}
	configs := []*rayConfig{}
	if err != nil {
		log.Println(err)
		return configs
	}
	for rows.Next() {
		var serverName string
		rows.Scan(&serverName)
		c := &rayConfig{
			V:    "2",
			Ps:   serverName,
			Add:  fmt.Sprint(serverName, ".thedanni.design"),
			Port: "443",
			ID:   userId,
			Aid:  "4",
			Net:  "ws",
			Type: "none",
			Host: "",
			Path: "/",
			TLS:  "tls",
		}
		configs = append(configs, c)
	}
	return configs
}
