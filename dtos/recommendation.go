package dtos

import "time"

type RecommendationInput struct {
	Tag string `json:"tag" form:"tag"`
}

type RecommendationResponse struct {
	RecommendationID uint      `json:"recommendation_id"`
	Tag              string    `json:"tag"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
