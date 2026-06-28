package main

import (
	"spotsync/internal/config"
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"
	"spotsync/internal/server"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)
	db.AutoMigrate(&user.User{}, &zone.ParkingZone{})
	user.SeedAdmin(user.NewRepository(db), cfg)
	server.Start(db, cfg)
}
