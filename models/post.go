package models

import (
	"time"
	"github.com/go-sql-driver/mysql"
)

type Post struct {
	ID     		int    		`json:"id"`
	Title  		string 		`json:"title"`
	Slug 		string 		`json:"slug"`
	Body   		string 		`json:"email"`
	UserID 		int 		`json:"user_id"`
	CreatedAt 	time.Time	`json:"createdAt"`
	UpdatedAt 	mysql.NullTime	`json:"updatedAt"`
}