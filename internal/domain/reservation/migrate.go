package reservation

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&Reservation{}); err != nil {
		return err
	}

	return db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS uniq_active_license_plate
		ON reservations (license_plate)
		WHERE status = 'active' AND deleted_at IS NULL
	`).Error
}
