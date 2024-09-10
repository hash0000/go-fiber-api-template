package schemas

type PaginationSchema struct {
	Page int16 `params:"page" binding:"required" validate:"required"`
}
