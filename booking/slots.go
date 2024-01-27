package booking

import (
	"context"
	"errors"
	"time"
)

const DefaultBookingDuration = 1 * time.Hour

type BookableSlot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type SlotsParams struct{}
type SlotsResponse struct{ Slots []BookableSlot }

//encore:api public method=GET path=/slots/:from
func GetBookableSlots(ctx context.Context, from string) (*SlotsResponse, error) {
	fromDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, err
	}

	const numDays = 7
	
	var slots []BookableSlot

	for i := 0; i <= numDays; i++ {
		date := fromDate.AddDate(0, 0, i)
		daySlots, err :=
		slotsResponse.Slots = append(times, fromDate)
	}

	if len(i2) < 1 {
		return nil, errors.New("No time slots available")
	}
	return &slotsResponse, nil
}
