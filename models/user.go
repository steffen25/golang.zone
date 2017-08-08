package models

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
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

// TODO: Maybe find a better solution to remove the password when marshalling to json
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID int 			`json:"id"`
		Name string 	`json:"name"`
		Email string 	`json:"email"`
		CreatedAt string	`json:"createdAt"`
	}{u.ID, u.Name, u.Email, u.CreatedAt.Format(time.RFC3339)})
}

func (u *User) NewUser(name string, email string, password string) {
	// validate length here?
}

func (u *User) SetPassword(password string) {
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.Password = string(pwhash)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}

	return true
}