package database

import (
	"github.com/go-pg/pg/v10"
	"github.com/steffen25/golang.zone/api/auth/models"
	"log"
)

type AuthStore struct {
	db *pg.DB
}

func NewAuthStore(db *pg.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (s *AuthStore) FindByEmail(email string) (models.User, error) {
	u := models.User{Email: email}
	err := s.db.Model(&u).
		Column("id", "email", "name", "password", "created_at", "updated_at", "deleted_at").
		Where("email = ?email").
		Relation("Roles").
		First()

	return u, err
}

func (s *AuthStore) FindById(userId int) (models.User, error) {
	u := models.User{ID: userId}
	err := s.db.Model(&u).
		Column("id", "email", "name", "password", "created_at", "updated_at", "deleted_at").
		Where("id = ?id").
		Relation("Roles").
		First()

	return u, err
}

func (s AuthStore) CreateUser(user *models.User) error {
	_, err := s.db.Model(user).Returning("*").Insert()
	if err != nil {
		return err
	}
	log.Println(user)

	return nil
}
