package responses

import (
	"fmt"
	"go-fiber-api-template/app/common/types"
)

type MainResponse struct {
	Status        int                         `json:"status" bson:"status" binding:"required"`
	Data          any                         `json:"data,omitempty" bson:"data,omitempty"`
	ErrorTypeCode int                         `json:"error_type_code,omitempty" bson:"error_type_code,omitempty"`
	ErrorInfo     []types.ValidationErrorType `json:"error_info,omitempty" bson:"error_info,omitempty"`
	Message       string                      `json:"message,omitempty" bson:"message,omitempty"`
}

func (mr *MainResponse) Error() string {
	return fmt.Sprintf("Status: %d, Message: %s, ErrorTypeCode: %d", mr.Status, mr.Message, mr.ErrorTypeCode)
}
