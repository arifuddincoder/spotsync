package reservation

import (
	"errors"

	"spotsync/internal/domain/zone"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneNotFound   = errors.New("parking zone not found")
	ErrZoneFull       = errors.New("parking zone is full, no available spots")
	ErrDuplicatePlate = errors.New("this license plate already has an active reservation")
)

type Repository interface {
	CreateReservation(res *Reservation) error
	GetByUserID(userID uint) ([]Reservation, error)
	GetAll() ([]Reservation, error)
	GetByID(id uint) (*Reservation, error)
	UpdateStatus(res *Reservation, status string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateReservation(res *Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var z zone.ParkingZone
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&z, res.ZoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrZoneNotFound
			}
			return err
		}

		var dup int64
		if err := tx.Model(&Reservation{}).
			Where("license_plate = ? AND status = ?", res.LicensePlate, StatusActive).
			Count(&dup).Error; err != nil {
			return err
		}
		if dup > 0 {
			return ErrDuplicatePlate
		}

		var active int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", res.ZoneID, StatusActive).
			Count(&active).Error; err != nil {
			return err
		}

		if int(active) >= z.TotalCapacity {
			return ErrZoneFull
		}

		res.Status = StatusActive
		if err := tx.Omit(clause.Associations).Create(res).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrDuplicatePlate
			}
			return err
		}

		return nil
	})
}

func (r *repository) GetByUserID(userID uint) ([]Reservation, error) {
	var list []Reservation
	err := r.db.
		Preload("Zone").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&list).Error
	return list, err
}

func (r *repository) GetAll() ([]Reservation, error) {
	var list []Reservation
	err := r.db.
		Preload("User").
		Preload("Zone").
		Order("id DESC").
		Find(&list).Error
	return list, err
}

func (r *repository) GetByID(id uint) (*Reservation, error) {
	var res Reservation
	if err := r.db.First(&res, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *repository) UpdateStatus(res *Reservation, status string) error {
	return r.db.Model(res).Update("status", status).Error
}
