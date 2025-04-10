-- +goose Up
CREATE TABLE weather_data (
    id UUID PRIMARY KEY,
    coordinates_id UUID REFERENCES coordinates(id) ON DELETE CASCADE,    
    city_name VARCHAR(50),
    country VARCHAR(50),
    temperature FLOAT,
    feels_like FLOAT,
    temp_min FLOAT,
    temp_max FLOAT,
    pressure INT,
    humidity INT,
    sea_level INT,
    grnd_level INT,
    visibility INT,
    wind_speed FLOAT,
    wind_deg INT,
    cloudiness INT,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    sunrise TIMESTAMPTZ,
    sunset TIMESTAMPTZ,
    timezone INT
);

-- +goose Down
DROP TABLE weather_data;