package main

import (
	"fmt"
	"log"
	"uber/internal/model"
	"uber/internal/ride"
)

func main() {
	// Initialize a simple RideManager and simulate the system.
	rideManager := ride.NewRideManager()

	// Create Users and Drivers
	users := []model.User{
		{ID: "1", Name: "Alice", Location: model.Location{Latitude: 40.748817, Longitude: -73.985428}},
		{ID: "2", Name: "Bob", Location: model.Location{Latitude: 40.730610, Longitude: -73.935242}},
	}

	drivers := []model.Driver{
		{ID: "d1", Name: "Driver1", Location: model.Location{Latitude: 40.748900, Longitude: -73.985500}, Status: "available"},
		{ID: "d2", Name: "Driver2", Location: model.Location{Latitude: 40.730500, Longitude: -73.935400}, Status: "available"},
	}

	rideManager.AddUsers(users)
	rideManager.AddDrivers(drivers)

	// Simulate a ride request and completion
	trip, err := rideManager.RequestRide(users[0])
	if err != nil {
		log.Fatal("Error requesting ride:", err)
	}

	trip, err = rideManager.AcceptRide(drivers[0], trip.ID)
	if err != nil {
		log.Fatal("Error accepting ride:", err)
	}

	trip, err = rideManager.CompleteRide(trip.ID, 5, "Great ride! Would travel again.")
	if err != nil {
		log.Fatal("Error completing ride:", err)
	}

	fmt.Println("Trip completed successfully:", *trip)



	fmt.Printf("\nCompleted Trip Details:\n")
	fmt.Printf("Trip ID: %s\n", trip.ID)
	fmt.Printf("User: %s\n", trip.User.Name)
	fmt.Printf("Driver: %s\n", trip.Driver.Name)
	fmt.Printf("Rating: %d/5 - %s\n", trip.Rating.Score, trip.Rating.Review)
	fmt.Printf("Payment: $%.2f\n", trip.Payment.Amount)
	fmt.Printf("Payment Status: %s\n", trip.Payment.Status)
}
