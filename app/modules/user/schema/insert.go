package schema

type InsertUserSchema struct {
	Name  string `json:"name" binding:"required" validate:"required"`
	Phone string `json:"phone" binding:"required" validate:"required"`
}
