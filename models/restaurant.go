package models

// Restaurant models a restaurant
type Restaurant struct {
	ID          string
	Name        string
	Latitude    float64
	Longitude   float64
	Cuisine     string
	Description string
	Website     string
	NYCLink     string
}
