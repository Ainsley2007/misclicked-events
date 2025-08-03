CREATE TABLE IF NOT EXISTS SERVER (
    id text PRIMARY KEY,
    name text NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS competition(
    id text PRIMARY KEY,
    server_id text NOT NULL,
    name text NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    status text NOT NULL DEFAULT 'pending',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (server_id) REFERENCES SERVER (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_competition_server_id ON competition(server_id);

CREATE TABLE IF NOT EXISTS config (
  server_id           TEXT      PRIMARY KEY,
  ranking_channel_id  TEXT,
  hiscore_channel_id  TEXT,
  category_channel_id TEXT,
  ranking_message_id  TEXT,
  hiscore_message_id  TEXT,
  created_at          DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at          DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (server_id) REFERENCES server(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_config_server_id ON config(server_id);

CREATE TABLE IF NOT EXISTS participant(
    id integer PRIMARY KEY AUTOINCREMENT,
    server_id text NOT NULL,
    discord_id text NOT NULL,
    discord_name text NOT NULL,
    botm_score integer NOT NULL DEFAULT 0,
    kots_score integer NOT NULL DEFAULT 0,
    botm_enabled boolean NOT NULL DEFAULT FALSE,
    kots_enabled boolean NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (server_id) REFERENCES SERVER (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_participant_server ON participant(server_id);

CREATE TABLE IF NOT EXISTS account(
    id integer PRIMARY KEY AUTOINCREMENT,
    participant_id integer NOT NULL,
    username text NOT NULL,
    botm_starting_kc integer NOT NULL,
    botm_current_kc integer NOT NULL,
    kots_starting_xp integer NOT NULL,
    kots_current_xp integer NOT NULL,
    failed_fetch_count integer NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (participant_id) REFERENCES participant(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_account_participant ON account(participant_id);

