package dtos

type TrainStationInput struct {
	StationID  uint   `json:"station_id" form:"station_id" example:"1"`
	ArriveTime string `json:"arrive_time" form:"arrive_time" example:"00:00"`
}

type TrainStationResponse struct {
	StationID  uint         `json:"station_id" example:"1"`
	Station    StationInput `json:"station"`
	ArriveTime string       `json:"arrive_time"`
}
