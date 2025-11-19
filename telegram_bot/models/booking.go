package models

type Booking struct {
	BookingID int64 `json:"booking_id,omitempty"`

	DateFrom *string `json:"date_from"`
	DateTo   *string `json:"date_to"`

	FullCost int64 `json:"full_cost,omitempty"`

	ParkingPlaceID *int64 `json:"parking_place_id"`

	Status string `json:"status,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

func NewBooking() *Booking {
	booking := new(Booking)
	booking.DateTo = new(string)
	booking.DateFrom = new(string)
	booking.ParkingPlaceID = new(int64)
	return booking
}
