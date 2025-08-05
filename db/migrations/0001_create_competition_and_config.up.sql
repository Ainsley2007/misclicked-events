CREATE TABLE IF NOT EXISTS SERVER (
    id text PRIMARY KEY,
    name text NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS botm(
    id integer PRIMARY KEY AUTOINCREMENT,
    server_id text NOT NULL REFERENCES SERVER (id) ON DELETE CASCADE,
    current_boss text NOT NULL,
    password TEXT NOT NULL,
    status text NOT NULL DEFAULT 'pending'
);

CREATE INDEX IF NOT EXISTS idx_botm_server_id ON botm(server_id);

CREATE TABLE IF NOT EXISTS kots(
    id integer PRIMARY KEY AUTOINCREMENT,
    server_id text NOT NULL REFERENCES SERVER (id) ON DELETE CASCADE,
    current_skill text NOT NULL,
    current_king_participant integer NOT NULL, -- FK → participant(id)
    streak integer NOT NULL DEFAULT 0,
    start_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date DATETIME,
    status text NOT NULL DEFAULT 'pending',
    FOREIGN KEY (current_king_participant) REFERENCES participant(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_kots_server_id ON kots(server_id);

CREATE INDEX IF NOT EXISTS idx_kots_king_participant ON kots(current_king_participant);

CREATE TABLE IF NOT EXISTS config(
    server_id text PRIMARY KEY REFERENCES SERVER (id) ON DELETE CASCADE,
    ranking_channel_id text,
    hiscore_channel_id text,
    category_channel_id text,
    ranking_message_id text,
    hiscore_message_id text
);

CREATE INDEX IF NOT EXISTS idx_config_server_id ON config(server_id);

CREATE TABLE IF NOT EXISTS participant(
    id integer PRIMARY KEY AUTOINCREMENT,
    server_id text NOT NULL REFERENCES SERVER (id) ON DELETE CASCADE,
    discord_id text NOT NULL,
    botm_enabled boolean NOT NULL DEFAULT FALSE,
    kots_enabled boolean NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_participant_server_id ON participant(server_id);

CREATE TABLE IF NOT EXISTS account(
    id integer PRIMARY KEY AUTOINCREMENT,
    participant_id integer NOT NULL REFERENCES participant(id) ON DELETE CASCADE,
    username text NOT NULL,
    failed_fetch_count integer NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_account_participant_id ON account(participant_id);

CREATE TABLE IF NOT EXISTS botm_participation(
    account_id integer NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    botm_id integer NOT NULL REFERENCES botm(id) ON DELETE CASCADE,
    starting_kc integer NOT NULL,
    current_kc integer NOT NULL,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (account_id, botm_id)
);

CREATE INDEX IF NOT EXISTS idx_botm_part_botm ON botm_participation(botm_id);

CREATE INDEX IF NOT EXISTS idx_botm_part_acc ON botm_participation(account_id);

CREATE TABLE IF NOT EXISTS kots_participation(
    account_id integer NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    kots_id integer NOT NULL REFERENCES kots(id) ON DELETE CASCADE,
    starting_xp integer NOT NULL,
    current_xp integer NOT NULL,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (account_id, kots_id)
);

CREATE INDEX IF NOT EXISTS idx_kots_part_kots ON kots_participation(kots_id);

CREATE INDEX IF NOT EXISTS idx_kots_part_acc ON kots_participation(account_id);

