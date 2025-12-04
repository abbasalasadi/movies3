// cmd/migrate-old-db/phase_core_person.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// MigrateCorePersonsPhase runs ONLY the person migration:
//   Tables."CastTable" -> public.person
func MigrateCorePersonsPhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	start := time.Now()
	log.Printf("=== Starting migration phase=\"core-person\" dryRun=%v ===", dryRun)

	if err := migratePersons(ctx, oldDB, newDB, dryRun); err != nil {
		return fmt.Errorf("migratePersons: %w", err)
	}

	log.Printf("=== Migration phase=\"core-person\" completed successfully in %s ===", time.Since(start))
	return nil
}

// ======================
//   PERSON MIGRATION
// ======================
//
// OLD: Tables."CastTable"
//   "CastID"        bigint NOT NULL
//   "CastName"      varchar NOT NULL
//   "IsDirector"    boolean
//   "IsWriter"      boolean
//   "IsCharacter"   boolean
//
// NEW: public.person
//   id                  BIGINT PK
//   imdb_id             TEXT (unused for now)
//   name                TEXT NOT NULL
//   primary_profession  TEXT
//   created_at          TIMESTAMPTZ NOT NULL
//   updated_at          TIMESTAMPTZ NOT NULL
//
// We keep IDs identical so junction tables can refer to them.

func migratePersons(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Println("--- Migrating person (CastTable → person) ---")

	// 1) Count rows in old CastTable
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "Tables"."CastTable"`).Scan(&total); err != nil {
		return fmt.Errorf("count CastTable: %w", err)
	}
	log.Printf("migratePersons: found %d rows in Tables.\"CastTable\"", total)

	if dryRun {
		log.Printf("migratePersons [DRY-RUN]: would read %d CastTable rows and insert/update into person", total)
		return nil
	}

	// 2) Prepare insert/upsert statement in new DB
	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO person (
			id,
			name,
			primary_profession,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, now(), now()
		)
		ON CONFLICT (id) DO UPDATE
		SET
			name               = EXCLUDED.name,
			primary_profession = EXCLUDED.primary_profession,
			updated_at         = now()
	`)
	if err != nil {
		return fmt.Errorf("prepare insert person: %w", err)
	}
	defer stmt.Close()

	// 3) Stream rows from old DB
	rows, err := oldDB.QueryContext(ctx, `
		SELECT
			"CastID",
			"CastName",
			COALESCE("IsDirector", false)  AS is_director,
			COALESCE("IsWriter", false)    AS is_writer,
			COALESCE("IsCharacter", false) AS is_character
		FROM "Tables"."CastTable"
		ORDER BY "CastID"
	`)
	if err != nil {
		return fmt.Errorf("select CastTable: %w", err)
	}
	defer rows.Close()

	var processed int64
	lastLog := time.Now()

	for rows.Next() {
		var (
			id          int64
			name        string
			isDirector  bool
			isWriter    bool
			isCharacter bool
		)

		if err := rows.Scan(&id, &name, &isDirector, &isWriter, &isCharacter); err != nil {
			return fmt.Errorf("scan CastTable row: %w", err)
		}

		// Build primary_profession from flags
		var profs []string
		if isDirector {
			profs = append(profs, "director")
		}
		if isWriter {
			profs = append(profs, "writer")
		}
		if isCharacter {
			profs = append(profs, "actor")
		}
		primaryProfession := strings.Join(profs, ",")

		if _, err := stmt.ExecContext(ctx, id, name, primaryProfession); err != nil {
			return fmt.Errorf("insert person id=%d: %w", id, err)
		}

		processed++
		if processed%50000 == 0 || time.Since(lastLog) > 10*time.Second {
			pct := float64(processed) * 100.0 / float64(total)
			log.Printf("migratePersons: inserted/updated %d/%d persons (%.1f%%)", processed, total, pct)
			lastLog = time.Now()
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CastTable: %w", err)
	}

	log.Printf("--- Done person: %d rows processed ---", processed)
	return nil
}
