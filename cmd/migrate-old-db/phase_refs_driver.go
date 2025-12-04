// cmd/migrate-old-db/phase_refs_driver.go
package main

import (
	"context"
	"database/sql"
	"log"
)

// MigrateRefsPhase is called by main.go for phase="refs".
// In your current workflow, the reference data has already been migrated
// successfully, so this implementation is a safe no-op.
// If you re-run it, it will simply log and return.
func MigrateRefsPhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("MigrateRefsPhase: reference data migration already completed earlier; no-op for now (dryRun=%v)", dryRun)
	return nil
}

// MigrateJunctionsPhase will later handle all link/junction tables
// (title_country, title_language, title_genre, title_cast, media_file, etc.).
// For now it's a placeholder so the program compiles and phases 'core-*' can run.
func MigrateJunctionsPhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("MigrateJunctionsPhase: junctions migration not implemented yet; TODO (dryRun=%v)", dryRun)
	// Later we will add the real implementation here.
	return nil
}
