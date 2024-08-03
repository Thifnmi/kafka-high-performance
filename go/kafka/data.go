package kafka

type MSGValue struct {
	EventType string `json:"event_type"`
	Publisher struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"publisher"`
	Payload       struct{} `json:"payload"`
	CorrelationId string   `json:"correlation_id"`
	EventId       string   `json:"event_id"`
	Actor         struct {
		Type        string `json:"type"`
		Id          string `json:"id"`
		Name        string `json:"name"`
		ProjectId   string `json:"project_id"`
		ProjectName string `json:"project_name"`
	} `json:"actor"`
	Timestamp string `json:"timestamp"`
	Resource  struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"resource"`
}