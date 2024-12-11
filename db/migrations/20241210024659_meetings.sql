-- +goose Up
-- +goose StatementBegin
CREATE TABLE booker_details (
    id SERIAL PRIMARY KEY ,
    meeting_id int,
    email VARCHAR(64),
    name VARCHAR(32)
);

CREATE TABLE meetings (
    id SERIAL PRIMARY KEY,
    uid int NOT NULL ,
    date date NOT NULL ,
    start_time time NOT NULL ,
    end_time time NOT NULL ,
    event_id int,
    booker_id int,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE,
    CONSTRAINT fk_booker FOREIGN KEY (booker_id) REFERENCES booker_details(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meetings;
DROP TABLE booker_details;
-- +goose StatementEnd
