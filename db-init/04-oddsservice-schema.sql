\connect oddsservice_db;

CREATE TABLE IF NOT EXISTS odds (
    match_id   UUID PRIMARY KEY,
    home_odds  NUMERIC(6,3) NOT NULL,
    draw_odds  NUMERIC(6,3) NOT NULL,
    away_odds  NUMERIC(6,3) NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);