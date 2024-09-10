package schema

type CreateUserSchema struct {
	UserID     int64  `json:"user_id" binding:"required" validate:"required"`
	InviteCode string `json:"invite_code" validate:"omitempty,uuid"`
}
