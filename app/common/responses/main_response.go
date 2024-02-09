package responses

import "go-fiber-api-template/app/common/types"

type MainResponse struct {
	Status          int                         `json:"status" bson:"status" binding:"required"`
	Data            any                         `json:"data,omitempty" bson:"data,omitempty"`
	ErrorTypeCode   int                         `json:"errorTypeCode,omitempty" bson:"errorTypeCode,omitempty"`
	ValidationError []types.ValidationErrorType `json:"validationError,omitempty" bson:"validationError,omitempty"`
}
