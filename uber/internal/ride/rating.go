package ride

import (
	"uber/internal/model"
	"fmt"
)

// SubmitRating allows a user to rate a driver after completing a trip.
func SubmitRating(trip *model.Trip, score int, review string) error {
	if trip.Status != "completed" {
		return fmt.Errorf("cannot rate a trip that is not completed")
	}

	// Save the rating
	trip.Rating = &model.Rating{
		DriverID: trip.Driver.ID,
		UserID:   trip.User.ID,
		Score:    score,
		Review:   review,
	}

	return nil
}
