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

var (
	oldDSN = flag.String("old", "", "Postgres DSN for OLD database (mediadb)")
	newDSN = flag.String("new", "", "Postgres DSN for NEW database (movies3db)")
	phase  = flag.String("phase", "refs", "Migration phase (refs | core-persons | core-title | junctions-country | junctions-language | junctions-genre | junctions-alias | junctions-certificate | junctions)")
	dryRun = flag.Bool("dry-run", false, "if set, do NOT write to new DB; just read and count")
)

func main() {
	log.SetOutput(os.Stdout)
	flag.Parse()

	if *oldDSN == "" || *newDSN == "" {
		log.Printf("ERROR: both -old and -new DSNs are required")
		flag.Usage()
		os.Exit(2)
	}

	ctx := context.Background()

	log.Printf("Connecting to OLD DB: %s", *oldDSN)
	oldDB, err := sql.Open("postgres", *oldDSN)
	if err != nil {
		log.Fatalf("open old DB: %v", err)
	}
	defer oldDB.Close()

	log.Printf("Connecting to NEW DB: %s", *newDSN)
	newDB, err := sql.Open("postgres", *newDSN)
	if err != nil {
		log.Fatalf("open new DB: %v", err)
	}
	defer newDB.Close()

	if err := oldDB.PingContext(ctx); err != nil {
		log.Fatalf("ping old DB: %v", err)
	}
	if err := newDB.PingContext(ctx); err != nil {
		log.Fatalf("ping new DB: %v", err)
	}

	log.Printf("=== Starting migration phase=%q dryRun=%v ===", *phase, *dryRun)
	start := time.Now()

	var phaseErr error

	switch *phase {
	case "refs":
		phaseErr = MigrateRefsPhase(ctx, oldDB, newDB, *dryRun)

	case "core-persons":
		// implemented earlier, not touched here
		phaseErr = MigrateCorePersonsPhase(ctx, oldDB, newDB, *dryRun)

	case "core-title":
		phaseErr = MigrateCoreTitlesPhase(ctx, oldDB, newDB, *dryRun)

	case "junctions":
		// OPTIONAL umbrella phase if you still use it:
		// basic + others, depending on how you wired it
		phaseErr = MigrateJunctionsPhase(ctx, oldDB, newDB, *dryRun)

	case "junctions-country":
		phaseErr = MigrateJunctionsCountryPhase(ctx, oldDB, newDB, *dryRun)

	case "junctions-language":
		phaseErr = MigrateJunctionsLanguagePhase(ctx, oldDB, newDB, *dryRun)

	case "junctions-genre":
		phaseErr = MigrateJunctionsGenrePhase(ctx, oldDB, newDB, *dryRun)

	case "junctions-alias":
		phaseErr = MigrateJunctionsAliasPhase(ctx, oldDB, newDB, *dryRun)

	case "junctions-certificate":
		phaseErr = MigrateJunctionsCertificatePhase(ctx, oldDB, newDB, *dryRun)

	default:
		log.Fatalf("unknown phase %q", *phase)
	}

	if phaseErr != nil {
		log.Fatalf("migration phase %q failed: %v", *phase, phaseErr)
	}

	log.Printf("=== Migration phase=%q completed successfully in %s ===",
		*phase, time.Since(start).Truncate(time.Millisecond))
}
