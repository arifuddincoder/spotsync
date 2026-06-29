package reservation

import "gorm.io/gorm"

// একই plate-এর একসাথে একটার বেশি active reservation আটকায়।
// status active না থাকলে (cancelled/completed) বা soft-deleted হলে
// index প্রযোজ্য নয় — তাই পরে আবার একই plate দিয়ে reserve করা যাবে।
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
