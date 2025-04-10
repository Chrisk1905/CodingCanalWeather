-- name: InsertCondition :one
INSERT INTO weather_conditions (
    id, condition_id, main, description, icon
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetConditionByConditionID :one
SELECT * FROM weather_conditions 
WHERE condition_id = $1;

