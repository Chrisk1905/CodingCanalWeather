CREATE TABLE coordinates (
    id UUID PRIMARY KEY,
    weather_id INT REFERENCES weather_data(id) ON DELETE CASCADE,
    lon FLOAT,
    lat FLOAT
);
