package handlers

import (
	"github.com/h4x4d/parking_net/parking/internal/models"
	"github.com/h4x4d/parking_net/pkg/domain"
)

func ToDomainParking(api *models.ParkingPlace) *domain.ParkingPlace {
	if api == nil {
		return nil
	}
	
	p := &domain.ParkingPlace{
		Name:       getStringValue(api.Name),
		City:       getStringValue(api.City),
		Address:    getStringValue(api.Address),
		HourlyRate: float64(api.HourlyRate),
		Capacity:   int(api.Capacity),
	}
	
	if api.ParkingType != "" {
		p.Type = domain.ParkingType(api.ParkingType)
	}
	
	return p
}

func ToDomainParkingUpdate(api *models.ParkingPlace) *domain.ParkingPlace {
	if api == nil {
		return nil
	}
	
	p := &domain.ParkingPlace{}
	
	if api.Name != nil {
		p.Name = *api.Name
	}
	if api.City != nil {
		p.City = *api.City
	}
	if api.Address != nil {
		p.Address = *api.Address
	}
	if api.ParkingType != "" {
		p.Type = domain.ParkingType(api.ParkingType)
	}
	if api.HourlyRate != 0 {
		p.HourlyRate = float64(api.HourlyRate)
	}
	if api.Capacity != 0 {
		p.Capacity = int(api.Capacity)
	}
	
	return p
}

func ToAPIParking(d *domain.ParkingPlace) *models.ParkingPlace {
	if d == nil {
		return nil
	}
	
	return &models.ParkingPlace{
		ID:          d.ID,
		Name:        stringPtr(d.Name),
		City:        stringPtr(d.City),
		Address:     stringPtr(d.Address),
		ParkingType: string(d.Type),
		HourlyRate:  int64(d.HourlyRate),
		Capacity:    int64(d.Capacity),
		OwnerID:     d.OwnerID,
	}
}

func ToAPIParkingList(domains []*domain.ParkingPlace) []*models.ParkingPlace {
	result := make([]*models.ParkingPlace, 0, len(domains))
	for _, d := range domains {
		result = append(result, ToAPIParking(d))
	}
	return result
}

func ToDomainUser(api *models.User) *domain.User {
	if api == nil {
		return nil
	}
	
	return &domain.User{
		ID:         api.UserID,
		Role:       domain.UserRole(api.Role),
		TelegramID: int64(api.TelegramID),
	}
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringPtr(s string) *string {
	return &s
}

