package ride

import (
  "fmt"
  "uber/internal/model"
  "math/rand"
  "time"
)

type RideManager struct {
   Drivers []model.Driver
   Trips   []model.Trip
   Users   []model.User
}

func NewRideManager() *RideManager {
	return &RideManager{}
}

func (rm * RideManager) AddUsers(user []model.User) {
	rm.Users =append(rm.Users, user...)
}

func (rm * RideManager) AddDrivers(driver []model.Driver) {
	rm.Drivers =append(rm.Drivers, driver...)
}


func (rm *RideManager) RequestRide(user model.User) (*model.Trip,error) {

	var selectedDriver *model.Driver
	
	for i,driver := range rm.Drivers {
		if driver.Status == "available" {
			selectedDriver = &rm.Drivers[i]
			break;
		}
	}

	if selectedDriver == nil {
		return nil,fmt.Errorf("no available drivers")
	}

	trip := model.Trip{
		ID: fmt.Sprintf("%d",rand.Int()),
		User: user,
		Driver: *selectedDriver,
		Pickup: user.Location,
		Dropoff: model.Location{Latitude: rand.Float64(),Longitude: rand.Float64()},
		Status: "requested",
		Distance: rand.Float64()*10,
		Duration: time.Duration(rand.Intn(30)+5)*time.Minute,
		StartTime: time.Time{},
		EndTime: time.Time{},
	}
     
	selectedDriver.Status = "unavailable"

	rm.Trips = append(rm.Trips, trip)

	return &trip,nil
}

func (rm *RideManager) AcceptRide(driver model.Driver,tripID string) (*model.Trip,error){
	var trip *model.Trip
	for i:= range rm.Trips{
		if rm.Trips[i].ID == tripID {
			trip = &rm.Trips[i]
			break
		}
	}

	if trip == nil {
		return nil,fmt.Errorf("trip not found")
	}

	trip.Status="in-progress"
	trip.StartTime = time.Now()

	return trip,nil
}

func (rm *RideManager) CompleteRide(tripID string, rating int, review string) (*model.Trip, error) {
	// Find the trip by ID
	var trip *model.Trip
	for i := range rm.Trips {
		if rm.Trips[i].ID == tripID {
			trip = &rm.Trips[i]
			break
		}
	}

	if trip == nil {
		return nil, fmt.Errorf("trip not found")
	}

	if trip.Status != "in-progress" {
		return nil, fmt.Errorf("ride must be in progress to complete")
	}

	// Set the trip status to completed
	trip.Status = "completed"
	trip.EndTime = time.Now()

	// Calculate the payment amount
	amount := CalculatePayment(trip.Distance, int(trip.Duration.Minutes()))

	// Create the Payment object and add it to the trip
	trip.Payment = &model.Payment{
		UserID:    trip.User.ID,
		Amount:    amount,
		Status:    "paid",  // Assume the payment is successful
		PaymentTime: time.Now(),
	}

	// Save the rating if the ride is completed
	trip.Rating = &model.Rating{
		DriverID: trip.Driver.ID,
		UserID:   trip.User.ID,
		Score:    rating,
		Review:   review,
	}

	// Update the driver status to "available"
	for i := range rm.Drivers {
		if rm.Drivers[i].ID == trip.Driver.ID {
			rm.Drivers[i].Status = "available"
		}
	}

	// Optionally, you can update the driver's average rating here.

	return trip, nil
}
