-- name: InsertWeatherCoordinates :one
INSERT INTO coordinates(id, lon, lat)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;


-- name: GetWeatherByID :one
SELECT * FROM coordinates WHERE id = $1;