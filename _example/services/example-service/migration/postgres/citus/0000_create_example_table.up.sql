-- Create new example table

CREATE TABLE example (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    age INTEGER,
    email VARCHAR(255) UNIQUE,
    tags JSONB,
    active BOOLEAN
);