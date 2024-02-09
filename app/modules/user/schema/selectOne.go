package schema

import "github.com/google/uuid"

type SelectOneUserSchema struct {
	Id uuid.UUID `json:"id" binding:"required" validate:"required"`
}
