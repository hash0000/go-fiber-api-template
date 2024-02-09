package entities

import "github.com/google/uuid"

type User struct {
	ID        *uuid.UUID `json:"id,omitempty" sql:"primary_key"`
	Name      *string    `json:"name,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	CreatedAt *string    `json:"createdAt,omitempty"`
}
