package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// MigrateJunctionsAliasPhase migrates KnownAsTitleLine -> title_alias.
func MigrateJunctionsAliasPhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	mode := "LIVE"
	if dryRun {
		mode = "DRY-RUN"
	}
	log.Printf("=== Starting migration phase=\"junctions-alias\" dryRun=%v [%s] ===", dryRun, mode)

	if err := migrateTitleAlias(ctx, oldDB, newDB, dryRun); err != nil {
		return fmt.Errorf("migrateTitleAlias: %w", err)
	}

	log.Printf("=== Migration phase=\"junctions-alias\" completed successfully ===")
	return nil
}

// loadImdbToNewTitleIDMap builds a map imdb_id -> new title.id from the NEW DB.
func loadImdbToNewTitleIDMap(ctx context.Context, newDB *sql.DB) (map[string]int64, error) {
	log.Printf("--- Building IMDbID -> new title.id map from movies3db.title ---")

	rows, err := newDB.QueryContext(ctx, `SELECT id, imdb_id FROM title WHERE imdb_id IS NOT NULL`)
	if err != nil {
		return nil, fmt.Errorf("query new title table: %w", err)
	}
	defer rows.Close()

	m := make(map[string]int64, 8_000_000) // big but OK for your dataset

	var (
		id     int64
		imdbID string
		count  int64
	)
	for rows.Next() {
		if err := rows.Scan(&id, &imdbID); err != nil {
			return nil, fmt.Errorf("scan new title row: %w", err)
		}
		m[imdbID] = id
		count++
		if count%500000 == 0 {
			log.Printf("loadImdbToNewTitleIDMap: loaded %d titles into map", count)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate new title rows: %w", err)
	}

	log.Printf("loadImdbToNewTitleIDMap: loaded %d titles into IMDbID map", count)
	return m, nil
}

// migrateTitleAlias copies KnownAsTitleLine into title_alias using IMDbID mapping.
func migrateTitleAlias(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	// Count source rows first.
	var total int64
	if err := oldDB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM "Lines"."KnownAsTitleLine"`).Scan(&total); err != nil {
		return fmt.Errorf("count KnownAsTitleLine: %w", err)
	}
	log.Printf("migrateTitleAlias: %d rows in Lines.\"KnownAsTitleLine\"", total)

	if dryRun {
		log.Printf("migrateTitleAlias [DRY-RUN]: would process %d rows", total)
		return nil
	}

	// Build IMDbID -> new title.id map from the NEW DB.
	imdbMap, err := loadImdbToNewTitleIDMap(ctx, newDB)
	if err != nil {
		return err
	}

	// Start a transaction for bulk inserts.
	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx for title_alias: %w", err)
	}
	defer func() {
		// If commit fails, tx.Rollback() will be called explicitly below.
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO title_alias (title_id, alias)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert title_alias: %w", err)
	}
	defer stmt.Close()

	// Query old KnownAsTitleLine joined to TitleTable so we can get IMDbID.
	rows, err := oldDB.QueryContext(ctx, `
		SELECT t."IMDbID", k."KnownAsTitle"
		FROM "Lines"."KnownAsTitleLine" k
		JOIN "Tables"."TitleTable" t ON t."TitleID" = k."TitleID"
		WHERE t."IMDbID" IS NOT NULL
	`)
	if err != nil {
		return fmt.Errorf("query KnownAsTitleLine join TitleTable: %w", err)
	}
	defer rows.Close()

	var (
		processed int64
		inserted  int64
		skipped   int64
	)

	const progressStep int64 = 50000

	for rows.Next() {
		var imdbID, alias string
		if err := rows.Scan(&imdbID, &alias); err != nil {
			return fmt.Errorf("scan KnownAsTitleLine row: %w", err)
		}

		newTitleID, ok := imdbMap[imdbID]
		if !ok {
			// We don't have this title in the new DB (e.g. not migrated or filtered out)
			skipped++
		} else {
			if _, err := stmt.ExecContext(ctx, newTitleID, alias); err != nil {
				return fmt.Errorf("insert title_alias (imdb_id=%s): %w", imdbID, err)
			}
			inserted++
		}

		processed++
		if processed%progressStep == 0 {
			percent := float64(processed) / float64(total) * 100.0
			log.Printf("migrateTitleAlias: processed %d/%d rows (%.1f%%), inserted=%d, skipped=%d",
				processed, total, percent, inserted, skipped)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate KnownAsTitleLine rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit title_alias tx: %w", err)
	}

	log.Printf("migrateTitleAlias: processed %d/%d rows, inserted=%d, skipped (no title mapping)=%d",
		processed, total, inserted, skipped)

	return nil
}
