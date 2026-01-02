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

CREATE TABLE IF NOT EXISTS booking_services
(
    booking_id      INTEGER NOT NULL,
    service_id      INTEGER NOT NULL,
    quantity        INTEGER NOT NULL DEFAULT 1,
    price           BIGINT NOT NULL,
    PRIMARY KEY (booking_id, service_id),
    FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_booking_services_booking_id ON booking_services(booking_id);
CREATE INDEX IF NOT EXISTS idx_booking_services_service_id ON booking_services(service_id);