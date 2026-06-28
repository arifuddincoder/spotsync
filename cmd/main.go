package main

import (
	"spotsync/internal/config"
	"spotsync/internal/domain/user"
	"spotsync/internal/server"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)
	db.AutoMigrate(&user.User{})

	server.Start(db, cfg)
}
