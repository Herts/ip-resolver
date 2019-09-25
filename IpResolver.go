package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

func CFRequest() {
	log.Println("Making request: CFRequest")
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=A", viper.GetString("cloudflare.zoneid")), nil)
	if err != nil {
		log.Println(err)
	}

	req.Header = *CFHeader()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(byteResp))
}

func CFCreateDNS(name, ip, zoneId string) {
	data := &DNSCreateReq{
		Type:    "A",
		Name:    name,
		Content: ip,
		Proxied: false,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneId),
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
	}
	req.Header = *CFHeader()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(byteResp))
}

func CFHeader() *http.Header {
	header := &http.Header{}
	header.Set("X-Auth-Email", viper.GetString("cloudflare.apiemail"))
	header.Set("X-Auth-Key", viper.GetString("cloudflare.apikey"))
	header.Set("Content-Type", "application/json")
	return header
}
