package entities

import "github.com/google/uuid"

type Url struct {
	ID        *uuid.UUID `json:"id,omitempty" sql:"primary_key"`
	Url       string     `json:"url,omitempty"`
	UpdatedAt *string    `json:"updatedAt,omitempty"`
	CreatedAt *string    `json:"createdAt,omitempty"`
}
