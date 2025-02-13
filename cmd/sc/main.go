package main

import (
	"context"
	"github.com/gwkeo/scaling-couscous/internal/app/endpoints/user"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/sqlite"
	"github.com/gwkeo/scaling-couscous/internal/app/services"
	"github.com/gwkeo/scaling-couscous/internal/config"
	"log"
	"os/signal"
	"syscall"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	usersRepo, err := sqlite.NewUsersRepo(ctx, cfg.DBPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	usersService := services.NewUsersService(usersRepo)

	userHandler := user.New(ctx, usersService, cfg)

	go func() {
		if err = userHandler.Start(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	<-ctx.Done()

	if err := usersRepo.CloseDBConnection(); err != nil {
		log.Println(err.Error())
	}

	stop()
	log.Println("shutting down...")
}
