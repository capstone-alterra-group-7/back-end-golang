package dtos

import (
	"time"
)

type UserRegisterInput struct {
	FullName        string  `form:"full_name" json:"full_name" example:"Mochammad Hanif"`
	Email           string  `form:"email" json:"email" example:"me@hanifz.com"`
	Password        string  `form:"password" json:"password" example:"qweqwe123"`
	ConfirmPassword string  `form:"confirm_password" json:"confirm_password" example:"qweqwe123"`
	BirthDate       *string `form:"birth_date" json:"birth_date,omitempty" example:"2002-09-12"`
	PhoneNumber     string  `form:"phone_number" json:"phone_number" example:"0851555555151"`
	Role            string  `form:"role" json:"role" example:"user"`
	IsActive        *bool   `form:"is_active" json:"is_active,omitempty" example:"true"`
}

type UserLoginInput struct {
	Email    string `form:"email" json:"email" example:"me@hanifz.com"`
	Password string `form:"password" json:"password" example:"qweqwe123"`
}

type UserUpdatePhotoProfileInput struct {
	ProfilePicture string `form:"file" json:"file" example:"https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"`
}

type UserUpdatePasswordInput struct {
	OldPassword     string `form:"old_password" json:"old_password" example:"qweqwe123"`
	NewPassword     string `form:"new_password" json:"new_password" example:"asdqwe123"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" example:"asdqwe123"`
}

type UserUpdateProfileInput struct {
	FullName    string `form:"full_name" json:"full_name" example:"Hanif Mochammad"`
	PhoneNumber string `form:"phone_number" json:"phone_number" example:"085199999999"`
	BirthDate   string `form:"birth_date" json:"birth_date" example:"2000-01-01"`
	Citizen     string `form:"citizen" json:"citizen" example:"Indonesia"`
}

type UserLoginResponse struct {
	FullName    string    `json:"full_name" example:"Mochammad Hanif"`
	Email       string    `json:"email" example:"me@hanifz.com"`
	PhoneNumber string    `json:"phone_number" example:"0851555555151"`
	Role        string    `json:"role" example:"user"`
	CreatedAt   time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type UserRegisterResponse struct {
	FullName    string    `json:"full_name" example:"Mochammad Hanif"`
	Email       string    `json:"email" example:"me@hanifz.com"`
	PhoneNumber string    `json:"phone_number" example:"0851555555151"`
	Role        string    `json:"role" example:"user"`
	CreatedAt   time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type UserInformationResponse struct {
	ID             uint      `json:"id" example:"1"`
	FullName       string    `json:"full_name" example:"Mochammad Hanif"`
	Email          string    `json:"email" example:"me@hanifz.com"`
	PhoneNumber    string    `json:"phone_number" example:"0851555555151"`
	BirthDate      string    `json:"birth_date" example:"2002-09-12"`
	ProfilePicture string    `json:"profile_picture_url" example:"https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"`
	Citizen        string    `json:"citizen" example:"Indonesia"`
	Role           *string   `json:"role,omitempty" example:"user"`
	Token          *string   `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2ODQ0MDYzMzMsInJvbGUiOiJ1c2VyIiwidXNlcklkIjozfQ.B8vBlMIiU4iZR0YHe4-Mo3DpJ2nwlTV3PuhEJc31pMo"`
	CreatedAt      time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt      time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
	DeletedAt      *string   `json:"deleted_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
}

type UserInformationResponses struct {
	ID             uint   `json:"user_id" example:"1"`
	FullName       string `json:"full_name" example:"Mochammad Hanif"`
	Email          string `json:"email" example:"me@hanifz.com"`
	PhoneNumber    string `json:"phone_number" example:"0851555555151"`
	BirthDate      string `json:"birth_date" example:"2002-09-12"`
	ProfilePicture string `json:"profile_picture_url" example:"https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"`
	Citizen        string `json:"citizen" example:"Indonesia"`
}
