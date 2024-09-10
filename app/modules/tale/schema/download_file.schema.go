package schema

type GetFileSchema struct {
	FileName string `params:"file_name" binding:"required" validate:"required"`
}
