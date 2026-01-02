CREATE TABLE IF NOT EXISTS services
(
    id              SERIAL PRIMARY KEY,
    parking_id      INTEGER NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT,
    price           BIGINT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parking_id) REFERENCES parking_places(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_services_parking_id ON services(parking_id);

CREATE TABLE IF NOT EXISTS booking_services
(
    booking_id      INTEGER NOT NULL,
    service_id      INTEGER NOT NULL,
    quantity        INTEGER NOT NULL DEFAULT 1,
    price           BIGINT NOT NULL,
    PRIMARY KEY (booking_id, service_id),
    FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_booking_services_booking_id ON booking_services(booking_id);
CREATE INDEX IF NOT EXISTS idx_booking_services_service_id ON booking_services(service_id);

