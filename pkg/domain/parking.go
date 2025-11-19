package domain

type ParkingType string

const (
	ParkingTypeOutdoor     ParkingType = "outdoor"
	ParkingTypeCovered     ParkingType = "covered"
	ParkingTypeUnderground ParkingType = "underground"
	ParkingTypeMultiLevel  ParkingType = "multi-level"
)

type ParkingPlace struct {
	ID         int64
	Name       string
	City       string
	Address    string
	Type       ParkingType
	HourlyRate float64
	Capacity   int
	OwnerID    string
}

func (p *ParkingPlace) IsValid() error {
	if p.Name == "" {
		return ErrInvalidParkingName
	}
	if p.City == "" {
		return ErrInvalidParkingCity
	}
	if p.Address == "" {
		return ErrInvalidParkingAddress
	}
	if p.HourlyRate <= 0 {
		return ErrInvalidHourlyRate
	}
	if p.Capacity <= 0 {
		return ErrInvalidCapacity
	}
	if p.Type == "" {
		return ErrInvalidParkingType
	}
	// OwnerID is set by the service layer, not validated here
	return nil
}
