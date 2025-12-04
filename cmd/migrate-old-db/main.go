// cmd/migrate-old-db/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	var oldDSNFlag, newDSNFlag, phaseFlag string
	var dryRunFlag bool

	flag.StringVar(&oldDSNFlag, "old", "", "Postgres DSN for OLD database (mediadb)")
	flag.StringVar(&newDSNFlag, "new", "", "Postgres DSN for NEW database (movies3db)")
	flag.BoolVar(&dryRunFlag, "dry-run", false, "if set, do NOT write to new DB; just read and count")
	flag.StringVar(&phaseFlag, "phase", "", "migration phase: refs | core | core-person | core-title | junctions")
	flag.Parse()

	oldDSN := firstNonEmpty(os.Getenv("OLD_DB_DSN"), oldDSNFlag)
	newDSN := firstNonEmpty(os.Getenv("NEW_DB_DSN"), newDSNFlag)
	phase := firstNonEmpty(os.Getenv("MIGRATION_PHASE"), phaseFlag)
	if phase == "" {
		phase = "refs"
	}

	dryRun := dryRunFlag
	if env := os.Getenv("DRY_RUN"); env != "" {
		switch env {
		case "1", "true", "TRUE", "True":
			dryRun = true
		case "0", "false", "FALSE", "False":
			dryRun = false
		}
	}

	if oldDSN == "" || newDSN == "" {
		log.Fatalf("both OLD_DB_DSN/--old and NEW_DB_DSN/--new must be set")
	}

	log.Printf("Connecting to OLD DB: %s", oldDSN)
	oldDB, err := sql.Open("postgres", oldDSN)
	if err != nil {
		log.Fatalf("open old DB: %v", err)
	}
	defer oldDB.Close()

	log.Printf("Connecting to NEW DB: %s", newDSN)
	newDB, err := sql.Open("postgres", newDSN)
	if err != nil {
		log.Fatalf("open new DB: %v", err)
	}
	defer newDB.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 24*time.Hour)
	defer cancel()

	if err := oldDB.PingContext(ctx); err != nil {
		log.Fatalf("ping old DB: %v", err)
	}
	if err := newDB.PingContext(ctx); err != nil {
		log.Fatalf("ping new DB: %v", err)
	}

	log.Printf("=== Starting migration phase=%q dryRun=%v ===", phase, dryRun)

	var phaseErr error

	switch phase {
	case "refs":
		phaseErr = MigrateRefsPhase(ctx, oldDB, newDB, dryRun)

	case "core":
		// Backwards-compatible: run persons then titles
		log.Printf(`phase "core" selected: running persons THEN titles`)
		if err := MigrateCorePersonsPhase(ctx, oldDB, newDB, dryRun); err != nil {
			phaseErr = err
		} else {
			phaseErr = MigrateCoreTitlesPhase(ctx, oldDB, newDB, dryRun)
		}

	case "core-person":
		phaseErr = MigrateCorePersonsPhase(ctx, oldDB, newDB, dryRun)

	case "core-title":
		phaseErr = MigrateCoreTitlesPhase(ctx, oldDB, newDB, dryRun)

	case "junctions":
		phaseErr = MigrateJunctionsPhase(ctx, oldDB, newDB, dryRun)

	default:
		log.Fatalf("unknown phase %q (expected refs|core|core-person|core-title|junctions)", phase)
	}

	if phaseErr != nil {
		log.Fatalf("migration phase %q failed: %v", phase, phaseErr)
	}

	log.Printf("=== Migration phase=%q completed successfully ===", phase)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
