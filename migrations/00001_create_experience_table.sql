-- +goose Up
CREATE TABLE experiences
(
    id      SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type    BIGINT NOT NULL,
    from    TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    to      TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    level   BIGINT NOT NULL
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS experiences;
-- +goose StatementBegin
-- +goose StatementEnd
