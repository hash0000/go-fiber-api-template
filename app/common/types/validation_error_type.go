package types

type ValidationErrorType struct {
	Property string `json:"status" binding:"required"`
	Code     string `json:"code" binding:"required"`
}
