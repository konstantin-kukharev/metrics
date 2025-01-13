-- +goose Up
CREATE TYPE mtype AS ENUM ('gauge','counter');

-- +goose Down
IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'mtype') THEN
    DROP TYPE mtype;
END IF;