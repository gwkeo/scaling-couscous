package main

import (
	"context"
	"github.com/gwkeo/scaling-couscous/internal/app/endpoints/user"
	sqliteRepo "github.com/gwkeo/scaling-couscous/internal/app/repo/sqlite"
	"github.com/gwkeo/scaling-couscous/internal/app/services"
	"github.com/gwkeo/scaling-couscous/internal/config"
	"log"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.TODO()
	usersRepo, err := sqliteRepo.NewUsersRepo(ctx, cfg.DBPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	usersService := services.NewUsersService(usersRepo)

	userHandler := user.New(ctx, usersService, cfg)
	if err = userHandler.Start(); err != nil {
		log.Fatal(err.Error())
	}
}
