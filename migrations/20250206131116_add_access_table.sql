-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access (
    access_id bigserial PRIMARY KEY,
    login TEXT NOT NULL,
    data_id TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE access;
-- +goose StatementEnd