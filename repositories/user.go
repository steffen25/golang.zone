package repositories

import (
	"database/sql"
	"errors"

	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/database"
)

type UserRepository interface {
	Create(u *models.User) error
	GetAll() ([]*models.User, error)
	FindById(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Exists(email string) bool
	Delete(id int) error
	Update(u *models.User) error
}

type userRepository struct {
	*database.DB
}

func NewUserRespository(db *database.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Create(u *models.User) error {

	// Check if an user already exists with the email
	// Prepare statement for inserting data
	stmt, err := ur.DB.Prepare("INSERT INTO users SET name=?, email=?, password=?, created_at=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.Name, u.Email, u.Password, u.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) Update(u *models.User) error {

	// Check if an user already exists with the email
	// Prepare statement for inserting data
	stmt, err := ur.DB.Prepare("UPDATE users SET name=?, email=?, password=? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.Name, u.Email, u.Password, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) GetAll() ([]*models.User, error) {
	var users []*models.User

	rows, err := ur.DB.Query("SELECT id, name, email, created_at, admin from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(models.User)
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.Admin)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *userRepository) FindByEmail(email string) (*models.User, error) {
	user := models.User{}

	err := ur.DB.QueryRow("SELECT id, name, email, password, created_at, admin FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.Admin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) FindById(id int) (*models.User, error) {
	user := models.User{}

	err := ur.DB.QueryRow("SELECT id, name, email, password, created_at, admin FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.Admin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) Exists(email string) bool {

	// Check if an user already exists with the email
	var exists bool
	stmt, err := ur.DB.Prepare("SELECT email FROM users WHERE email = ?")
	if err != nil {
		return true
	}
	defer stmt.Close()
	err = stmt.QueryRow(email).Scan(exists)
	if err != nil && err != sql.ErrNoRows {
		return true
	}

	return exists
}

func (ur *userRepository) Delete(id int) error {

	return errors.New("hej")
}