package main

import (
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

func (s *manageServer) rayConfigs(userId string) []*rayConfig {
	rows, err := s.mysqlDb.Query("SELECT server_name FROM v2ray.t_user_server WHERE userid = ?", userId)
	if err != nil {
		log.Println(err)
		return []*rayConfig{}
	}
	configs := []*rayConfig{}
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
