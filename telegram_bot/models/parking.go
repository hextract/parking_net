package models

type ParkingPlace struct {
	Address *string `json:"address"`
	City    *string `json:"city"`
	
	HourlyRate  int64  `json:"hourly_rate,omitempty"`
	Capacity    int64  `json:"capacity,omitempty"`
	ParkingType string `json:"parking_type,omitempty"`
	
	ID      int64  `json:"id,omitempty"`
	Name    *string `json:"name"`
	OwnerID string `json:"owner_id,omitempty"`
}

func NewParkingPlace() *ParkingPlace {
	parkingPlace := new(ParkingPlace)
	parkingPlace.Address = new(string)
	parkingPlace.City = new(string)
	parkingPlace.Name = new(string)
	return parkingPlace
}
