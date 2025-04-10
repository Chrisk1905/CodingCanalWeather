-- name: InsertWeatherDatum :one
INSERT INTO weather_data (
    id,
    coordinates_id,
    city_name,
    country,
    temperature,
    feels_like,
    temp_min,
    temp_max,
    pressure,
    humidity,
    sea_level,
    grnd_level,
    visibility,
    wind_speed,
    wind_deg,
    cloudiness,
    timestamp,
    sunrise,
    sunset,
    timezone
) VALUES (
    $1,  -- UUID
    $2,  -- coordinates_id (UUID)
    $3,  -- city_name
    $4,  -- country
    $5,  -- temperature
    $6,  -- feels_like
    $7,  -- temp_min
    $8,  -- temp_max
    $9, -- pressure
    $10, -- humidity
    $11, -- sea_level
    $12, -- grnd_level
    $13, -- visibility
    $14, -- wind_speed
    $15, -- wind_deg
    $16, -- cloudiness
    $17, -- timestamp
    $18, -- sunrise
    $19, -- sunset
    $20  -- timezone
)
RETURNING *;
