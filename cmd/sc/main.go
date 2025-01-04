package main

import (
	"context"
	"github.com/gwkeo/scaling-couscous/internal/config"
	"github.com/gwkeo/scaling-couscous/internal/repo/sqlite"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	println(cfg.DBPath)

	ctx := context.Background()
	_, err = sqlite.NewUsersRepo(ctx, cfg.DBPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = sqlite.NewProductsRepo(ctx, cfg.DBPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
		log.Println("hello world")
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err.Error())
	}
}
