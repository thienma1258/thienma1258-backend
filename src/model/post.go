package model

type Post struct {
	ID        *int    `json:"id"`
	Title     *string `json:"title"`
	UserID    *string `json:"user_id"`
	Slug      *string `json:"slug"`
	Views     *int    `json:"views"`
	Image     *string `json:"image"`
	Body      *string `json:"body"`
	Published *bool   `json:"published"`
	CreatedAt *string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
}
