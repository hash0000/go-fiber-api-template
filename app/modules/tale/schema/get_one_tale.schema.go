package schema

type GetOneTaleSchema struct {
	ID int64 `params:"id" binding:"required" validate:"required"`
}
