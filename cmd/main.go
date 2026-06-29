package main

import (
	"spotsync/internal/config"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"
	"spotsync/internal/server"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)
	db.AutoMigrate(&user.User{}, &zone.ParkingZone{})
	if err := reservation.Migrate(db); err != nil { // partial index সহ
		panic("failed to migrate reservations: " + err.Error())
	}
	user.SeedAdmin(user.NewRepository(db), cfg)
	server.Start(db, cfg)
}
