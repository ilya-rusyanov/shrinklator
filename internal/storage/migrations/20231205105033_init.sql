-- +goose Up
CREATE TABLE shorts
(short text,
 long text UNIQUE,
 user_id text,
 is_deleted boolean NOT NULL DEFAULT false,
 PRIMARY KEY (short));

-- +goose Down
DROP TABLE shorts;
