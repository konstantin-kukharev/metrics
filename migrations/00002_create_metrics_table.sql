-- +goose Up
CREATE TABLE metrics (
    id VARCHAR NOT NULL,
    mtype mtype NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION
);

CREATE unique INDEX IF NOT EXISTS metrics_mname_idx ON metrics (id, mtype);

-- +goose Down
DROP TABLE IF EXISTS metrics;