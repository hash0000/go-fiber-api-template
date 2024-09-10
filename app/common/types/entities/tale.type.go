package entities

import (
	"encoding/json"
	"time"
)

type TaleType struct {
	ID                  *int64           `json:"id,omitempty"`
	Name                *string          `json:"name,omitempty"`
	UserID              *int64           `json:"user_id,omitempty"`
	FileName            *string          `json:"file_name,omitempty"`
	IsPayed             *bool            `json:"is_payed,omitempty"`
	TaleGenerationID    *string          `json:"tale_generation_id,omitempty"`
	ChildData           *string          `json:"child_data,omitempty"`
	BackgroundCharacter *string          `json:"background_characters,omitempty"`
	Preferences         *string          `json:"preferences,omitempty"`
	Moral               *string          `json:"moral,omitempty"`
	OpenAiAnswer        *json.RawMessage `json:"open_ai_answer,omitempty"`
	FabulaImgToTextJson *json.RawMessage `json:"fabula_img_to_text_json,omitempty"`
	CreatedAt           *time.Time       `json:"created_at,omitempty"`
}

type TalesListWithPaginationType struct {
	Count int64      `json:"count,omitempty"`
	Data  []TaleType `json:"data"`
}
