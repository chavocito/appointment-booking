package booking

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
	for i := 0; i < numDays; i++ {
		date := fromDate.AddDate(0, 0, i)
		daySlots, err := bookableSlotsForDay(date)
		if err != nil {
			return nil, err
		}
		slots = append(slots, daySlots...)
	}

	return &SlotsResponse{Slots: slots}, nil
}

func bookableSlotsForDay(date time.Time) ([]BookableSlot, error) {
	//hardcode current start to end time as 09:00 to 17:00
	availableStartTime := pgtype.Time{
		Valid:        true,
		Microseconds: int64(9*3600) * 1e6,
	}

	availableEndTime := pgtype.Time{
		Valid:        true,
		Microseconds: int64(17*3600) * 1e6,
	}

	availStart := date.Add(time.Duration(availableStartTime.Microseconds) * time.Microsecond)
	availEnd := date.Add(time.Duration(availableEndTime.Microseconds) * time.Microsecond)

	// Compute the bookable slots in this day, based on availability.
	var slots []BookableSlot
	start := availStart
	for {
		end := start.Add(DefaultBookingDuration)
		if end.After(availEnd) {
			break
		}
		slots = append(slots, BookableSlot{
			Start: start,
			End:   end,
		})
		start = end
	}

	return slots, nil
}
