package main

import (
	"log"
	"net/http"

	"market_system/config"
	"market_system/controller"
	"market_system/storage/postgres"
)

func main() {

	var cfg = config.Load()

	pgStorage, err := postgres.NewConnectionPostgres(&cfg)
	if err != nil {
		panic(err)
	}

	handler := controller.NewController(&cfg, pgStorage)

	http.HandleFunc("/client", handler.Client)
	http.HandleFunc("/order_products", handler.OrderProduct)
	http.HandleFunc("/order", handler.Order)
	http.HandleFunc("/category", handler.Category)
	http.HandleFunc("/product", handler.Product)
	http.HandleFunc("/branch", handler.Branch)

	log.Println("Listening:", cfg.ServiceHost+cfg.ServiceHTTPPort, "...")
	if err := http.ListenAndServe(cfg.ServiceHost+cfg.ServiceHTTPPort, nil); err != nil {
		panic("Listent and service panic:" + err.Error())
	}
}
