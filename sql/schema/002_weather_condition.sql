-- +goose Up
CREATE TABLE weather_conditions (
    id UUID PRIMARY KEY,
    condition_id INT UNIQUE,
    main VARCHAR(50),
    description VARCHAR(100),
    icon VARCHAR(10)
);

-- +goose Down
DROP TABLE weather_conditions;