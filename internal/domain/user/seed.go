package user

import (
	"log"

	"spotsync/internal/config"
)

func SeedAdmin(repo Repository, cfg *config.Config) {
	if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
		log.Println("Admin credentials not set, skipping seeding")
		return
	}

	existing, err := repo.GetUserByEmail(cfg.AdminEmail)
	if err != nil {
		log.Println("Failed to check admin existence:", err)
		return
	}
	if existing != nil {
		log.Println("Admin already exists, skipping seeding")
		return
	}

	name := cfg.AdminName
	if name == "" {
		name = "Admin"
	}

	admin := User{
		Name:  name,
		Email: cfg.AdminEmail,
		Role:  RoleAdmin,
	}

	if err := admin.hashPassword(cfg.AdminPassword); err != nil {
		log.Println("Failed to hash admin password:", err)
		return
	}

	if err := repo.CreateUser(&admin); err != nil {
		log.Println("Failed to seed admin:", err)
		return
	}

	log.Println("Admin seeded successfully")
}
