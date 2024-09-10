package schema

type UpdateTokenNumberSchema struct {
	ID          int64 `json:"id" params:"id" binding:"required" validate:"required"`
	TokenNumber int16 `json:"token_number" binding:"required" validate:"required"`
}
