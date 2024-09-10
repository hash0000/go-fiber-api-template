package tale

var selectTalesList = `
SELECT
       t.id                      AS tale_id,
       t.name                    AS tale_name,
       t.file_name               AS tale_file_name,
       t.is_payed                AS tale_is_payed,
       t.tale_generation_id      AS tale_generation_id,
       t.child_data              AS tale_child_data,
       t.background_characters   AS tale_background_characters,
       t.preferences             AS tale_preferences,
       t.moral                   AS tale_moral,
       t.open_ai_answer          AS tale_open_ai_answer,
       t.fabula_img_to_text_json AS tale_fabula_img_to_text_json,
	   t.user_id				 AS tale_user_ud,
       t.created_at              AS tale_created_at
FROM tale t
WHERE t.created_at BETWEEN $1 AND $2 
ORDER BY %s
LIMIT $3 OFFSET $4;
`

var selectOneQ = `
SELECT json_build_object('id', t.id,
                'name', t.name,
                'file_name', t.file_name,
                'is_payed', t.is_payed,
                'tale_generation_id', t.tale_generation_id,
                'created_at', t.created_at,
                'child_data', t.child_data,
                'background_characters', t.background_characters,
                'preferences', t.preferences,
                'moral', t.moral,
                'open_ai_answer', t.open_ai_answer,
                'fabula_img_to_text_json', t.fabula_img_to_text_json
                                          )
FROM "tale" t
WHERE t.id = $1
LIMIT 1;
`
