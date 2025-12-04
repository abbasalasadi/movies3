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
