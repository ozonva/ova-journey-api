-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS journeys (
                              journey_id SERIAL PRIMARY KEY,
                              user_id bigint NOT NULL,
                              address text NOT NULL DEFAULT '',
                              description text NOT NULL DEFAULT '',
                              start_time date NOT NULL,
                              end_time date NOT NULL,
                              is_deleted boolean NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS "journeys.user_id_index" ON "journeys"("user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE journeys;
-- +goose StatementEnd