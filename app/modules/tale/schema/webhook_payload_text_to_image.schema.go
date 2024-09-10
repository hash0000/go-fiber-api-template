package schema

type WebhookPayloadTextToImageSchema struct {
	OrderID     int      `json:"order_id" binding:"required" validate:"required"`
	Status      string   `json:"status" binding:"required" validate:"required"`
	Images      []string `json:"images" binding:"required" validate:"required"`
	SecondsTook int      `json:"seconds_took" binding:"required" validate:"required"`
	CustomData  string   `json:"custom_data" binding:"required" validate:"required"`
}

type WebhookPayloadTextToImageCustomDataSchema struct {
	SessionID string `json:"session_id" binding:"required" validate:"required"`
	Index     string `json:"index" binding:"required" validate:"required"`
}
