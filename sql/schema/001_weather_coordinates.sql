-- +goose Up
CREATE TABLE coordinates (
    id UUID PRIMARY KEY,
    lon FLOAT,
    lat FLOAT,
    CONSTRAINT unique_coordinates UNIQUE (lon, lat)
);

-- +goose Down
DROP TABLE coordinates;