package reservation

import (
	"errors"
	"time"

	"spotsync/internal/domain/reservation/dto"
)

var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrNotOwner            = errors.New("you can only cancel your own reservations")
	ErrAlreadyCancelled    = errors.New("reservation is already cancelled")
)

type Service interface {
	Reserve(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	GetAllReservations() ([]dto.AdminReservationResponse, error)
	CancelReservation(userID uint, role string, reservationID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Reserve(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	res := Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       StatusActive,
	}

	if err := s.repo.CreateReservation(&res); err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    res.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	rows, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.MyReservationResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneBrief{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}
	return out, nil
}

func (s *service) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	rows, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	out := make([]dto.AdminReservationResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.UserBrief{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
				Role:  r.User.Role,
			},
			Zone: dto.ZoneBrief{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}
	return out, nil
}

func (s *service) CancelReservation(userID uint, role string, reservationID uint) error {
	res, err := s.repo.GetByID(reservationID)
	if err != nil {
		return err
	}
	if res == nil {
		return ErrReservationNotFound
	}

	// driver শুধু নিজেরটা cancel করতে পারবে; admin যেকোনোটা পারবে
	if res.UserID != userID && role != "admin" {
		return ErrNotOwner
	}

	if res.Status == StatusCancelled {
		return ErrAlreadyCancelled
	}

	// status = cancelled করলেই spot ফ্রি হয়ে যায় (available_spots শুধু active গোনে)
	return s.repo.UpdateStatus(res, StatusCancelled)
}
