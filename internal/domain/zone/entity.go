package zone

import "gorm.io/gorm"

const (
	TypeGeneral    = "general"
	TypeEVCharging = "ev_charging"
	TypeCovered    = "covered"
)

type ParkingZone struct {
	gorm.Model
	Name          string  `gorm:"type:varchar(100);not null"`
	Type          string  `gorm:"type:varchar(20);not null"`
	TotalCapacity int     `gorm:"not null"`
	PricePerHour  float64 `gorm:"type:decimal(10,2);not null"`
}
