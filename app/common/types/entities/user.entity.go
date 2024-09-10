package entities

import (
	"github.com/google/uuid"
	"time"
)

type UserUpdateableRowType struct {
	ID          int64     `json:"id,omitempty" sql:"primary_key"`
	TokenNumber int16     `json:"token_number,omitempty"`
	UseTrial    bool      `json:"use_trial,omitempty"`
	InviteCode  uuid.UUID `json:"invite_code,omitempty"`
	CreatedAt   string    `json:"createdAt,omitempty"`
}

type User struct {
	ID               *int64     `json:"id,omitempty" sql:"primary_key"`
	TokenNumber      *int16     `json:"token_number,omitempty"`
	UseTrial         *bool      `json:"use_trial,omitempty"`
	InviteCode       *uuid.UUID `json:"invite_code,omitempty"`
	IsPayedTale      *bool      `json:"is_payed_tale,omitempty"`
	FirstName        *string    `json:"first_name,omitempty"`
	LastName         *string    `json:"last_name,omitempty"`
	TelegramUsername *string    `json:"telegram_username,omitempty"`
	CreatedAt        *time.Time `json:"created_at"`
}
