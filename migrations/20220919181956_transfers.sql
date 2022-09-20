-- +goose Up
-- +goose StatementBegin

CREATE TABLE transfers (
    id TEXT PRIMARY KEY default gen_random_uuid(),
    type TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    receiver_id TEXT NOT NULL,
    amount INTEGER NOT NULL,
    description TEXT NOT NULL,
    created_at timestamptz default now()
);

CREATE TABLE entry_parts (
    id TEXT PRIMARY KEY default gen_random_uuid(),
    transfer_id TEXT NOT NULL,
    type INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    amount INTEGER NOT NULL,
    created_at timestamptz default now(),

    FOREIGN KEY (transfer_id) REFERENCES transfers(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE entry_parts;

DROP TABLE transfers;

-- +goose StatementEnd
