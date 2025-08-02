-- Create competitions table
CREATE TABLE IF NOT EXISTS competition (
    id           TEXT      PRIMARY KEY,
    name         TEXT      NOT NULL,
    start_date   DATETIME  NOT NULL,
    end_date     DATETIME  NOT NULL,
    status       TEXT      NOT NULL  DEFAULT 'pending',
    created_at   DATETIME  NOT NULL  DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME  NOT NULL  DEFAULT CURRENT_TIMESTAMP
);

-- Create config key/value store
CREATE TABLE IF NOT EXISTS config (
    key          TEXT      PRIMARY KEY,
    value        TEXT      NOT NULL,
    description  TEXT,
    created_at   DATETIME  NOT NULL  DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME  NOT NULL  DEFAULT CURRENT_TIMESTAMP
);
