package types

type ValidationErrorType struct {
	Property string `json:"property,omitempty"`
	Type     string `json:"type,omitempty"`
	Message  string `json:"message,omitempty"`
}
