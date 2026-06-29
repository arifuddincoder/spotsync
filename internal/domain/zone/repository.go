package zone

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	CreateZone(zone *ParkingZone) error
	GetAllZones() ([]ZoneWithCount, error)
	GetZoneByID(id uint) (*ZoneWithCount, error)
	UpdateZone(id uint, fields map[string]any) (*ParkingZone, error)
	DeleteZone(id uint) error
}

type ZoneWithCount struct {
	ParkingZone
	AvailableSpots int
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

const availableSpotsWithReservations = `total_capacity - COALESCE((
	SELECT COUNT(*) FROM reservations
	WHERE reservations.zone_id = parking_zones.id
	  AND reservations.status = 'active'
	  AND reservations.deleted_at IS NULL
), 0)`

func (r *repository) availableSpotsExpr() string {
	if r.db.Migrator().HasTable("reservations") {
		return availableSpotsWithReservations
	}
	return "total_capacity"
}

func (r *repository) CreateZone(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) GetAllZones() ([]ZoneWithCount, error) {
	var zones []ZoneWithCount
	expr := r.availableSpotsExpr()

	err := r.db.Table("parking_zones").
		Select("parking_zones.*, (" + expr + ") AS available_spots").
		Where("parking_zones.deleted_at IS NULL").
		Order("parking_zones.id ASC").
		Scan(&zones).Error
	if err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *repository) GetZoneByID(id uint) (*ZoneWithCount, error) {
	var zone ZoneWithCount
	expr := r.availableSpotsExpr()

	result := r.db.Table("parking_zones").
		Select("parking_zones.*, ("+expr+") AS available_spots").
		Where("parking_zones.id = ? AND parking_zones.deleted_at IS NULL", id).
		Take(&zone)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &zone, nil
}

func (r *repository) UpdateZone(id uint, fields map[string]any) (*ParkingZone, error) {
	var z ParkingZone

	if err := r.db.First(&z, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if len(fields) > 0 {
		if err := r.db.Model(&z).Updates(fields).Error; err != nil {
			return nil, err
		}
	}
	return &z, nil
}

func (r *repository) DeleteZone(id uint) error {

	result := r.db.Delete(&ParkingZone{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
