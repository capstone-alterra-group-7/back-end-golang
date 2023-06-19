package dtos

type HistorySearchHotelInput struct {
	Name    string `form:"name" json:"name" binding:"required"`
	HotelID uint   `form:"hotel_id" json:"hotel_id"`
}

type HistorySearchHotelResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	HotelID   uint   `json:"hotel_id"`
	HotelName string `json:"hotel_name"`
}
