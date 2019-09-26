package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/matryer/way"
	"github.com/spf13/viper"
	"github.com/tomasen/realip"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

type manageServer struct {
	router  *way.Router
	mysqlDb *sql.DB
}

type IpInfoResp struct {
	Bogon    string `json:"bogon"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Error    struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	} `json:"error"`
}

func newManageServer(conn string) *manageServer {
	var err error
	s := &manageServer{}
	s.mysqlDb, err = mysqlDb(conn)
	if err != nil {
		log.Panic(err)
	}
	s.routes()
	return s
}

func IpInfo(ip, ipInfoKey string) *IpInfoResp {
	url := fmt.Sprintf("http://ipinfo.io/%s?token=%s", ip, ipInfoKey)
	log.Println(url)
	resp, err := http.Get(url)
	log.Println("Getting")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	var ipInfo IpInfoResp
	log.Println("Decoding")
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	if err != nil {
		log.Println(err)
	}
	return &ipInfo
}

func IpCountry(ip string) string {
	url := fmt.Sprintf("http://ipinfo.io/%s/country", ip)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode != 200 {
		return ""
	}
	bytesCountry, err := ioutil.ReadAll(resp.Body)
	country := strings.ReplaceAll(string(bytesCountry), "\n", "")
	return country

}

type DNSCreateReq struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Priority int    `json:"priority"`
	Proxied  bool   `json:"proxied"`
}

func (s *manageServer) handleIp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host := r.URL.Query()["host"]
		addrs, err := net.LookupIP(host[0])
		if err != nil {
			log.Println(err)
		}
		log.Println(addrs)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "%s", addrs[0])
	}
}

func (s *manageServer) handleDnsAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)
		log.Println(ip)
		country := IpCountry(ip)
		if country == "" {
			http.NotFound(w, r)
			return
		}
		name := s.SetDNS(ip, country)
		name = strings.ToLower(name)
		fmt.Fprint(w, name)
	}
}

func (s *manageServer) SetDNS(ip string, country string) string {
	row := s.mysqlDb.QueryRow("SELECT name FROM t_server WHERE ip = ?", ip)
	var name string
	row.Scan(&name)
	log.Println(name)
	if len(name) == 0 {
		row := s.mysqlDb.QueryRow("SELECT MAX(region_idx) FROM t_server WHERE region = ?", country)
		var index sql.NullInt32
		err := row.Scan(&index)
		if err != nil {
			log.Println(err)
		}
		var intIdx int32
		if index.Valid {
			intIdx = index.Int32 + 1
		}
		name = fmt.Sprint(country, intIdx)
		CFCreateDNS(name, ip, viper.GetString("cloudflare.zoneid"))
		_, err = s.mysqlDb.Exec("INSERT INTO t_server (ip, name, region, region_idx) VALUES (?,?,?,?)", ip, name, country, intIdx)
		if err != nil {
			log.Println(err)
		}
	}
	return name
}
