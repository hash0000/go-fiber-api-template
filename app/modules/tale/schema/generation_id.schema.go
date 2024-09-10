package schema

type GenerationIdSchema struct {
	GenerationId *string `query:"generation_id" validate:"omitempty"`
}
