-- +goose Up
-- +goose StatementBegin
CREATE TABLE meetings (
    id SERIAL PRIMARY KEY,
    uid int NOT NULL ,
    date date NOT NULL ,
    start_time time NOT NULL ,
    end_time time NOT NULL ,
    event_id int,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meetings;
-- +goose StatementEnd
