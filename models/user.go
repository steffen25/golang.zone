package models

import (
	"time"
)

// User represents a user account
// Make sure not to expose the password field when marshalling to json
type User struct {
	ID     int    		`json:"id"`
	Name  string 		`json:"name"`
	Email   string 		`json:"email"`
	Password string 	`json:"password"`
	CreatedAt time.Time	`json:"createdAt"`
	UpdatedAt time.Time	`json:"updatedAt"`
}