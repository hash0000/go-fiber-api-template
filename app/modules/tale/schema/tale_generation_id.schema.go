package schema

type FinishTrialTaleSchema struct {
	TaleGenerationId string `json:"tale_generation_id" binding:"required" validate:"required"`
	UserId           int64  `json:"user_id" binding:"required" validate:"required"`
}
