CREATE TABLE IF NOT EXISTS bookings
(
    id               SERIAL PRIMARY KEY,
    date_from        TIMESTAMP    NOT NULL,
    date_to          TIMESTAMP    NOT NULL,
    parking_place_id INTEGER NOT NULL,
    full_cost        INTEGER                                                                     DEFAULT 0,
    status           TEXT CHECK ( status in ('Waiting', 'Confirmed', 'Canceled') ) DEFAULT 'Waiting',
    user_id          TEXT    NOT NULL
);