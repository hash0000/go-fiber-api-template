package schema

type UpdateUserSchema struct {
	UserId                  int64   `json:"user_id" validate:"required"`
	FirstName               *string `json:"first_name" validate:"omitempty"`
	FirstNamePresent        bool    `json:"first_name_present"`
	LastName                *string `json:"last_name" validate:"omitempty"`
	LastNamePresent         bool    `json:"last_name_present"`
	TelegramUsername        *string `json:"telegram_username" validate:"omitempty"`
	TelegramUsernamePresent bool    `json:"telegram_username_present"`
}
