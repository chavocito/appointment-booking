package booking

import (
	"context"
	"log"
	"time"

	"encore.app/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	bookingDB = sqldb.NewDatabase("booking", sqldb.DatabaseConfig{
		Migrations: "./db/migrations",
	})
	pgxdb = sqldb.Driver[*pgxpool.Pool](bookingDB)
	query = db.New(pgxdb)
)

type Booking struct {
	ID    int64     `json:"id"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Email string    `json:"email"`
}

type BookParams struct {
	ID    int64     `json:"id"`
	Start time.Time `json:"start"`
	Email string    `json:"email"`
}

//encore:api public method=POST path=/booking
func Book(ctx context.Context, p *BookParams) *errs.Builder {
	errorBuilder := errs.B()

	now := time.Now()
	if p.Start.Before(now) {
		return errorBuilder.Code(errs.InvalidArgument).Msg("start time must be in the future").Err()
	}

	tx, err := pgxdb.Begin(ctx)
	if err != nil {
		return errorBuilder.Cause(err).Code(errs.Unavailable).Msg("failed to start transaction").Err()
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		e := tx.Rollback(ctx)
		if e != nil {
			log.Fatalf("Error rolling back transaction %s", e)
		}
	}(tx, context.Background())

	_, err = query.InsertBooking(ctx, db.InsertBookingParams{
		StartTime: pgtype.Timestamp{Time: p.Start, Valid: true},
		EndTime:   pgtype.Timestamp{Time: p.Start.Add(DefaultBookingDuration), Valid: true},
		Email:     p.Email,
	})

	if err != nil {
		return errorBuilder.Cause(err).Code(errs.Aborted).Msg("failed to insert record into database")
	}

	if err := tx.Commit(ctx); err != nil {
		
	}
}
