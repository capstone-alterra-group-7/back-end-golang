package dtos

type DashboardResponse struct {
	CountUser  interface{}              `json:"count_user"`
	CountHotel interface{}              `json:"count_hotel"`
	CountTrain interface{}              `json:"count_train"`
	CountOrder interface{}              `json:"count_order"`
	NewOrder   []map[string]interface{} `json:"new_order"`
	NewUser    []map[string]interface{} `json:"new_user"`
}
