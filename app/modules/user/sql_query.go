package user

var selectUsersWithTalesQuery = `
WITH user_subset AS (
    SELECT 
        u.id,
        u.created_at
    FROM "user" u
    WHERE u.created_at BETWEEN $1 AND $2
    ORDER BY %s
    LIMIT $3 OFFSET $4
)
SELECT u.id                      AS user_id,
       u.token_number,
       u.use_trial,
       u.invite_code,
       u.is_payed_tale,
       u.first_name,
       u.last_name,
       u.telegram_username,
       u.created_at              AS user_created_at,
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
       t.created_at              AS tale_created_at
FROM "user" u
         LEFT JOIN tale t ON u.id = t.user_id
WHERE u.created_at BETWEEN $1 AND $2 
  AND u.id IN (SELECT id FROM user_subset)
ORDER BY %s;
`

var selectOneWithTalesQuery = `
SELECT json_build_object(
               'id', u.id,
               'token_number', u.token_number,
               'use_trial', u.use_trial,
               'invite_code', u.invite_code,
               'is_payed_tale', u.is_payed_tale,
               'first_name', u.first_name,
               'last_name', u.last_name,
               'telegram_username', u.telegram_username,
               'created_at', u.created_at,
               'tales', COALESCE(json_agg(json_build_object(
                'id', t.id,
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
                                          )) FILTER (WHERE t.id IS NOT NULL), '[]')
       ) AS user
FROM "user" u
         LEFT JOIN
     "tale" t ON u.id = t.user_id
WHERE u.id = $1
GROUP BY u.id;

`
