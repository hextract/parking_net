package database_service

import (
	"context"
	"fmt"
)

func (ds *DatabaseService) Delete(ctx context.Context, bookingID int64) error {
	query := `DELETE FROM bookings WHERE id = $1`
	
	result, err := ds.pool.Exec(ctx, query, bookingID)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}
	
	return nil
}

