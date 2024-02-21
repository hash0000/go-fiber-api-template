package schema

type RedirectSchema struct {
	Range string `json:"range" binding:"required"`
	Gid   string `json:"gid" binding:"required"`
}
