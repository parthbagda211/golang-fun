package model

import "time"

type Location struct {
	Latitude float64
	Longitude float64
}

type User struct {
	ID string
	Name string
	Location Location
}

type Driver struct {
	ID string
	Name string
	Location Location
	Status string
}

type Rating struct {
	DriverID string
	UserID string
	Score int 
	Review string
}

type Payment struct{
	UserID string
	Amount float64
	Status string
	PaymentTime time.Time
}

type Trip struct {
	ID         string
	User       User
	Driver     Driver
	Pickup     Location
	Dropoff    Location
	Status     string // "requested", "in-progress", "completed", "canceled"
	Distance   float64
	Duration   time.Duration
	StartTime  time.Time
	EndTime    time.Time
	Rating     *Rating
	Payment    *Payment
}