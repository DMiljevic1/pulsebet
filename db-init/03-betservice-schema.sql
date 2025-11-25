\connect betservice_db;

CREATE TABLE IF NOT EXISTS matches (
    id   UUID PRIMARY KEY,
    home TEXT NOT NULL,
    away TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS bets (
    id         UUID PRIMARY KEY,
    match_id   UUID NOT NULL,
    bettor_id  TEXT NOT NULL,
    selection  TEXT NOT NULL,
    stake      NUMERIC(12,2) NOT NULL CHECK (stake > 0),
    odds       NUMERIC(6,3)  NOT NULL CHECK (odds > 0),
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now()
);

ALTER TABLE bets
    ADD CONSTRAINT fk_bets_match
    FOREIGN KEY (match_id) REFERENCES matches(id);
