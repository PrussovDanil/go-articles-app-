package models

import "time"

type Article struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	AuthorID  int       `db:"author_id"`
	Published bool      `db:"published"`
	Views     int       `db:"views"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
