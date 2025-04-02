CREATE TABLE weather_conditions (
    id UUID PRIMARY KEY,
    condition_id INT,
    main VARCHAR(50),
    description VARCHAR(100),
    icon VARCHAR(10)
);
