CREATE TABLE IF NOT EXISTS repository (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(10),
    full_name VARCHAR(100),
    description VARCHAR(255)
);

INSERT INTO repository (provider, full_name, description) VALUES ('GitHub', 'quantonganh/ssr', 'Security Scan Result');

CREATE TABLE IF NOT EXISTS scan (
    id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    status VARCHAR(11),
    repository_id SERIAL REFERENCES repository (id) ON UPDATE CASCADE ON DELETE CASCADE,
    findings JSONB,
    queued_at TIMESTAMP,
    scanning_at TIMESTAMP,
    finished_at TIMESTAMP
);

CREATE INDEX repository_id_idx ON scan (repository_id);

CREATE INDEX pagination_idx ON scan (finished_at, id);