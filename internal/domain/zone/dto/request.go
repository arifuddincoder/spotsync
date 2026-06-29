package dto

type CreateZoneRequest struct {
	Name          string  `json:"name" validate:"required"`
	Type          string  `json:"type" validate:"required,oneof=general ev_charging covered"`
	TotalCapacity int     `json:"total_capacity" validate:"required,gt=0"`
	PricePerHour  float64 `json:"price_per_hour" validate:"required,gt=0"`
}

type UpdateZoneRequest struct {
	Name          *string  `json:"name" validate:"omitempty"`
	Type          *string  `json:"type" validate:"omitempty,oneof=general ev_charging covered"`
	TotalCapacity *int     `json:"total_capacity" validate:"omitempty,gt=0"`
	PricePerHour  *float64 `json:"price_per_hour" validate:"omitempty,gt=0"`
}
