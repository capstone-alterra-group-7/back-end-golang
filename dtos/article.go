package dtos

import "time"

type ArticleInput struct {
	Title       string `json:"title" form:"title" example:"Title"`
	Image       string `json:"image" form:"image" example:"image.png"`
	Description string `json:"description" form:"description" example:"Description"`
	Label       string `json:"label" form:"label" example:"Label"`
}

type ArticleResponse struct {
	ArticleID   uint      `json:"article_id" example:"1"`
	Title       string    `json:"title" example:"Title"`
	Image       string    `json:"image" example:"image.png"`
	Description string    `json:"description" example:"Description"`
	Label       string    `json:"label" example:"Label"`
	CreatedAt   time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}
