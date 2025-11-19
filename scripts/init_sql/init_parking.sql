CREATE TABLE IF NOT EXISTS parking_places
(
    id           SERIAL PRIMARY KEY,
    name         TEXT NOT NULL,
    city         TEXT NOT NULL,
    address      TEXT NOT NULL,
    parking_type TEXT CHECK ( parking_type IN ('outdoor', 'covered', 'underground', 'multi-level') ),
    hourly_rate  INT  NOT NULL,
    capacity     INT  NOT NULL DEFAULT 0,
    owner_id     TEXT
);