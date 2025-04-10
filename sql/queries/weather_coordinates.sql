-- name: InsertWeatherCoordinates :one
INSERT INTO coordinates(id, lon, lat)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetWeatherCoord :one
SELECT * FROM coordinates 
WHERE lon = $1 AND lat = $2;

-- name: GetWeatherByID :one
SELECT * FROM coordinates WHERE id = $1;