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
	fmt.Fprintf(w, "%s", addrs[0])
}

func main() {
	http.HandleFunc("/ipapi", handleIp)
	log.Fatal(http.ListenAndServe(":1096", nil))

}
