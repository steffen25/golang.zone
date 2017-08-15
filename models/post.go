package models

import (
	"time"
	"github.com/go-sql-driver/mysql"
	"encoding/json"
)

type Post struct {
	ID     		int    			`json:"id"`
	Title  		string 			`json:"title"`
	Slug 		string 			`json:"slug"`
	Body   		string 			`json:"body"`
	UserID 		int 			`json:"userId"`
	CreatedAt 	time.Time		`json:"createdAt"`
	UpdatedAt 	mysql.NullTime	`json:"updatedAt"`
}

func (p *Post) MarshalJSON() ([]byte, error) {
	// TODO: Find a better way to set updatedAt to nil
	if !p.UpdatedAt.Valid {
		return json.Marshal(struct {
			ID int 						`json:"id"`
			Title string 				`json:"title"`
			Slug string 				`json:"slug"`
			Body string					`json:"body"`
			UserID int					`json:"userId"`
			CreatedAt time.Time			`json:"createdAt"`
			UpdatedAt *mysql.NullTime	`json:"updatedAt"`
		}{p.ID, p.Title, p.Slug, p.Body, p.UserID, p.CreatedAt, nil})
	}

	return json.Marshal(struct {
		ID int 					`json:"id"`
		Title string 			`json:"title"`
		Slug string 			`json:"slug"`
		Body string				`json:"body"`
		UserID int				`json:"userId"`
		CreatedAt time.Time		`json:"createdAt"`
		UpdatedAt time.Time		`json:"updatedAt"`
	}{p.ID, p.Title, p.Slug, p.Body, p.UserID, p.CreatedAt, p.UpdatedAt.Time})
}