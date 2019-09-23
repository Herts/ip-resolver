package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func handleIp(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query()["host"]
	addrs, err := net.LookupIP(host[0])
	if err != nil {
		log.Println(err)
	}
	log.Println(addrs)
	w.Header().Set("Content-Type", "text/plain")
	for _, addr := range addrs {
		fmt.Fprint(w, addr)
	}
}

func main() {
	http.HandleFunc("/", handleIp)
	log.Fatal(http.ListenAndServe(":1096", nil))

}
