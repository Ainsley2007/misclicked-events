## Misclicked Events — Implemented Functionality

This document describes the features currently implemented in the project, how they work end-to-end, and where to find the relevant code.

### Architecture overview
- **Runtime**: Discord bot built with `discordgo`.
- **Entry point**: `cmd/misclickedevents/main.go` — initializes the data layer, creates the Discord session, registers slash commands, and starts background services.
- **Configuration**: `internal/config/config.go` reads the `DISCORD_BOT_TOKEN` from `.env`.
- **Services**: A background hiscore update service periodically updates leaderboards.
- **Domain & Use cases**: Implemented under `internal/domain/usecase` and orchestrated via repositories in `internal/data/repository`.
- **Persistence**: SQLite (`data.db`) with schemas/migrations in `db/migrations/0001_create_competition_and_config.up.sql`.
- **External APIs**: OSRS Hiscores (`https://secure.runescape.com/m=hiscore_oldschool`) for account validation and KC reads.

### Discord commands (slash)
Registered in `internal/commands/commands.go`, handled by `internal/handlers/interactions.go`.

- `/setup-channels` — Configure where results are displayed.
  - Handler: `internal/commands/config_command.go`
  - Stores channel IDs in `config` table, used by the hiscore service to post or update embeds.

- `/add-account` — Link an OSRS account to your Discord user.
  - Handler: `internal/commands/add_account_command.go`
  - Use case: `internal/domain/usecase/add_account.go`
    - Validates player existence against OSRS Hiscores.
    - Persists account under the participant.
    - If a BOTM event is running, takes a starting KC snapshot and adds participation for that account.

- `/remove-account` — Stop tracking a linked OSRS account.
  - Handler: `internal/commands/remove_account_command.go`
  - Repository call: `ParticipantRepo.RemoveAccount`.

- `/tracked-accounts` — View your tracked accounts; if an event is active, shows KC per account and total.
  - Handler: `internal/commands/tracked_accounts_command.go`
  - Efficiently fetches accounts with KC when a BOTM is active.

- `/rename-account` — Rename a tracked account.
  - Handler: `internal/commands/rename_account_command.go`
  - Use case: `RenameAccountUseCase` inits in `internal/data/app.go` and updates storage.

- `/add-activity` — Add an activity (boss/skill) to the catalog, with optional hiscore name mapping and threshold.
  - Handler: `internal/commands/add_activity_command.go`
  - Use case: `internal/domain/usecase/add_activity.go`.

- `/remove-activity` — Remove an activity from the catalog.
  - Handler: `internal/commands/remove_activity_command.go`
  - Use case: `internal/domain/usecase/remove_activity.go`.

- `/start` — Start a BOTM event for a chosen activity (admin only).
  - Handler: `internal/commands/start_activity_command.go`
  - Use case: `internal/domain/usecase/start_activity.go`
    - Validates there is no current event.
    - Loads the activity and its configured hiscore names.
    - Creates a BOTM record with status and password.
    - For each participant, snapshots per-account KC across the mapped bosses and registers participation.
    - Optionally renames the configured category channel to reflect the current event.

Autocomplete is implemented for account and activity selections in `internal/commands/command_helpers.go`.

### Background service: Hiscore updates
File: `internal/services/hiscore_service.go`.

- Starts on bot launch and runs hourly.
- For each guild, executes `UpdateHiscoresUseCase` to:
  - Check if an event is active, list tracked participants, and gather per-account KC deltas.
  - Build a leaderboard with ranks and totals.
  - Post or update an embed in the configured hiscore channel.
- If no event is active, posts/updates a "No Ongoing Event" embed.

Supporting use case: `internal/domain/usecase/update_hiscores.go`.

### Data model and persistence
Storage: SQLite (`data.db`). Initialization and wiring in `internal/data/app.go`.

- Core tables: `server`, `config`, `activity`, `botm`, `participant`, `account`, `botm_participation`, with indices.
- Activities include type, hiscore name mappings (comma-separated), and a KC threshold.
- BOTM stores the current activity and status per server; participation tracks starting and current KC per account.
- Migrations and seed activities are in `db/migrations/0001_create_competition_and_config.up.sql`.

Repositories under `internal/data/repository` encapsulate storage and external API access:
- `CompetitionRepository` — start/stop/fetch BOTM and map activity/domain.
- `ParticipantRepository` — add/remove/rename accounts, list tracked accounts, manage BOTM participation, and fetch per-account KC.
- `HiscoreRepository` — talks to OSRS Hiscores via `internal/data/datasource/api`, maps to domain models, and provides KC queries (single and combined).
- `ConfigRepository`, `ServerRepository`, `ActivityRepository` — configuration and catalog management.

### External hiscore integration
Source: `internal/data/datasource/api/hiscore_datasource.go` and `internal/data/repository/hiscore_repository.go`.

- Validates player existence (`index_lite.ws`).
- Fetches JSON hiscores (`index_lite.json`).
- Maps skills and activities to domain types and exposes helpers to compute combined KC across multiple boss names.

### Bot readiness and server registration
File: `internal/handlers/ready_handler.go` — on ready, registers each guild (server) in the database through `ServerRepo`.

### CLI utility: Bulk import
File: `cmd/bulk_import/main.go`.

- Reads `participants.json` and bulk inserts participants and their accounts into `data.db`.
- Targets the first server found in the database.

### How it fits together (flow)
1. Bot starts (`cmd/misclickedevents/main.go`), initializes SQLite and repositories, loads token, connects to Discord.
2. Registers slash commands and starts the hiscore service.
3. Admin runs `/setup-channels` to configure where to display results.
4. Users link accounts via `/add-account` (validated against OSRS hiscores).
5. Admin starts an event with `/start` for a chosen activity; system snapshots starting KC for all participants.
6. The background service periodically updates leaderboards based on current KC and posts/edits the embed.

### Notes and limitations
- Ending events is scaffolded in handlers but not implemented (commented in `interactions.go`).
- Skilling activities are modeled but the primary flow focuses on bosses/KC.
- The repository includes a seed list of activities/bosses with combined mappings.

### Roadmap — Phased plan

#### Phase 1 — Code audit and cleanup
- Identify unused code per package (`internal/commands`, `handlers`, `domain/usecase`, `data/**`, `services`, `utils`), remove dead/duplicate helpers.
- Tighten exports: unexport functions/types not used outside their package.

#### Phase 2 — Logging overhaul
- Unify logging via `utils` with levels (Error, Info, Debug) and env-driven level.
- Remove verbose prints and chatty debug in hot paths; keep a light audit trail:
  - Info: add/remove/rename account; start/stop BOTM; add/remove activity.
- Keep minimal Debug at key checkpoints; gate by level.

#### Phase 3 — Architecture alignment
- Ensure dependencies direction: handlers → commands → usecases → interfaces → repository/datasource.
- Use cases depend on interfaces, not concrete repos; persistence details stay in repos.
- Keep `internal/data/app.go` as the composition root only.

#### Phase 4 — Hiscore updates: resilience and correctness
- Add bounded retries with jitter in hiscore datasource; handle rate limiting.
- HTTP status handling: 400 → rate-limited (retry/backoff, increment `failed_fetch_count`), 404 → account missing (mark and notify participant).
- Notify participant (Discord DM) once per account per day to use `/rename-account`.
- Serialize or lightly throttle account fetches per guild to avoid bursts.

#### Phase 5 — Generalize the scheduler loop
- Refactor `HiscoreService` into a generic periodic scheduler supporting multiple tasks (hourly, weekly).
- Tasks: hourly BOTM leaderboard refresh; weekly KOTS evaluation.
- Preserve graceful shutdown and task cancelation.

#### Phase 6 — King of the Skill (KOTS)
- Use skill-type activities (`activity.type = 'skill'`).
- Weekly flow: select skill, compute XP deltas per participant’s accounts, determine King.
- Scoring rules:
  - Week 1: 0 points (first crown)
  - Retaining crown: 3 × current streak
  - Dethroning: always 5 points
  - Streak resets when dethroned
- Update points/streaks on `kots` and participants; announce weekly King; simple standings view.

#### Phase 7 — Validation and polish
- Migrations for any new fields (e.g., last-notified for 404s).
- Smoke tests in a test guild; verify commands, rate limit handling, DMs.
- Update this doc and `readme.md` links; add KOTS rules summary.

