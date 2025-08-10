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
    hiscore_names text NOT NULL
);

INSERT INTO activity (name, type, hiscore_names) VALUES
('COLO', 'boss', 'Sol Heredit'),
('Corp', 'boss', 'Corporeal Beast'),
('Wildy', 'boss', 'Artio,Callisto,Cal''varion,Vet''ion,Venenatis,Spindel'),
('Callisto', 'boss', 'Callisto,Artio'),
('Vet''ion', 'boss', 'Vet''ion,Cal''varion'),
('Venenatis', 'boss', 'Venenatis,Spindel'),
('COX', 'boss', 'Chambers of Xeric,Chambers of Xeric: Challenge Mode'),
('Chambers of Xeric REG', 'boss', 'Chambers of Xeric'),
('Chambers of Xeric: CM', 'boss', 'Chambers of Xeric: Challenge Mode'),
('Huey', 'boss', 'The Hueycoatl'),
('Inferno', 'boss', 'TzKal-Zuk'),
('Nex', 'boss', 'Nex'),
('NM', 'boss', 'Nightmare,Phosani''s Nightmare'),
('Sarachnis', 'boss', 'Sarachnis'),
('TOA', 'boss', 'Tombs of Amascut,Tombs of Amascut: Expert Mode'),
('Tombs of Amascut REG', 'boss', 'Tombs of Amascut'),
('Tombs of Amascut: Expert', 'boss', 'Tombs of Amascut: Expert Mode'),
('TOB', 'boss', 'Theatre of Blood,Theatre of Blood: Hard Mode'),
('Theatre of Blood REG', 'boss', 'Theatre of Blood'),
('Theatre of Blood: HM', 'boss', 'Theatre of Blood: Hard Mode'),
('Zulrah', 'boss', 'Zulrah'),
('DT2', 'boss', 'Vardorvis,Duke Sucellus,The Whisperer,The Leviathan'),
('Vardorvis', 'boss', 'Vardorvis'),
('Duke Sucellus', 'boss', 'Duke Sucellus'),
('The Whisperer', 'boss', 'The Whisperer'),
('The Leviathan', 'boss', 'The Leviathan'),
('MOKHA', 'boss', 'Doom of Mokhaiotl');

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
    participant_id text NOT NULL REFERENCES participant(discord_id) ON DELETE CASCADE,
    botm_id integer NOT NULL REFERENCES botm(id) ON DELETE CASCADE,
    start_amount integer NOT NULL,
    current_amount integer NOT NULL,
    PRIMARY KEY (participant_id, botm_id)
);

CREATE INDEX IF NOT EXISTS idx_botm_part_botm ON botm_participation(botm_id);
CREATE INDEX IF NOT EXISTS idx_botm_part_participant ON botm_participation(participant_id);

CREATE TABLE IF NOT EXISTS kots_participation(
    participant_id text NOT NULL REFERENCES participant(discord_id) ON DELETE CASCADE,
    kots_id integer NOT NULL REFERENCES kots(id) ON DELETE CASCADE,
    start_amount integer NOT NULL,
    current_amount integer NOT NULL,
    PRIMARY KEY (participant_id, kots_id)
);

CREATE INDEX IF NOT EXISTS idx_kots_part_kots ON kots_participation(kots_id);
CREATE INDEX IF NOT EXISTS idx_kots_part_participant ON kots_participation(participant_id);