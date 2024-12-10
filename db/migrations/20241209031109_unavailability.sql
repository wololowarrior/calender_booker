-- +goose Up
-- +goose StatementBegin
CREATE TABLE unavailable_slots(
    id SERIAL PRIMARY KEY,
    uid int,
    unavailable_date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users FOREIGN KEY (uid) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE unavailable_slots;
-- +goose StatementEnd
