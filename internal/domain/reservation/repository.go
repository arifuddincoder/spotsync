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

// CreateReservation — "EV Spot Bottleneck" সমাধান।
// পুরো check + insert একটা Transaction-এর ভেতরে, এবং zone row-টা
// FOR UPDATE দিয়ে lock করা হয়, যাতে একই zone-এর জন্য আসা concurrent
// request গুলো একটার পর একটা (serialized) চলে — কখনো over-capacity হবে না।
func (r *repository) CreateReservation(res *Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// ১) zone row-টা lock করো (SELECT ... FOR UPDATE)
		var z zone.ParkingZone
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&z, res.ZoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrZoneNotFound
			}
			return err
		}

		// ২) একই plate-এর active reservation আছে কিনা চেক (duplicate guard → 409)
		var dup int64
		if err := tx.Model(&Reservation{}).
			Where("license_plate = ? AND status = ?", res.LicensePlate, StatusActive).
			Count(&dup).Error; err != nil {
			return err
		}
		if dup > 0 {
			return ErrDuplicatePlate
		}

		// ৩) এই zone-এ এখন কতগুলো active reservation আছে গুনি
		var active int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", res.ZoneID, StatusActive).
			Count(&active).Error; err != nil {
			return err
		}

		// ৪) capacity চেক — জায়গা না থাকলে ZoneFull
		if int(active) >= z.TotalCapacity {
			return ErrZoneFull
		}

		// ৫) একই tx-এর ভেতরেই reservation তৈরি করি।
		//    Omit(clause.Associations) → শুধু UserID/ZoneID column বসবে,
		//    ভুল করে User/Zone টেবিল upsert হবে না।
		res.Status = StatusActive
		if err := tx.Omit(clause.Associations).Create(res).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrDuplicatePlate // → handler 409 দেবে
			}
			return err
		}

		return nil // nil রিটার্ন করলে tx commit হবে; error দিলে rollback
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
