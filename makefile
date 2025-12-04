# ===========================
# movies3 Makefile
# ===========================

# ---- Project & tools ----
PROJECT_NAME := movies3
GO           := go

# ERD generator Go package
ERD_GEN_PKG  := ./cmd/gen_erd
ERD_GEN_CMD  := $(GO) run $(ERD_GEN_PKG)

# Migration tool (old -> new DB reference data)
MIGRATE_OLD_PKG := ./cmd/migrate-old-db
MIGRATE_OLD_CMD := $(GO) run $(MIGRATE_OLD_PKG)


# ---- Database schema paths ----
OLD_SCHEMA   := db/old/schema.sql
NEW_SCHEMA   := db/new/schema.sql

OLD_ERD      := db/old/schema.mmd
NEW_ERD      := db/new/schema.mmd

# ---- Utility ----
.PHONY: help
## Show all make targets with descriptions
help:
	@echo "Available make commands:"
	@grep -E '^[a-zA-Z0-9_-]+:.*## ' $(MAKEFILE_LIST) | sort \
		| awk 'BEGIN {FS=":.*## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ===========================
# Go basics
# ===========================

.PHONY: fmt
fmt: ## Run gofmt on all Go files
	@echo ">> gofmt"
	@find . -name '*.go' -not -path './vendor/*' -exec gofmt -w {} +

.PHONY: tidy
tidy: ## Run go mod tidy
	@echo ">> go mod tidy"
	@$(GO) mod tidy

.PHONY: test
test: ## Run go test ./...
	@echo ">> go test ./..."
	@$(GO) test ./...

.PHONY: build
build: ## Build main application binary (if/when you add it)
	@echo ">> go build ./..."
	@$(GO) build ./...

# ===========================
# ERD generation (old & new)
# ===========================

.PHONY: create-db-dirs
create-db-dirs: ## Ensure db/old and db/new directories exist
	@mkdir -p db/old db/new

.PHONY: build-erd-old
build-erd-old: create-db-dirs ## Generate ERD for the OLD (LabVIEW-era) schema
	@echo ">> Generating OLD ERD from $(OLD_SCHEMA) -> $(OLD_ERD)"
	@$(ERD_GEN_CMD) -sql $(OLD_SCHEMA) -out $(OLD_ERD)

.PHONY: build-erd-new
build-erd-new: create-db-dirs ## Generate ERD for the NEW movies3 schema
	@echo ">> Generating NEW ERD from $(NEW_SCHEMA) -> $(NEW_ERD)"
	@$(ERD_GEN_CMD) -sql $(NEW_SCHEMA) -out $(NEW_ERD)

.PHONY: clean-erd
clean-erd: ## Remove generated .mmd ERD files
	@echo ">> Removing ERD files"
	@rm -f $(OLD_ERD) $(NEW_ERD)

# ===========================
# Dev helpers
# ===========================

.PHONY: check
check: fmt test ## Run formatting and tests
	@echo ">> All checks passed"

# ===========================
# movies3 Makefile
# ===========================

# ---- Project & tools ----
PROJECT_NAME := movies3
GO           := go

# ERD generator Go package
ERD_GEN_PKG  := ./cmd/gen_erd
ERD_GEN_CMD  := $(GO) run $(ERD_GEN_PKG)

# ---- Database schema paths ----
OLD_SCHEMA   := db/old/schema.sql
NEW_SCHEMA   := db/new/schema.sql

OLD_ERD      := db/old/schema.mmd
NEW_ERD      := db/new/schema.mmd

# ---- Utility ----
.PHONY: help
## Show all make targets with descriptions
help:
	@echo "Available make commands:"
	@grep -E '^[a-zA-Z0-9_-]+:.*## ' $(MAKEFILE_LIST) | sort \
		| awk 'BEGIN {FS=":.*## "}; {printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2}'

# ===========================
# Go basics
# ===========================

.PHONY: fmt
fmt: ## Run gofmt on all Go files
	@echo ">> gofmt"
	@find . -name '*.go' -not -path './vendor/*' -exec gofmt -w {} +

.PHONY: tidy
tidy: ## Run go mod tidy
	@echo ">> go mod tidy"
	@$(GO) mod tidy

.PHONY: test
test: ## Run go test ./...
	@echo ">> go test ./..."
	@$(GO) test ./...

.PHONY: build
build: ## Build main application binary (if/when you add it)
	@echo ">> go build ./..."
	@$(GO) build ./...

# ===========================
# ERD generation (old & new)
# ===========================

.PHONY: create-db-dirs
create-db-dirs: ## Ensure db/old and db/new directories exist
	@mkdir -p db/old db/new

.PHONY: build-erd-old
build-erd-old: create-db-dirs ## Generate ERD for the OLD (LabVIEW-era) schema
	@echo ">> Generating OLD ERD from $(OLD_SCHEMA) -> $(OLD_ERD)"
	@$(ERD_GEN_CMD) -sql $(OLD_SCHEMA) -out $(OLD_ERD)

.PHONY: build-erd-new
build-erd-new: create-db-dirs ## Generate ERD for the NEW movies3 schema
	@echo ">> Generating NEW ERD from $(NEW_SCHEMA) -> $(NEW_ERD)"
	@$(ERD_GEN_CMD) -sql $(NEW_SCHEMA) -out $(NEW_ERD)

.PHONY: clean-erd
clean-erd: ## Remove generated .mmd ERD files
	@echo ">> Removing ERD files"
	@rm -f $(OLD_ERD) $(NEW_ERD)

# ===========================
# Old → New DB migration
# ===========================

DB_PASSWORD ?= 12345678
OLD_DB_DSN ?= host=127.0.0.1 user=postgres password=$(DB_PASSWORD) dbname=mediadb sslmode=disable
NEW_DB_DSN ?= host=127.0.0.1 user=postgres password=$(DB_PASSWORD) dbname=movies3db sslmode=disable

MIGRATE_CMD := $(GO) run ./cmd/migrate-old-db

# ---- Reference data (already done) ----

.PHONY: migrate-ref-dry-run
migrate-ref-dry-run: ## DRY-RUN reference migration (no writes)
	@echo ">> DRY-RUN reference migration (no writes)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	@OLD_DB_DSN='$(OLD_DB_DSN)' NEW_DB_DSN='$(NEW_DB_DSN)' MIGRATION_PHASE=refs DRY_RUN=1 $(MIGRATE_CMD)

.PHONY: migrate-ref
migrate-ref: ## REAL reference migration
	@echo ">> REAL reference migration"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	@OLD_DB_DSN='$(OLD_DB_DSN)' NEW_DB_DSN='$(NEW_DB_DSN)' MIGRATION_PHASE=refs DRY_RUN=0 $(MIGRATE_CMD)

# ---- Core migration (persons + titles) ----
# (You already ran persons successfully; see titles-only targets below.)

.PHONY: migrate-core-dry-run
migrate-core-dry-run: ## DRY-RUN core migration (persons + titles)
	@echo ">> DRY-RUN core migration (persons + titles)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	@OLD_DB_DSN='$(OLD_DB_DSN)' NEW_DB_DSN='$(NEW_DB_DSN)' MIGRATION_PHASE=core DRY_RUN=1 $(MIGRATE_CMD)

.PHONY: migrate-core
migrate-core: ## REAL core migration (persons + titles)
	@echo ">> REAL core migration (persons + titles)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	@OLD_DB_DSN='$(OLD_DB_DSN)' NEW_DB_DSN='$(NEW_DB_DSN)' MIGRATION_PHASE=core DRY_RUN=0 $(MIGRATE_CMD)

# ---- Titles-only migration (what you want to run now) ----

.PHONY: migrate-titles-dry-run
migrate-titles-dry-run: ## DRY-RUN titles migration only
	@echo ">> DRY-RUN core TITLE migration only"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	@OLD_DB_DSN='$(OLD_DB_DSN)' NEW_DB_DSN='$(NEW_DB_DSN)' MIGRATION_PHASE=core-title DRY_RUN=1 $(MIGRATE_CMD)

.PHONY: migrate-titles
migrate-titles: ## REAL titles migration only
	@echo ">> REAL core TITLE migration only"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	@OLD_DB_DSN='$(OLD_DB_DSN)' NEW_DB_DSN='$(NEW_DB_DSN)' MIGRATION_PHASE=core-title DRY_RUN=0 $(MIGRATE_CMD)

# ---- Junctions (later: countries, languages, genres, cast, media_file, etc.) ----

.PHONY: migrate-junctions-dry-run
migrate-junctions-dry-run: ## DRY-RUN basic junction migrations (no writes)
	@echo ">> DRY-RUN junctions (basic title_* junctions) [stub]"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	@$(GO) run ./cmd/migrate-old-db \
		-old "$(OLD_DB_DSN)" \
		-new "$(NEW_DB_DSN)" \
		-phase junctions \
		-dry-run

.PHONY: migrate-junctions
migrate-junctions: ## REAL basic junction migrations (title_* junctions)
	@echo ">> REAL junctions migration (basic title_* junctions)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	@$(GO) run ./cmd/migrate-old-db \
		-old "$(OLD_DB_DSN)" \
		-new "$(NEW_DB_DSN)" \
		-phase junctions

# ---------------------------
# Junctions: split phases
# ---------------------------

migrate-junctions-country-dry-run:
	@echo ">> DRY-RUN junctions COUNTRY (CountryTitleLine -> title_country)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-country \
	  -dry-run

migrate-junctions-country:
	@echo ">> REAL junctions COUNTRY (CountryTitleLine -> title_country)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-country

migrate-junctions-language-dry-run:
	@echo ">> DRY-RUN junctions LANGUAGE (LanguageTitleLine -> title_language)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-language \
	  -dry-run

migrate-junctions-language:
	@echo ">> REAL junctions LANGUAGE (LanguageTitleLine -> title_language)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-language

migrate-junctions-genre-dry-run:
	@echo ">> DRY-RUN junctions GENRE (GenreTitleLine -> title_genre)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-genre \
	  -dry-run

migrate-junctions-genre:
	@echo ">> REAL junctions GENRE (GenreTitleLine -> title_genre)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-genre

migrate-junctions-alias-dry-run:
	@echo ">> DRY-RUN junctions ALIAS (KnownAsTitleLine -> title_alias)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-alias \
	  -dry-run

migrate-junctions-alias:
	@echo ">> REAL junctions ALIAS (KnownAsTitleLine -> title_alias)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-alias

migrate-junctions-certificate-dry-run:
	@echo ">> DRY-RUN junctions CERTIFICATE (CertificateTitleLine -> title_certificate)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-certificate \
	  -dry-run

migrate-junctions-certificate:
	@echo ">> REAL junctions CERTIFICATE (CertificateTitleLine -> title_certificate)"
	@echo "   OLD_DB: $(OLD_DB_DSN)"
	@echo "   NEW_DB: $(NEW_DB_DSN)"
	go run ./cmd/migrate-old-db \
	  -old "$(OLD_DB_DSN)" \
	  -new "$(NEW_DB_DSN)" \
	  -phase junctions-certificate
