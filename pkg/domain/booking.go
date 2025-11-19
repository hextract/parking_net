package domain

import (
	"time"
)

type BookingStatus string

const (
	BookingStatusWaiting   BookingStatus = "Waiting"
	BookingStatusConfirmed BookingStatus = "Confirmed"
	BookingStatusCanceled  BookingStatus = "Canceled"
)

type Booking struct {
	ID            int64
	DateFrom      time.Time
	DateTo        time.Time
	ParkingPlaceID int64
	FullCost      float64
	Status        BookingStatus
	UserID        string
}

func (b *Booking) IsValid() error {
	if b.DateFrom.IsZero() {
		return ErrInvalidDateFrom
	}
	if b.DateTo.IsZero() {
		return ErrInvalidDateTo
	}
	if b.DateFrom.After(b.DateTo) {
		return ErrInvalidDateRange
	}
	if b.ParkingPlaceID == 0 {
		return ErrInvalidParkingPlaceID
	}
	if b.UserID == "" {
		return ErrInvalidUserID
	}
	return nil
}

func (b *Booking) CalculateCost(hourlyRate float64) {
	duration := b.DateTo.Sub(b.DateFrom)
	hours := duration.Hours()
	b.FullCost = hourlyRate * hours
}

