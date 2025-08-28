CREATE TABLE IF NOT EXISTS SERVER (
    id text PRIMARY KEY,
    name text NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS activity (
    id integer PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL UNIQUE,
    type text NOT NULL,
    hiscore_names text NOT NULL,
    threshold integer NOT NULL DEFAULT 5
);

INSERT INTO activity (name, type, hiscore_names, threshold) VALUES
('COLO', 'boss', 'Sol Heredit', 1),
('Corp', 'boss', 'Corporeal Beast', 25),
('Wildy', 'boss', 'Artio,Callisto,Cal''varion,Vet''ion,Venenatis,Spindel', 25),
('Callisto', 'boss', 'Callisto,Artio', 25),
('Vet''ion', 'boss', 'Vet''ion,Cal''varion', 25),
('Venenatis', 'boss', 'Venenatis,Spindel', 25),
('COX', 'boss', 'Chambers of Xeric,Chambers of Xeric: Challenge Mode', 5),
('Chambers of Xeric REG', 'boss', 'Chambers of Xeric', 5),
('Chambers of Xeric: CM', 'boss', 'Chambers of Xeric: Challenge Mode', 5),
('Huey', 'boss', 'The Hueycoatl', 25),
('Inferno', 'boss', 'TzKal-Zuk', 1),
('Nex', 'boss', 'Nex', 25),
('NM', 'boss', 'Nightmare,Phosani''s Nightmare', 25),
('Sarachnis', 'boss', 'Sarachnis', 25),
('TOA', 'boss', 'Tombs of Amascut,Tombs of Amascut: Expert Mode', 5),
('Tombs of Amascut REG', 'boss', 'Tombs of Amascut', 5),
('Tombs of Amascut: Expert', 'boss', 'Tombs of Amascut: Expert Mode', 5),
('TOB', 'boss', 'Theatre of Blood,Theatre of Blood: Hard Mode', 5),
('Theatre of Blood REG', 'boss', 'Theatre of Blood', 5),
('Theatre of Blood: HM', 'boss', 'Theatre of Blood: Hard Mode', 5),
('Zulrah', 'boss', 'Zulrah', 25),
('DT2', 'boss', 'Vardorvis,Duke Sucellus,The Whisperer,The Leviathan', 25),
('Vardorvis', 'boss', 'Vardorvis', 25),
('Duke Sucellus', 'boss', 'Duke Sucellus', 25),
('The Whisperer', 'boss', 'The Whisperer', 25),
('The Leviathan', 'boss', 'The Leviathan', 25),
('MOKHA', 'boss', 'Doom of Mokhaiotl', 1);

CREATE TABLE IF NOT EXISTS botm(
    id integer PRIMARY KEY AUTOINCREMENT,
    server_id text NOT NULL REFERENCES SERVER (id) ON DELETE CASCADE,
    activity_id integer NOT NULL REFERENCES activity (id) ON DELETE CASCADE,
    password TEXT NOT NULL,
    status text NOT NULL DEFAULT 'pending'
);

CREATE INDEX IF NOT EXISTS idx_botm_server_id ON botm(server_id);
CREATE INDEX IF NOT EXISTS idx_botm_activity_id ON botm(activity_id);

CREATE TABLE IF NOT EXISTS kots(
    id integer PRIMARY KEY AUTOINCREMENT,
    server_id text NOT NULL REFERENCES SERVER (id) ON DELETE CASCADE,
    activity_id integer NOT NULL REFERENCES activity (id) ON DELETE CASCADE,
    current_king_participant text NOT NULL,
    streak integer NOT NULL DEFAULT 0,
    start_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date DATETIME,
    status text NOT NULL DEFAULT 'pending',
    FOREIGN KEY (current_king_participant) REFERENCES participant(discord_id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_kots_server_id ON kots(server_id);
CREATE INDEX IF NOT EXISTS idx_kots_activity_id ON kots(activity_id);
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
    discord_id text PRIMARY KEY,
    server_id text NOT NULL REFERENCES SERVER (id) ON DELETE CASCADE,
    botm_points integer NOT NULL DEFAULT 0,
    kots_points integer NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_participant_server_id ON participant(server_id);

CREATE TABLE IF NOT EXISTS account(
    id integer PRIMARY KEY AUTOINCREMENT,
    participant_id text NOT NULL REFERENCES participant(discord_id) ON DELETE CASCADE,
    username text NOT NULL,
    failed_fetch_count integer NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_account_participant_id ON account(participant_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_account_participant_username_lower ON account(participant_id, LOWER(username));

CREATE TABLE IF NOT EXISTS botm_participation(
    account_id integer NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    botm_id integer NOT NULL REFERENCES botm(id) ON DELETE CASCADE,
    start_amount integer NOT NULL,
    current_amount integer NOT NULL,
    PRIMARY KEY (account_id, botm_id)
);

CREATE INDEX IF NOT EXISTS idx_botm_part_botm ON botm_participation(botm_id);
CREATE INDEX IF NOT EXISTS idx_botm_part_account ON botm_participation(account_id);

CREATE TABLE IF NOT EXISTS kots_participation(
    account_id integer NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    kots_id integer NOT NULL REFERENCES kots(id) ON DELETE CASCADE,
    start_amount integer NOT NULL,
    current_amount integer NOT NULL,
    PRIMARY KEY (account_id, kots_id)
);

CREATE INDEX IF NOT EXISTS idx_kots_part_kots ON kots_participation(kots_id);
CREATE INDEX IF NOT EXISTS idx_kots_part_account ON kots_participation(account_id);