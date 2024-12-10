-- +goose Up
-- +goose StatementBegin
CREATE TYPE slots_type as ENUM ('30','60');
CREATE TABLE event (
    id SERIAL PRIMARY KEY,
    uid int NOT NULL,
    name VARCHAR(255) NOT NULL,
    message VARCHAR(255),
    slotted boolean DEFAULT false,
    slots slots_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table event;
DROP TYPE slots_type;
-- +goose StatementEnd
