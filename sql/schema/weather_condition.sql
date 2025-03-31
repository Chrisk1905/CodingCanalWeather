CREATE TABLE weather_conditions (
    id UUID PRIMARY KEY,
    weather_id INT REFERENCES weather_data(id) ON DELETE CASCADE,
    condition_id INT,
    main VARCHAR(50),
    description VARCHAR(100),
    icon VARCHAR(10)
);
