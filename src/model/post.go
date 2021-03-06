package model

type Post struct {
	ID                *int                   `json:"id"`
	Title             *string                `json:"title"`
	Description       *string                `json:"description"`
	UserID            *string                `json:"user_id"`
	Slug              *string                `json:"slug"`
	Views             *int                   `json:"views"`
	SocialDescription *string                `json:"socialDescription"`
	SocialTitle       *string                `json:"socialTitle"`
	SocialImage       *string                `json:"socialImage"`
	Meta              map[string]interface{} `json:"meta"`
	Author            *string                 `json:"author"`
	Image             *string                `json:"image"`
	Body              *string                `json:"body"`
	Published         *bool                  `json:"published"`
	CreatedAt         *string                `json:"createdAt"`
	UpdatedAt         *string                `json:"updatedAt"`
}
