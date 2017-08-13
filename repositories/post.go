package repositories

import (
	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/database"
	"strconv"
	"log"
)

type PostRepository interface {
	Create (p *models.Post) error
	GetAll() ([]*models.Post, error)
	FindById(id int) (*models.Post, error)
	FindByUser(u *models.User) ([]*models.Post, error)
	Delete(id int) error
	Update(p *models.Post) error
}

type postRepository struct {
	*database.DB
}

func NewPostRepository(db *database.DB) PostRepository {
	return &postRepository{db}
}

func (pr *postRepository) Create(p *models.Post) error {
	stmt, err := pr.DB.Prepare("INSERT INTO posts SET title=?, slug=?, body=?, created_at=?, user_id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(p.Title, p.Slug, p.Body, p.CreatedAt, p.UserID)
	if err != nil {
		return err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	id := strconv.FormatInt(lastId, 10)

	_, err = pr.DB.Exec("UPDATE posts SET slug=? WHERE id=?", p.Slug+"-"+id, lastId)
	if err != nil {
		return err
	}

	return nil
}

func (pr *postRepository) GetAll() ([]*models.Post, error) {
	var posts []*models.Post

	rows, err := pr.DB.Query("SELECT id, title, slug, body, created_at, updated_at, user_id from posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := new(models.Post)
		err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Body, &p.CreatedAt, &p.UpdatedAt, &p.UserID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (pr *postRepository) FindById(id int) (*models.Post, error) {
	return nil, nil
}

func (pr *postRepository) FindByUser(u *models.User) ([]*models.Post, error) {
	return nil, nil
}

func (pr *postRepository) Delete(id int) error {
	return nil
}

func (pr *postRepository) Update(p *models.Post) error {
	return nil
}

