package dtos

import "time"

type StationInput struct {
	Origin  string `form:"origin" json:"origin" example:"Jakarta"`
	Name    string `form:"name" json:"name" example:"Pasar Senen"`
	Initial string `form:"initial" json:"initial" example:"PSE"`
}

type StationResponse struct {
	StationID uint      `json:"station_id" example:"1"`
	Origin    string    `json:"origin" example:"Jakarta"`
	Name      string    `json:"name" example:"Pasar Senen"`
	Initial   string    `json:"initial" example:"PSE"`
	CreatedAt time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
	DeletedAt *string   `json:"deleted_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
}

type StationResponseSimply struct {
	StationID  uint   `json:"station_id" example:"1"`
	Origin     string `json:"origin" example:"Jakarta"`
	Name       string `json:"name" example:"Pasar Senen"`
	Initial    string `json:"initial" example:"PSE"`
	ArriveTime string `json:"arrive_time" example:"00:00"`
}
