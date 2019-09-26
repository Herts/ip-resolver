package main

import (
	"fmt"
	"github.com/matryer/way"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func (s *manageServer) routes() {
	s.router = way.NewRouter()
	s.router.HandleFunc("GET", "/dns/add", s.handleDnsAdd())
	s.router.HandleFunc("GET", "/dns/ray/:userId", s.handleConfig())
}

func main() {
	log.Println("Running")
	viper.SetConfigName("config1")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	s := newManageServer(viper.GetString("database.conn"))
	http.Handle("/dns/", s.router)

	log.Fatal(http.ListenAndServe(":1096", nil))

}
