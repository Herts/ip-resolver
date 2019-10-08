package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

type rayConfig struct {
	V     string `json:"v"`
	Ps    string `json:"ps"`
	Add   string `json:"add"`
	Port  string `json:"port"`
	ID    string `json:"id"`
	Aid   string `json:"aid"`
	Net   string `json:"net"`
	Type  string `json:"type"`
	Host  string `json:"host"`
	Path  string `json:"path"`
	TLS   string `json:"tls"`
	Group string `json:"-"`
}

func (s *manageServer) rayConfigs(userId, email string) []*rayConfig {
	stmt, err := s.mysqlDb.Prepare("SELECT server_name, userid, server_region FROM v_user_server WHERE userid = ?")
	if len(userId) == 0 {
		stmt, err = s.mysqlDb.Prepare("SELECT server_name, userid, server_region FROM v_user_server WHERE useremail = ?")
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
		var serverName, id, serverRegion string
		rows.Scan(&serverName, &id, &serverRegion)
		c := &rayConfig{
			V:     "2",
			Ps:    serverName,
			Add:   fmt.Sprint(serverName, ".thedanni.design"),
			Port:  "443",
			ID:    id,
			Aid:   "4",
			Net:   "ws",
			Type:  "none",
			Host:  "",
			Path:  "/",
			TLS:   "tls",
			Group: serverRegion,
		}
		configs = append(configs, c)
	}
	return configs
}

func ConfigToLinks(configs []*rayConfig) string {
	links := ""
	for _, config := range configs {
		byteConfig, err := json.Marshal(config)
		if err != nil {
			log.Println(err)
		}
		link := fmt.Sprint("vmess://", base64.StdEncoding.EncodeToString(byteConfig), "\n")
		links += link
	}
	return links
}

func ConfigToQuantumult(configs []*rayConfig) string {
	links := ""
	templateURL := "%s = vmess, %s, 443, aes-128-cfb, \"%s\", group=%s, over-tls=true, certificate=1, obfs=ws\n"
	for _, config := range configs {
		links += fmt.Sprintf(templateURL, config.Ps, config.Add, config.ID, config.Group)
	}
	return links
}
