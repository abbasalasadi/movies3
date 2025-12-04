// cmd/migrate-old-db/phase_junctions.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

const junctionProgressEvery = 50000

// MigrateJunctionsPhase keeps the old "all basic junctions" behaviour,
// but in practice you'll usually call the split phases instead.
func MigrateJunctionsPhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
    log.Printf("=== Starting migration phase=\"junctions\" dryRun=%v ===", dryRun)

    if err := MigrateJunctionsCountryPhase(ctx, oldDB, newDB, dryRun); err != nil {
        return fmt.Errorf("junctions-country: %w", err)
    }
    if err := MigrateJunctionsLanguagePhase(ctx, oldDB, newDB, dryRun); err != nil {
        return fmt.Errorf("junctions-language: %w", err)
    }
    if err := MigrateJunctionsGenrePhase(ctx, oldDB, newDB, dryRun); err != nil {
        return fmt.Errorf("junctions-genre: %w", err)
    }
    if err := MigrateJunctionsAliasPhase(ctx, oldDB, newDB, dryRun); err != nil {
        return fmt.Errorf("junctions-alias: %w", err)
    }
    if err := MigrateJunctionsCertificatePhase(ctx, oldDB, newDB, dryRun); err != nil {
        return fmt.Errorf("junctions-certificate: %w", err)
    }

    log.Printf("=== Migration phase=\"junctions\" completed successfully ===")
    return nil
}


//
// Split phases
//

// Country: Lines."CountryTitleLine" -> title_country
func MigrateJunctionsCountryPhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("=== Starting migration phase=\"junctions-country\" dryRun=%v ===", dryRun)

	log.Printf("--- Building country ID map (old References.\"CountryRef\" -> new country_ref) ---")
	countryIDMap, err := buildCountryIDMap(ctx, oldDB, newDB)
	if err != nil {
		return fmt.Errorf("buildCountryIDMap: %w", err)
	}

	if err := migrateTitleCountry(ctx, oldDB, newDB, countryIDMap, dryRun); err != nil {
		return fmt.Errorf("migrateTitleCountry: %w", err)
	}

	log.Printf("=== Migration phase=\"junctions-country\" completed successfully ===")
	return nil
}

// Language: Lines."LanguageTitleLine" -> title_language
func MigrateJunctionsLanguagePhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("=== Starting migration phase=\"junctions-language\" dryRun=%v ===", dryRun)

	log.Printf("--- Building language ID map (old References.\"LanguageRef\" -> new language_ref) ---")
	langIDMap, err := buildLanguageIDMap(ctx, oldDB, newDB)
	if err != nil {
		return fmt.Errorf("buildLanguageIDMap: %w", err)
	}

	if err := migrateTitleLanguage(ctx, oldDB, newDB, langIDMap, dryRun); err != nil {
		return fmt.Errorf("migrateTitleLanguage: %w", err)
	}

	log.Printf("=== Migration phase=\"junctions-language\" completed successfully ===")
	return nil
}

// Genre: Lines."GenreTitleLine" -> title_genre
func MigrateJunctionsGenrePhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("=== Starting migration phase=\"junctions-genre\" dryRun=%v ===", dryRun)

	log.Printf("--- Building genre ID map (old References.\"GenreRef\" -> new genre_ref) ---")
	genreIDMap, err := buildGenreIDMap(ctx, oldDB, newDB)
	if err != nil {
		return fmt.Errorf("buildGenreIDMap: %w", err)
	}

	if err := migrateTitleGenre(ctx, oldDB, newDB, genreIDMap, dryRun); err != nil {
		return fmt.Errorf("migrateTitleGenre: %w", err)
	}

	log.Printf("=== Migration phase=\"junctions-genre\" completed successfully ===")
	return nil
}

// Certificate: Lines."CertificateTitleLine" -> title_certificate
func MigrateJunctionsCertificatePhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("=== Starting migration phase=\"junctions-certificate\" dryRun=%v ===", dryRun)

	log.Printf("--- Building country ID map (for title_certificate) ---")
	countryIDMap, err := buildCountryIDMap(ctx, oldDB, newDB)
	if err != nil {
		return fmt.Errorf("buildCountryIDMap (for certificate): %w", err)
	}

	log.Printf("--- Building certificate ID map (old References.\"CertificateRef\" -> new certificate_ref) ---")
	certIDMap, err := buildCertificateIDMap(ctx, oldDB, newDB)
	if err != nil {
		return fmt.Errorf("buildCertificateIDMap: %w", err)
	}

	if err := migrateTitleCertificate(ctx, oldDB, newDB, countryIDMap, certIDMap, dryRun); err != nil {
		return fmt.Errorf("migrateTitleCertificate: %w", err)
	}

	log.Printf("=== Migration phase=\"junctions-certificate\" completed successfully ===")
	return nil
}

//
// ID map builders
//

// old: "References"."CountryRef"(CountryID, CountryName, CountryCode)
// new: country_ref(id, name, iso2_code, iso3_code)
func buildCountryIDMap(ctx context.Context, oldDB, newDB *sql.DB) (map[int32]int16, error) {
	type newCountry struct {
		id   int16
		name string
	}
	newByName := make(map[string]int16)

	rows, err := newDB.QueryContext(ctx, `SELECT id, name FROM country_ref`)
	if err != nil {
		return nil, fmt.Errorf("select new country_ref: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c newCountry
		if err := rows.Scan(&c.id, &c.name); err != nil {
			return nil, fmt.Errorf("scan new country_ref: %w", err)
		}
		key := strings.ToLower(strings.TrimSpace(c.name))
		if key != "" {
			newByName[key] = c.id
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate new country_ref: %w", err)
	}

	result := make(map[int32]int16)
	var mapped, missing int64

	rows, err = oldDB.QueryContext(ctx, `
        SELECT "CountryID", "CountryName"
        FROM "References"."CountryRef"
    `)
	if err != nil {
		return nil, fmt.Errorf("select old CountryRef: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var oldID int32
		var name string
		if err := rows.Scan(&oldID, &name); err != nil {
			return nil, fmt.Errorf("scan old CountryRef: %w", err)
		}
		key := strings.ToLower(strings.TrimSpace(name))
		if newID, ok := newByName[key]; ok {
			result[oldID] = newID
			mapped++
		} else {
			if missing < 20 {
				log.Printf("WARN: buildCountryIDMap: no new country_ref for old CountryID=%d name=%q", oldID, name)
			}
			missing++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate old CountryRef: %w", err)
	}

	log.Printf("buildCountryIDMap: mapped %d old countries; %d missing", mapped, missing)
	return result, nil
}

// old: "References"."LanguageRef"(LanguageID, LanguageName, LanguageCode)
// new: language_ref(id, name, iso_code)
func buildLanguageIDMap(ctx context.Context, oldDB, newDB *sql.DB) (map[int32]int16, error) {
	type newLang struct {
		id   int16
		name string
		code string
	}
	byCode := make(map[string]int16)
	byName := make(map[string]int16)

	rows, err := newDB.QueryContext(ctx, `SELECT id, name, iso_code FROM language_ref`)
	if err != nil {
		return nil, fmt.Errorf("select new language_ref: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l newLang
		if err := rows.Scan(&l.id, &l.name, &l.code); err != nil {
			return nil, fmt.Errorf("scan new language_ref: %w", err)
		}
		codeKey := strings.ToLower(strings.TrimSpace(l.code))
		nameKey := strings.ToLower(strings.TrimSpace(l.name))
		if codeKey != "" {
			byCode[codeKey] = l.id
		}
		if nameKey != "" {
			byName[nameKey] = l.id
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate new language_ref: %w", err)
	}

	result := make(map[int32]int16)
	var mapped, missing int64

	rows, err = oldDB.QueryContext(ctx, `
        SELECT "LanguageID", "LanguageName", "LanguageCode"
        FROM "References"."LanguageRef"
    `)
	if err != nil {
		return nil, fmt.Errorf("select old LanguageRef: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var oldID int32
		var name, code string
		if err := rows.Scan(&oldID, &name, &code); err != nil {
			return nil, fmt.Errorf("scan old LanguageRef: %w", err)
		}
		codeKey := strings.ToLower(strings.TrimSpace(code))
		nameKey := strings.ToLower(strings.TrimSpace(name))

		var newID int16
		var ok bool

		if codeKey != "" && codeKey != "undefined" {
			newID, ok = byCode[codeKey]
		}
		if !ok && nameKey != "" {
			newID, ok = byName[nameKey]
		}

		if ok {
			result[oldID] = newID
			mapped++
		} else {
			if missing < 20 {
				log.Printf("WARN: buildLanguageIDMap: no new language_ref.id for old LanguageID=%d name=%q code=%q", oldID, name, code)
			}
			missing++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate old LanguageRef: %w", err)
	}

	log.Printf("buildLanguageIDMap: mapped %d old languages; %d missing", mapped, missing)
	return result, nil
}

// old: "References"."GenreRef"(GenreID, GenreName)
// new: genre_ref(id, name)
func buildGenreIDMap(ctx context.Context, oldDB, newDB *sql.DB) (map[int32]int16, error) {
	newByName := make(map[string]int16)

	rows, err := newDB.QueryContext(ctx, `SELECT id, name FROM genre_ref`)
	if err != nil {
		return nil, fmt.Errorf("select new genre_ref: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int16
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("scan new genre_ref: %w", err)
		}
		key := strings.ToLower(strings.TrimSpace(name))
		if key != "" {
			newByName[key] = id
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate new genre_ref: %w", err)
	}

	result := make(map[int32]int16)
	var mapped, missing int64

	rows, err = oldDB.QueryContext(ctx, `
        SELECT "GenreID", "GenreName"
        FROM "References"."GenreRef"
    `)
	if err != nil {
		return nil, fmt.Errorf("select old GenreRef: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var oldID int32
		var name string
		if err := rows.Scan(&oldID, &name); err != nil {
			return nil, fmt.Errorf("scan old GenreRef: %w", err)
		}
		key := strings.ToLower(strings.TrimSpace(name))
		if newID, ok := newByName[key]; ok {
			result[oldID] = newID
			mapped++
		} else {
			if missing < 20 {
				log.Printf("WARN: buildGenreIDMap: no new genre_ref.id for old GenreID=%d name=%q", oldID, name)
			}
			missing++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate old GenreRef: %w", err)
	}

	log.Printf("buildGenreIDMap: mapped %d old genres; %d missing", mapped, missing)
	return result, nil
}

// old: "References"."CertificateRef"(CertificateID, CertificateName)
// new: certificate_ref(id, name, description)
func buildCertificateIDMap(ctx context.Context, oldDB, newDB *sql.DB) (map[int32]int16, error) {
	newByName := make(map[string]int16)

	rows, err := newDB.QueryContext(ctx, `SELECT id, name FROM certificate_ref`)
	if err != nil {
		return nil, fmt.Errorf("select new certificate_ref: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int16
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("scan new certificate_ref: %w", err)
		}
		key := strings.ToLower(strings.TrimSpace(name))
		if key != "" {
			newByName[key] = id
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate new certificate_ref: %w", err)
	}

	result := make(map[int32]int16)
	var mapped, missing int64

	rows, err = oldDB.QueryContext(ctx, `
        SELECT "CertificateID", "CertificateName"
        FROM "References"."CertificateRef"
    `)
	if err != nil {
		return nil, fmt.Errorf("select old CertificateRef: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var oldID int32
		var name string
		if err := rows.Scan(&oldID, &name); err != nil {
			return nil, fmt.Errorf("scan old CertificateRef: %w", err)
		}
		key := strings.ToLower(strings.TrimSpace(name))
		if newID, ok := newByName[key]; ok {
			result[oldID] = newID
			mapped++
		} else {
			if missing < 20 {
				log.Printf("WARN: buildCertificateIDMap: no new certificate_ref.id for old CertificateID=%d name=%q", oldID, name)
			}
			missing++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate old CertificateRef: %w", err)
	}

	log.Printf("buildCertificateIDMap: mapped %d old certificates; %d missing", mapped, missing)
	return result, nil
}

//
// Actual junction migrations
//

// Assumes new title.id == old Tables."TitleTable"."TitleID" (as in your core-title migration)

// CountryTitleLine -> title_country
func migrateTitleCountry(
	ctx context.Context,
	oldDB, newDB *sql.DB,
	countryIDMap map[int32]int16,
	dryRun bool,
) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "Lines"."CountryTitleLine"`).Scan(&total); err != nil {
		return fmt.Errorf("count CountryTitleLine: %w", err)
	}
	log.Printf("migrateTitleCountry: %d rows in Lines.\"CountryTitleLine\"", total)

	if dryRun {
		log.Printf("migrateTitleCountry [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
        SELECT "TitleID", "CountryID"
        FROM "Lines"."CountryTitleLine"
        ORDER BY "TitleID"
    `)
	if err != nil {
		return fmt.Errorf("query CountryTitleLine: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx (title_country): %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO title_country (title_id, country_id)
        VALUES ($1, $2)
        ON CONFLICT (title_id, country_id) DO NOTHING
    `)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare insert title_country: %w", err)
	}
	defer stmt.Close()

	var processed int64
	var skipped int64

	for rows.Next() {
		var titleID int32
		var oldCountryID int32
		if err := rows.Scan(&titleID, &oldCountryID); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("scan CountryTitleLine: %w", err)
		}

		newCountryID, ok := countryIDMap[oldCountryID]
		if !ok {
			// No mapped country -> skip
			skipped++
			continue
		}

		if _, err := stmt.ExecContext(ctx, titleID, newCountryID); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert title_country title_id=%d oldCountryID=%d -> newCountryID=%d: %w",
				titleID, oldCountryID, newCountryID, err)
		}

		processed++
		if processed%junctionProgressEvery == 0 && total > 0 {
			pct := float64(processed) * 100.0 / float64(total)
			log.Printf("migrateTitleCountry: inserted %d/%d rows (%.1f%%)", processed, total, pct)
		}
	}
	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("iterate CountryTitleLine: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit title_country: %w", err)
	}
	log.Printf("--- Done title_country: %d rows processed, %d skipped (no country mapping) ---", processed, skipped)
	return nil
}

// LanguageTitleLine -> title_language
// NOTE: old schema has ONLY (LanguageTitleLineID, TitleID, LanguageID).
//       There is NO "IsOriginalLanguage" column, so we set is_original = false for all rows.
func migrateTitleLanguage(
	ctx context.Context,
	oldDB, newDB *sql.DB,
	langIDMap map[int32]int16,
	dryRun bool,
) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "Lines"."LanguageTitleLine"`).Scan(&total); err != nil {
		return fmt.Errorf("count LanguageTitleLine: %w", err)
	}
	log.Printf("migrateTitleLanguage: %d rows in Lines.\"LanguageTitleLine\"", total)

	if dryRun {
		log.Printf("migrateTitleLanguage [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
        SELECT "TitleID", "LanguageID"
        FROM "Lines"."LanguageTitleLine"
        ORDER BY "TitleID"
    `)
	if err != nil {
		return fmt.Errorf("query LanguageTitleLine: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx (title_language): %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO title_language (title_id, language_id, is_original)
        VALUES ($1, $2, $3)
        ON CONFLICT (title_id, language_id) DO NOTHING
    `)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare insert title_language: %w", err)
	}
	defer stmt.Close()

	var processed int64
	var skipped int64

	for rows.Next() {
		var titleID int32
		var oldLangID int32
		if err := rows.Scan(&titleID, &oldLangID); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("scan LanguageTitleLine: %w", err)
		}

		newLangID, ok := langIDMap[oldLangID]
		if !ok {
			skipped++
			continue
		}

		isOriginal := false // old DB has no flag; default to false

		if _, err := stmt.ExecContext(ctx, titleID, newLangID, isOriginal); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert title_language title_id=%d oldLanguageID=%d -> newLanguageID=%d: %w",
				titleID, oldLangID, newLangID, err)
		}

		processed++
		if processed%junctionProgressEvery == 0 && total > 0 {
			pct := float64(processed) * 100.0 / float64(total)
			log.Printf("migrateTitleLanguage: inserted %d/%d rows (%.1f%%)", processed, total, pct)
		}
	}
	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("iterate LanguageTitleLine: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit title_language: %w", err)
	}
	log.Printf("--- Done title_language: %d rows processed, %d skipped (no language mapping) ---", processed, skipped)
	return nil
}

// GenreTitleLine -> title_genre
func migrateTitleGenre(
	ctx context.Context,
	oldDB, newDB *sql.DB,
	genreIDMap map[int32]int16,
	dryRun bool,
) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "Lines"."GenreTitleLine"`).Scan(&total); err != nil {
		return fmt.Errorf("count GenreTitleLine: %w", err)
	}
	log.Printf("migrateTitleGenre: %d rows in Lines.\"GenreTitleLine\"", total)

	if dryRun {
		log.Printf("migrateTitleGenre [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
        SELECT "TitleID", "GenreID"
        FROM "Lines"."GenreTitleLine"
        ORDER BY "TitleID"
    `)
	if err != nil {
		return fmt.Errorf("query GenreTitleLine: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx (title_genre): %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO title_genre (title_id, genre_id)
        VALUES ($1, $2)
        ON CONFLICT (title_id, genre_id) DO NOTHING
    `)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare insert title_genre: %w", err)
	}
	defer stmt.Close()

	var processed int64
	var skipped int64

	for rows.Next() {
		var titleID int32
		var oldGenreID int32
		if err := rows.Scan(&titleID, &oldGenreID); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("scan GenreTitleLine: %w", err)
		}

		newGenreID, ok := genreIDMap[oldGenreID]
		if !ok {
			skipped++
			continue
		}

		if _, err := stmt.ExecContext(ctx, titleID, newGenreID); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert title_genre title_id=%d oldGenreID=%d -> newGenreID=%d: %w",
				titleID, oldGenreID, newGenreID, err)
		}

		processed++
		if processed%junctionProgressEvery == 0 && total > 0 {
			pct := float64(processed) * 100.0 / float64(total)
			log.Printf("migrateTitleGenre: inserted %d/%d rows (%.1f%%)", processed, total, pct)
		}
	}
	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("iterate GenreTitleLine: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit title_genre: %w", err)
	}
	log.Printf("--- Done title_genre: %d rows processed, %d skipped (no genre mapping) ---", processed, skipped)
	return nil
}

// CertificateTitleLine -> title_certificate
func migrateTitleCertificate(
	ctx context.Context,
	oldDB, newDB *sql.DB,
	countryIDMap map[int32]int16,
	certIDMap map[int32]int16,
	dryRun bool,
) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "Lines"."CertificateTitleLine"`).Scan(&total); err != nil {
		return fmt.Errorf("count CertificateTitleLine: %w", err)
	}
	log.Printf("migrateTitleCertificate: %d rows in Lines.\"CertificateTitleLine\"", total)

	if dryRun {
		log.Printf("migrateTitleCertificate [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
        SELECT "TitleID", "CertificateID", "CountryID"
        FROM "Lines"."CertificateTitleLine"
        ORDER BY "TitleID"
    `)
	if err != nil {
		return fmt.Errorf("query CertificateTitleLine: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx (title_certificate): %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO title_certificate (title_id, certificate_id, country_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (title_id, certificate_id, country_id) DO NOTHING
    `)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare insert title_certificate: %w", err)
	}
	defer stmt.Close()

	var processed int64
	var skipped int64

	for rows.Next() {
		var titleID int32
		var oldCertID int32
		var oldCountryID int32
		if err := rows.Scan(&titleID, &oldCertID, &oldCountryID); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("scan CertificateTitleLine: %w", err)
		}

		newCertID, okCert := certIDMap[oldCertID]
		newCountryID, okCountry := countryIDMap[oldCountryID]
		if !okCert || !okCountry {
			skipped++
			continue
		}

		if _, err := stmt.ExecContext(ctx, titleID, newCertID, newCountryID); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert title_certificate title_id=%d oldCertID=%d->%d oldCountryID=%d->%d: %w",
				titleID, oldCertID, newCertID, oldCountryID, newCountryID, err)
		}

		processed++
		if processed%junctionProgressEvery == 0 && total > 0 {
			pct := float64(processed) * 100.0 / float64(total)
			log.Printf("migrateTitleCertificate: inserted %d/%d rows (%.1f%%)", processed, total, pct)
		}
	}
	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("iterate CertificateTitleLine: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit title_certificate: %w", err)
	}
	log.Printf("--- Done title_certificate: %d rows processed, %d skipped (no cert/country mapping) ---", processed, skipped)
	return nil
}
