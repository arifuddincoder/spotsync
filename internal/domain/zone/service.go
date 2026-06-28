package zone

import (
	"errors"
	"time"

	"spotsync/internal/domain/zone/dto"
)

var ErrZoneNotFound = errors.New("parking zone not found")

type Service interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneWithAvailability, error)
	GetZoneByID(id uint) (*dto.ZoneWithAvailability, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.CreateZone(&zone); err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     zone.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) GetAllZones() ([]dto.ZoneWithAvailability, error) {
	zones, err := s.repo.GetAllZones()
	if err != nil {
		return nil, err
	}

	result := make([]dto.ZoneWithAvailability, 0, len(zones))
	for _, z := range zones {
		result = append(result, toAvailabilityDTO(z))
	}
	return result, nil
}

func (s *service) GetZoneByID(id uint) (*dto.ZoneWithAvailability, error) {
	zone, err := s.repo.GetZoneByID(id)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, ErrZoneNotFound
	}

	out := toAvailabilityDTO(*zone)
	return &out, nil
}

func toAvailabilityDTO(z ZoneWithCount) dto.ZoneWithAvailability {
	return dto.ZoneWithAvailability{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: z.AvailableSpots,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt.Format(time.RFC3339),
	}
}
