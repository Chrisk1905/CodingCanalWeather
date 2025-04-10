-- name: InsertWeatherDataConditions :one
INSERT INTO weather_data_conditions(
    weather_data_id,
    condition_id
)
VALUES(
    $1,
    $2
)
RETURNING *;