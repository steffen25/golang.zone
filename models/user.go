package models

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User represents a user account
// Make sure not to expose the password field when marshalling to json
type User struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Admin     bool           `json:"admin"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt mysql.NullTime `json:"updatedAt"`
}

// TODO: Maybe find a better solution to remove the password when marshalling to json
func (u *User) MarshalJSON() ([]byte, error) {
	if !u.UpdatedAt.Valid {
		return json.Marshal(struct {
			ID        int             `json:"id"`
			Name      string          `json:"name"`
			Email     string          `json:"email"`
			CreatedAt time.Time          `json:"createdAt"`
			UpdatedAt *mysql.NullTime `json:"updatedAt"`
		}{u.ID, u.Name, u.Email, u.CreatedAt, nil})
	}
	return json.Marshal(struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}{u.ID, u.Name, u.Email, u.CreatedAt, u.UpdatedAt.Time})
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

func (u *User) IsAdmin() bool {
	return u.Admin == true
}
