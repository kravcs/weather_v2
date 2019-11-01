package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/kravcs/gogo/cache"
	"github.com/kravcs/gogo/cache/redis"
	"github.com/kravcs/gogo/config"
	"github.com/kravcs/gogo/handler"
)

var storage cache.Storage

func init() {
	var err error
	connectionInfo := fmt.Sprintf(
		"%s://%s:%d",
		config.Configuration.Database.Driver,
		config.Configuration.Database.Host,
		config.Configuration.Database.Port,
	)
	if storage, err = redis.NewStorage(connectionInfo); err != nil {
		log.Fatalf("Redis error: %v", err)
	}
	handler.Storage = storage
	fmt.Println("Redis started!")
}

func main() {
	router := newRouter()
	serverStart(router)
}

func newRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/weather", handler.GetWeatherHandler).Methods("GET")

	return router
}

func serverStart(router *mux.Router) {
	port := strconv.Itoa(config.Configuration.Server.Port)
	server := http.Server{
		Addr:        ":" + port,
		Handler:     router,
		ReadTimeout: 10 * time.Second,
	}
	fmt.Println("Server started!")
	log.Fatal(server.ListenAndServe())
}
