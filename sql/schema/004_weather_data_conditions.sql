-- +goose Up
CREATE TABLE weather_data_conditions (
    weather_data_id UUID REFERENCES weather_data(id) ON DELETE CASCADE,
    condition_id INT REFERENCES weather_conditions(condition_id) ON DELETE CASCADE,
    PRIMARY KEY (weather_data_id, condition_id)
);

-- +goose Down
DROP TABLE weather_data_conditions;