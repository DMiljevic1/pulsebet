\connect gameservice_db;

CREATE TABLE IF NOT EXISTS matches (
    id   UUID PRIMARY KEY,
    home TEXT NOT NULL,
    away TEXT NOT NULL
);