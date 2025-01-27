package main

import (
	"context"
	"fmt"
	sqlite2 "github.com/gwkeo/scaling-couscous/internal/app/repo/sqlite"
	"github.com/gwkeo/scaling-couscous/internal/app/services"
	"github.com/gwkeo/scaling-couscous/internal/config"
	"log"
	"net/http"
)

const (
	PORT = "8000"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	usersRepo, err := sqlite2.NewUsersRepo(ctx, cfg.DBPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = sqlite2.NewProductsRepo(ctx, cfg.DBPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	_ = services.NewUsersService(usersRepo)

	mux := http.NewServeMux()
	mux.HandleFunc(`/user/`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
		log.Println("hello world")
	})

	fmt.Printf("Listening: http://localhost:%v\n", PORT)

	if err = http.ListenAndServe(":"+PORT, mux); err != nil {
		log.Fatal(err.Error())
	}
}
