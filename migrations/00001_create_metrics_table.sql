-- +goose Up
-- +goose StatementBegin 
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'mtype') THEN
        CREATE TYPE mtype AS ENUM ('gauge','counter');
    END IF;
END $$;
-- +goose StatementEnd 
CREATE TABLE IF NOT EXISTS metrics (
    id VARCHAR NOT NULL,
    mtype mtype NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION
);

CREATE unique INDEX IF NOT EXISTS metrics_mname_idx ON metrics (id, mtype);

-- +goose Down
DROP TABLE metrics;
DROP TYPE mtype;