package schema

type IdSchema struct {
	ID int64 `json:"id" params:"id" binding:"required" validate:"required"`
}
