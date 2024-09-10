-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    tale
ADD
    COLUMN child_data VARCHAR(1050) NOT NULL DEFAULT 'no data',
ADD
    COLUMN background_characters VARCHAR(1050) NOT NULL DEFAULT 'no data',
ADD
    COLUMN preferences VARCHAR(1050) NOT NULL DEFAULT 'no data',
ADD
    COLUMN moral VARCHAR(1050) NOT NULL DEFAULT 'no data',
ADD
    COLUMN open_ai_answer JSONB NOT NULL DEFAULT '{
  "title": "NO DATA",
  "questions_about_tale": [
    "NO DATA"
  ],
  "chapters": [
    {
      "title": "NO DATA",
      "text": "NO DATA",
      "pic_generation": "NO DATA"
    },
    {
      "title": "NO DATA",
      "text": "NO DATA",
      "pic_generation": "NO DATA"
    },
    {
      "title": "NO DATA",
      "text": "NO DATA",
      "pic_generation": "NO DATA"
    }
  ]
}',
ADD
    COLUMN fabula_img_to_text_json JSONB;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    tale DROP COLUMN child_data,
    DROP COLUMN background_characters,
    DROP COLUMN preferences,
    DROP COLUMN moral,
    DROP COLUMN open_ai_answer,
    DROP COLUMN fabula_img_to_text_json;

-- +goose StatementEnd