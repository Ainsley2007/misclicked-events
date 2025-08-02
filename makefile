# Makefile at project root

# Path to your migrations directory
MIGRATE_PATH := db/migrations

# SQLite connection string (data.db lives in project root)
DB_URL := sqlite3://./data.db

# The migrate CLI (must be in your $PATH)
MIGRATE := migrate

.PHONY: help migrate-up migrate-down migrate-force migrate-status migrate-create

help:
	@echo "Usage:"
	@echo "  make migrate-up        # apply all up migrations"
	@echo "  make migrate-down      # rollback the last migration"
	@echo "  make migrate-force     # force schema version (prompt)"
	@echo "  make migrate-status    # show current migration version"
	@echo "  make migrate-create    # scaffold new SQL migration"

migrate-up:
	@echo "→ Applying migrations..."
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	@echo "→ Rolling back the last migration..."
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

migrate-status:
	@echo "→ Migration status:"
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DB_URL)" version

migrate-force:
	@read -p "Enter target version: " version; \
	echo "→ Forcing to $$version..."; \
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DB_URL)" force $$version

migrate-create:
	@read -p "Enter migration name (e.g. add_users_table): " name; \
	echo "→ Creating new migration: $$name"; \
	$(MIGRATE) create -ext sql -dir $(MIGRATE_PATH) $$name
