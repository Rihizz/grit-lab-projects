CREATE TABLE event_responses (
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    response INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id)
);
