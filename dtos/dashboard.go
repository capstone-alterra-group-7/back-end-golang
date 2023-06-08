package dtos

type DashboardResponse struct {
	CountUser  interface{}              `json:"count_user"`
	CountTrain interface{}              `json:"count_train"`
	CountOrder interface{}              `json:"count_order"`
	NewOrder   []map[string]interface{} `json:"new_order"`
	NewUser    []map[string]interface{} `json:"new_user"`
}
