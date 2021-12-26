package model

import "time"

type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Views     int    `json:"views"`
	Image     string `json:"image"`
	Body      string `json:"body"`
	Published bool   `json:"published"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
