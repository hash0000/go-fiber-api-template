package constants

const (
	GptJsonOption = `
{
		"type": "json_schema",
		"json_schema": {
			"name": "math_response",
			"strict": true,
			"schema": {
				"type": "object",
				"properties": {
					"title": {
						"type": "string"
					},
					"questions_about_tale": {
						"type": "array",
						"items": {
							"type": "string"
						}
					},
					"chapters": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"title": {
									"type": "string"
								},
								"text": {
									"type": "string"
								},
								"pic_generation": {
									"type": "string"
								}
							},
							"required": [
								"title",
								"text",
								"pic_generation"
							],
							"additionalProperties": false
						}
					}
				},
				"required": [
					"title",
					"questions_about_tale",
					"chapters"
				],
				"additionalProperties": false
			}
		}
	}`

	GptSystemPromptV1 = ``
	GptPrompt         = ``
)
