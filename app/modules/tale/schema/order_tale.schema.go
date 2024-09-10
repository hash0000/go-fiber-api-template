package schema

type OrderTaleSchema struct {
	ChildData            string `json:"child_data" binding:"required" validate:"required"`
	BackgroundCharacters string `json:"background_characters" binding:"required" validate:"required"`
	Preferences          string `json:"preferences" binding:"required" validate:"required"`
	Moral                string `json:"moral" binding:"required" validate:"required"`
	ID                   int64  `json:"id" binding:"required" validate:"required"`
	Url                  string `json:"url" validate:"omitempty"`
}
