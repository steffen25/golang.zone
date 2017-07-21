package repositories

import (
	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/database"
	"database/sql"
)

type UserInterface interface {
	Create(*models.User) error
	GetAll() ([]*models.User, error)
	FindById(id int) (*models.User, error)
	Exists(email string) bool
	Delete(id int) error
	Update(id int) error
}

type UserRepository struct {
	*database.DB
}

func (ur *UserRepository) Create(u *models.User) error {

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

func (ur *UserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User

	rows, err := ur.DB.Query("SELECT id, name, email, created_at from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(models.User)
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
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

func (ur *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := models.User{}

	err := ur.DB.QueryRow("SELECT id, name, email, password, created_at FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) FindById(id int) (*models.User, error) {
	user := models.User{}

	err := ur.DB.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Exists(email string) bool {

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