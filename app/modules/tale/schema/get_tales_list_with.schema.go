package schema

type GetTalesListSchema struct {
	SortBy   string `query:"sort_by" binding:"required" validate:"required,oneof='DESC' 'ASC'"`
	OrderBy  string `query:"order_by" binding:"required" validate:"required,oneof='created_at'"`
	DateFrom string `query:"date_from" binding:"required" validate:"required"`
	DateTo   string `query:"date_to" binding:"required" validate:"required"`
	Limit    int16  `query:"limit" binding:"required" validate:"required"`
}
