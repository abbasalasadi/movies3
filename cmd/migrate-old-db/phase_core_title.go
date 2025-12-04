// cmd/migrate-old-db/phase_titles.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// MigrateCoreTitlesPhase runs the "core-title" phase: TitleTable → title.
func MigrateCoreTitlesPhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("=== Starting migration phase=%q dryRun=%v ===", "core-title", dryRun)

	start := time.Now()
	if err := migrateTitles(ctx, oldDB, newDB, dryRun); err != nil {
		return fmt.Errorf("migrateTitles: %w", err)
	}

	if !dryRun {
		if err := backfillTitleParents(ctx, oldDB, newDB); err != nil {
			return fmt.Errorf("backfillTitleParents: %w", err)
		}
	}

	log.Printf("=== Migration phase=%q completed successfully in %s ===", "core-title", time.Since(start))
	return nil
}

func migrateTitles(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Println("--- Migrating title (TitleTable → title) ---")

	// 1) Count rows in old TitleTable
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "Tables"."TitleTable"`).Scan(&total); err != nil {
		return fmt.Errorf("count TitleTable: %w", err)
	}
	log.Printf("migrateTitles: found %d rows in Tables.\"TitleTable\"", total)

	if dryRun {
		log.Printf("migrateTitles [DRY-RUN]: would read %d TitleTable rows and insert/update into title", total)
		return nil
	}

	// 2) Prepare INSERT for new title table
	//
	// Columns we populate:
	//   id,
	//   title_type_id,
	//   primary_title,
	//   original_title,
	//   start_year,
	//   runtime_minutes,
	//   primary_country_id,
	//   poster_url,
	//   metacritic_rating,
	//   revenue,
	//   imdb_rating,
	//   imdb_votes,
	//   popularity,
	//   parent_title_id   <-- initially NULL; backfilled in a second pass
	//   season_number,
	//   episode_number,
	//   total_seasons,
	//   total_episodes,
	//   date_released,
	//   date_added,
	//   date_updated,     <-- NEVER NULL (fallback to date_added)
	//   is_available,
	//   viewed_count,
	//   played_count,
	//   liked_count,
	//   disliked_count,
	//   folder_name,
	//   folder_path
	//
	// 28 columns → 28 VALUES placeholders.
	const insertSQL = `
INSERT INTO title (
	id,
	title_type_id,
	primary_title,
	original_title,
	start_year,
	runtime_minutes,
	primary_country_id,
	poster_url,
	metacritic_rating,
	revenue,
	imdb_rating,
	imdb_votes,
	popularity,
	parent_title_id,
	season_number,
	episode_number,
	total_seasons,
	total_episodes,
	date_released,
	date_added,
	date_updated,
	is_available,
	viewed_count,
	played_count,
	liked_count,
	disliked_count,
	folder_name,
	folder_path
) VALUES (
	$1,  $2,  $3,  $4,  $5,  $6,  $7,
	$8,  $9,  $10, $11, $12, $13, $14,
	$15, $16, $17, $18, $19, $20, $21,
	$22, $23, $24, $25, $26, $27, $28
)
ON CONFLICT (id) DO UPDATE SET
	title_type_id      = EXCLUDED.title_type_id,
	primary_title      = EXCLUDED.primary_title,
	original_title     = EXCLUDED.original_title,
	start_year         = EXCLUDED.start_year,
	runtime_minutes    = EXCLUDED.runtime_minutes,
	primary_country_id = EXCLUDED.primary_country_id,
	poster_url         = EXCLUDED.poster_url,
	metacritic_rating  = EXCLUDED.metacritic_rating,
	revenue            = EXCLUDED.revenue,
	imdb_rating        = EXCLUDED.imdb_rating,
	imdb_votes         = EXCLUDED.imdb_votes,
	popularity         = EXCLUDED.popularity,
	-- parent_title_id is backfilled later
	season_number      = EXCLUDED.season_number,
	episode_number     = EXCLUDED.episode_number,
	total_seasons      = EXCLUDED.total_seasons,
	total_episodes     = EXCLUDED.total_episodes,
	date_released      = EXCLUDED.date_released,
	date_added         = EXCLUDED.date_added,
	date_updated       = EXCLUDED.date_updated,
	is_available       = EXCLUDED.is_available,
	viewed_count       = EXCLUDED.viewed_count,
	played_count       = EXCLUDED.played_count,
	liked_count        = EXCLUDED.liked_count,
	disliked_count     = EXCLUDED.disliked_count,
	folder_name        = EXCLUDED.folder_name,
	folder_path        = EXCLUDED.folder_path;
`

	stmt, err := newDB.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert title: %w", err)
	}
	defer stmt.Close()

	// 3) Stream rows from old TitleTable
	const selectSQL = `
SELECT
	"TitleID",
	"TitleType",
	"TitleName",
	"OriginalTitle",
	"TitleYear",
	"TitleLength",
	"TitleCountry",
	"PosterURL",
	"MetacriticRating",
	"Revenue",
	"IMDbRating",
	"IMDbVotes",
	"Popularity",
	"ParentID",
	"EpisodeSeason",
	"EpisodeNumber",
	"TotalSeasons",
	"TotalEpisodes",
	"DateReleased",
	"DateAdded",
	"DateUpdated",
	"Available",
	"Viewed",
	"Played",
	"Liked",
	"UnLiked",
	"FolderName",
	"FolderPath"
FROM "Tables"."TitleTable"
ORDER BY "TitleID"
`

	rows, err := oldDB.QueryContext(ctx, selectSQL)
	if err != nil {
		return fmt.Errorf("query TitleTable: %w", err)
	}
	defer rows.Close()

	start := time.Now()
	var (
		processed int64
		inserted  int64
	)

	for rows.Next() {
		var (
			titleID       int64
			titleType     sql.NullInt64
			titleName     string
			originalTitle sql.NullString
			titleYear     sql.NullInt64
			titleLength   sql.NullInt64
			titleCountry  sql.NullInt64
			posterURL     sql.NullString
			metacritic    sql.NullInt64
			revenue       sql.NullInt64
			imdbRating    sql.NullFloat64
			imdbVotes     sql.NullInt64
			popularity    sql.NullInt64
			parentID      sql.NullInt64
			episodeSeason sql.NullString
			episodeNumber sql.NullInt64
			totalSeasons  sql.NullInt64
			totalEpisodes sql.NullInt64
			dateReleased  sql.NullTime
			dateAdded     time.Time
			dateUpdated   sql.NullTime
			available     sql.NullBool
			viewed        sql.NullInt64
			played        sql.NullInt64
			liked         sql.NullInt64
			unliked       sql.NullInt64
			folderName    sql.NullString
			folderPath    sql.NullString
		)

		if err := rows.Scan(
			&titleID,
			&titleType,
			&titleName,
			&originalTitle,
			&titleYear,
			&titleLength,
			&titleCountry,
			&posterURL,
			&metacritic,
			&revenue,
			&imdbRating,
			&imdbVotes,
			&popularity,
			&parentID,
			&episodeSeason,
			&episodeNumber,
			&totalSeasons,
			&totalEpisodes,
			&dateReleased,
			&dateAdded,
			&dateUpdated,
			&available,
			&viewed,
			&played,
			&liked,
			&unliked,
			&folderName,
			&folderPath,
		); err != nil {
			return fmt.Errorf("scan TitleTable row: %w", err)
		}

		processed++

		// Convert nullable values to appropriate Go / SQL types
		titleTypeID := nullInt64OrNil(titleType)
		startYear := nullInt64OrNil(titleYear)
		runtimeMinutes := nullInt64OrNil(titleLength)
		primaryCountryID := nullInt64OrNil(titleCountry)
		posterURLVal := nullStringOrNil(posterURL)
		metacriticVal := nullInt64OrNil(metacritic)
		revenueVal := nullInt64OrNil(revenue)
		imdbRatingVal := nullFloat64OrNil(imdbRating)
		imdbVotesVal := nullInt64OrNil(imdbVotes)
		popularityVal := nullInt64OrNil(popularity)
		// parentID is handled in a second pass, so we don't set parent_title_id here
		seasonNumber := parseSeasonToInt64(episodeSeason)
		episodeNumberVal := nullInt64OrNil(episodeNumber)
		totalSeasonsVal := nullInt64OrNil(totalSeasons)
		totalEpisodesVal := nullInt64OrNil(totalEpisodes)
		dateReleasedVal := nullTimeOrNil(dateReleased)

		// date_updated is NOT NULL in the new schema.
		// If DateUpdated is NULL, we fall back to DateAdded.
		dateUpdatedVal := dateAdded
		if dateUpdated.Valid {
			dateUpdatedVal = dateUpdated.Time
		}

		isAvailable := boolOrFalse(available)

		// Ensure counts are never NULL to satisfy NOT NULL
		viewedCount := int64(0)
		if viewed.Valid {
			viewedCount = viewed.Int64
		}
		playedCount := int64(0)
		if played.Valid {
			playedCount = played.Int64
		}
		likedCount := int64(0)
		if liked.Valid {
			likedCount = liked.Int64
		}
		dislikedCount := int64(0)
		if unliked.Valid {
			dislikedCount = unliked.Int64
		}

		folderNameVal := nullStringOrNil(folderName)
		folderPathVal := nullStringOrNil(folderPath)

		if _, err := stmt.ExecContext(
			ctx,
			titleID,              // id
			titleTypeID,          // title_type_id
			titleName,            // primary_title
			originalTitle.String, // original_title ("" if NULL)
			startYear,            // start_year
			runtimeMinutes,       // runtime_minutes
			primaryCountryID,     // primary_country_id
			posterURLVal,         // poster_url
			metacriticVal,        // metacritic_rating
			revenueVal,           // revenue
			imdbRatingVal,        // imdb_rating
			imdbVotesVal,         // imdb_votes
			popularityVal,        // popularity
			nil,                  // parent_title_id (backfilled later)
			seasonNumber,         // season_number
			episodeNumberVal,     // episode_number
			totalSeasonsVal,      // total_seasons
			totalEpisodesVal,     // total_episodes
			dateReleasedVal,      // date_released
			dateAdded,            // date_added
			dateUpdatedVal,       // date_updated (never NULL)
			isAvailable,          // is_available
			viewedCount,          // viewed_count
			playedCount,          // played_count
			likedCount,           // liked_count
			dislikedCount,        // disliked_count
			folderNameVal,        // folder_name
			folderPathVal,        // folder_path
		); err != nil {
			return fmt.Errorf("insert title id=%d: %w", titleID, err)
		}

		inserted++

		if inserted%500000 == 0 {
			percent := float64(inserted) * 100.0 / float64(total)
			log.Printf("migrateTitles: inserted/updated %d/%d titles (%.1f%%)", inserted, total, percent)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate TitleTable rows: %w", err)
	}

	percent := float64(inserted)
	if total > 0 {
		percent = percent * 100.0 / float64(total)
	}
	log.Printf("migrateTitles: inserted/updated %d/%d titles (%.1f%%)", inserted, total, percent)
	log.Printf("--- Done title: %d rows processed in %s ---", processed, time.Since(start))

	return nil
}

// backfillTitleParents runs AFTER all titles are inserted.
// It reads (TitleID, ParentID) from old TitleTable and updates title.parent_title_id.
func backfillTitleParents(ctx context.Context, oldDB, newDB *sql.DB) error {
	log.Println("--- Backfilling title.parent_title_id from TitleTable.ParentID ---")

	// Count rows with a positive ParentID
	var total int64
	if err := oldDB.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM "Tables"."TitleTable"
		WHERE "ParentID" IS NOT NULL
		  AND "ParentID" > 0
	`).Scan(&total); err != nil {
		return fmt.Errorf("count ParentID>0: %w", err)
	}
	log.Printf("backfillTitleParents: %d rows with ParentID>0", total)

	if total == 0 {
		log.Println("backfillTitleParents: no parent relationships to backfill")
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "TitleID", "ParentID"
		FROM "Tables"."TitleTable"
		WHERE "ParentID" IS NOT NULL
		  AND "ParentID" > 0
		ORDER BY "TitleID"
	`)
	if err != nil {
		return fmt.Errorf("query ParentID rows: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		UPDATE title
		SET parent_title_id = $2
		WHERE id = $1
	`)
	if err != nil {
		return fmt.Errorf("prepare UPDATE title.parent_title_id: %w", err)
	}
	defer stmt.Close()

	start := time.Now()
	var processed int64

	for rows.Next() {
		var (
			titleID  int64
			parentID int64
		)
		if err := rows.Scan(&titleID, &parentID); err != nil {
			return fmt.Errorf("scan ParentID row: %w", err)
		}

		// ParentID should be >0 due to WHERE clause, but double-check
		if parentID <= 0 {
			continue
		}

		if _, err := stmt.ExecContext(ctx, titleID, parentID); err != nil {
			return fmt.Errorf("update parent_title_id for title id=%d: %w", titleID, err)
		}

		processed++
		if processed%50000 == 0 {
			pct := float64(processed) * 100.0 / float64(total)
			log.Printf("migrateTitles: inserted/updated %d/%d titles (%.1f%%)", processed, total, pct)
		}

	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate ParentID rows: %w", err)
	}

	percent := float64(processed)
	if total > 0 {
		percent = percent * 100.0 / float64(total)
	}
	log.Printf("backfillTitleParents: updated %d/%d titles (%.1f%%) in %s", processed, total, percent, time.Since(start))

	return nil
}

// Helpers for nullable conversions

func nullInt64OrNil(n sql.NullInt64) interface{} {
	if !n.Valid {
		return nil
	}
	return n.Int64
}

func nullFloat64OrNil(n sql.NullFloat64) interface{} {
	if !n.Valid {
		return nil
	}
	return n.Float64
}

func nullTimeOrNil(n sql.NullTime) interface{} {
	if !n.Valid {
		return nil
	}
	return n.Time
}

func nullStringOrNil(n sql.NullString) interface{} {
	if !n.Valid {
		return nil
	}
	s := strings.TrimSpace(n.String)
	if s == "" {
		return nil
	}
	return s
}

// parseSeasonToInt64 tries to extract an integer season from strings like "1", "S1", "Season 1".
func parseSeasonToInt64(n sql.NullString) interface{} {
	if !n.Valid {
		return nil
	}
	s := strings.TrimSpace(n.String)
	if s == "" {
		return nil
	}

	var digits strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	if digits.Len() == 0 {
		return nil
	}

	v, err := strconv.Atoi(digits.String())
	if err != nil {
		return nil
	}
	return int64(v)
}

func boolOrFalse(n sql.NullBool) bool {
	if !n.Valid {
		return false
	}
	return n.Bool
}
