package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// MigrateRefsPhase migrates all lookup / reference tables that do not depend
// on person/title rows.
func MigrateRefsPhase(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("=== Starting reference data migration [%s] ===", modeLabel(dryRun))

	steps := []struct {
		name string
		fn   func(context.Context, *sql.DB, *sql.DB, bool) error
	}{
		{"country_ref", migrateCountryRef},
		{"language_ref", migrateLanguageRef},
		{"genre_ref", migrateGenreRef},
		{"certificate_ref", migrateCertificateRef},
		{"title_type_ref", migrateTitleTypeRef},
		{"connection_type_ref", migrateConnectionTypeRef},
		{"parental_guide_category_ref", migrateParentalGuideRef},
		{"quality_ref", migrateQualityRef},
		{"display_ref", migrateDisplayRef},
		{"cast_role_type_ref", migrateCastRoleTypeRef},
		{"award_event_ref", migrateAwardEventRef},
		{"award_nomination_type_ref", migrateAwardNominationTypeRef},
		{"certificate_country", migrateCertificateCountry},
	}

	start := time.Now()
	for _, step := range steps {
		log.Printf("--- Migrating %s ---", step.name)
		stepStart := time.Now()

		if err := step.fn(ctx, oldDB, newDB, dryRun); err != nil {
			return fmt.Errorf("migration step %s failed: %w", step.name, err)
		}

		log.Printf("--- Done %s in %s ---", step.name, time.Since(stepStart))
	}

	log.Printf("=== All reference migrations completed in %s [%s] ===",
		time.Since(start), modeLabel(dryRun))

	return nil
}

func modeLabel(dryRun bool) string {
	if dryRun {
		return "DRY-RUN (no writes)"
	}
	return "LIVE"
}

// migrateCountryRef migrates Countries."CountryRef" -> country_ref.
func migrateCountryRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "CountryID", "CountryName", "CountryCode"
		FROM "Countries"."CountryRef"
		ORDER BY "CountryID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query CountryRef: %w", err)
	}
	defer rows.Close()

	type countryRow struct {
		id   int64
		name string
		code sql.NullString
	}

	var allRows []countryRow
	for rows.Next() {
		var r countryRow
		if err := rows.Scan(&r.id, &r.name, &r.code); err != nil {
			return fmt.Errorf("scan CountryRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CountryRef: %w", err)
	}

	log.Printf("migrateCountryRef: read %d rows from Countries.\"CountryRef\"", len(allRows))

	if dryRun {
		// Just log some stats and return.
		var nonISO int
		for _, r := range allRows {
			code := strings.TrimSpace(strings.ToLower(r.code.String))
			if code == "" {
				continue
			}
			// IMDb-style 4-letter codes for historical countries, etc.
			if len(code) != 2 && len(code) != 3 {
				nonISO++
			}
		}
		log.Printf("migrateCountryRef [DRY-RUN]: %d rows total, %d rows have non-ISO-like codes", len(allRows), nonISO)
		return nil
	}

	// Insert into new.country_ref
	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx country_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO country_ref (id, name, iso2, iso3)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    iso2 = EXCLUDED.iso2,
		    iso3 = EXCLUDED.iso3
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert country_ref: %w", err)
	}
	defer stmt.Close()

	var nonISO int
	for _, r := range allRows {
		var iso2, iso3 sql.NullString
		code := strings.TrimSpace(strings.ToLower(r.code.String))
		if code != "" {
			switch len(code) {
			case 2:
				iso2 = sql.NullString{String: strings.ToUpper(code), Valid: true}
			case 3:
				iso3 = sql.NullString{String: strings.ToUpper(code), Valid: true}
			default:
				nonISO++
				log.Printf("WARN: country_ref id=%d name=%q has non-ISO code %q (len=%d); inserting with NULL iso2/iso3",
					r.id, r.name, code, len(code))
			}
		}

		if _, err := stmt.ExecContext(ctx, r.id, r.name, iso2, iso3); err != nil {
			return fmt.Errorf("insert country_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit country_ref: %w", err)
	}

	if nonISO > 0 {
		log.Printf("migrateCountryRef: %d rows had non-ISO codes; inserted with NULL iso2/iso3", nonISO)
	}

	return nil
}

// migrateLanguageRef migrates Languages."LanguageRef" -> language_ref.
func migrateLanguageRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "LanguageID", "LanguageName", "LanguageCode"
		FROM "Languages"."LanguageRef"
		ORDER BY "LanguageID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query LanguageRef: %w", err)
	}
	defer rows.Close()

	type langRow struct {
		id   int64
		name string
		code sql.NullString
	}

	var allRows []langRow
	for rows.Next() {
		var r langRow
		if err := rows.Scan(&r.id, &r.name, &r.code); err != nil {
			return fmt.Errorf("scan LanguageRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate LanguageRef: %w", err)
	}

	log.Printf("migrateLanguageRef: read %d rows from Languages.\"LanguageRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx language_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO language_ref (id, name, code)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    code = EXCLUDED.code
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert language_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		var code sql.NullString
		if trimmed := strings.TrimSpace(r.code.String); trimmed != "" && trimmed != "Undefined" {
			code = sql.NullString{String: trimmed, Valid: true}
		}
		if _, err := stmt.ExecContext(ctx, r.id, r.name, code); err != nil {
			return fmt.Errorf("insert language_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit language_ref: %w", err)
	}

	return nil
}

// migrateGenreRef migrates Genres."GenreRef" -> genre_ref.
func migrateGenreRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "GenreID", "GenreName"
		FROM "Genres"."GenreRef"
		ORDER BY "GenreID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query GenreRef: %w", err)
	}
	defer rows.Close()

	type genreRow struct {
		id   int64
		name string
	}

	var allRows []genreRow
	for rows.Next() {
		var r genreRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan GenreRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate GenreRef: %w", err)
	}

	log.Printf("migrateGenreRef: read %d rows from Genres.\"GenreRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx genre_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO genre_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert genre_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert genre_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit genre_ref: %w", err)
	}

	return nil
}

// migrateCertificateRef migrates Certificates."CertificateRef" -> certificate_ref.
func migrateCertificateRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "CertificateID", "CertificateName"
		FROM "Certificates"."CertificateRef"
		ORDER BY "CertificateID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query CertificateRef: %w", err)
	}
	defer rows.Close()

	type certRow struct {
		id   int64
		name string
	}

	var allRows []certRow
	for rows.Next() {
		var r certRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan CertificateRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CertificateRef: %w", err)
	}

	log.Printf("migrateCertificateRef: read %d rows from Certificates.\"CertificateRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx certificate_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO certificate_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert certificate_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert certificate_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit certificate_ref: %w", err)
	}

	return nil
}

// migrateTitleTypeRef migrates TitleTypes."TitleTypeRef" -> title_type_ref.
func migrateTitleTypeRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "TitleTypeID", "TitleTypeName"
		FROM "TitleTypes"."TitleTypeRef"
		ORDER BY "TitleTypeID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query TitleTypeRef: %w", err)
	}
	defer rows.Close()

	type ttRow struct {
		id   int64
		name string
	}

	var allRows []ttRow
	for rows.Next() {
		var r ttRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan TitleTypeRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate TitleTypeRef: %w", err)
	}

	log.Printf("migrateTitleTypeRef: read %d rows from TitleTypes.\"TitleTypeRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx title_type_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO title_type_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert title_type_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert title_type_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit title_type_ref: %w", err)
	}

	return nil
}

// migrateConnectionTypeRef migrates Connections."ConnectionTypeRef" -> connection_type_ref.
func migrateConnectionTypeRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "ConnectionTypeID", "ConnectionTypeName"
		FROM "Connections"."ConnectionTypeRef"
		ORDER BY "ConnectionTypeID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query ConnectionTypeRef: %w", err)
	}
	defer rows.Close()

	type connRow struct {
		id   int64
		name string
	}

	var allRows []connRow
	for rows.Next() {
		var r connRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan ConnectionTypeRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate ConnectionTypeRef: %w", err)
	}

	log.Printf("migrateConnectionTypeRef: read %d rows from Connections.\"ConnectionTypeRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx connection_type_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO connection_type_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert connection_type_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert connection_type_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit connection_type_ref: %w", err)
	}

	return nil
}

// migrateParentalGuideRef migrates ParentsGuide."ParentsGuideCategoryRef" -> parental_guide_category_ref.
func migrateParentalGuideRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "ParentsGuideCategoryID", "ParentsGuideCategoryName"
		FROM "ParentsGuide"."ParentsGuideCategoryRef"
		ORDER BY "ParentsGuideCategoryID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query ParentsGuideCategoryRef: %w", err)
	}
	defer rows.Close()

	type pgRow struct {
		id   int64
		name string
	}

	var allRows []pgRow
	for rows.Next() {
		var r pgRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan ParentsGuideCategoryRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate ParentsGuideCategoryRef: %w", err)
	}

	log.Printf("migrateParentalGuideRef: read %d rows from ParentsGuide.\"ParentsGuideCategoryRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx parental_guide_category_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO parental_guide_category_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert parental_guide_category_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert parental_guide_category_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit parental_guide_category_ref: %w", err)
	}

	return nil
}

// migrateQualityRef migrates Quality."QualityRef" -> quality_ref.
func migrateQualityRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "QualityID", "QualityName"
		FROM "Quality"."QualityRef"
		ORDER BY "QualityID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query QualityRef: %w", err)
	}
	defer rows.Close()

	type qRow struct {
		id   int64
		name string
	}

	var allRows []qRow
	for rows.Next() {
		var r qRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan QualityRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate QualityRef: %w", err)
	}

	log.Printf("migrateQualityRef: read %d rows from Quality.\"QualityRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx quality_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO quality_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert quality_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert quality_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit quality_ref: %w", err)
	}

	return nil
}

// migrateDisplayRef migrates Displays."DisplayRef" -> display_ref.
func migrateDisplayRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "DisplayID", "DisplayName"
		FROM "Displays"."DisplayRef"
		ORDER BY "DisplayID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query DisplayRef: %w", err)
	}
	defer rows.Close()

	type dRow struct {
		id   int64
		name string
	}

	var allRows []dRow
	for rows.Next() {
		var r dRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan DisplayRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate DisplayRef: %w", err)
	}

	log.Printf("migrateDisplayRef: read %d rows from Displays.\"DisplayRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx display_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO display_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert display_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert display_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit display_ref: %w", err)
	}

	return nil
}

// migrateCastRoleTypeRef migrates Cast."CastRoleTypeRef" -> cast_role_type_ref.
func migrateCastRoleTypeRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "CastRoleTypeID", "CastRoleTypeName"
		FROM "Cast"."CastRoleTypeRef"
		ORDER BY "CastRoleTypeID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query CastRoleTypeRef: %w", err)
	}
	defer rows.Close()

	type crtRow struct {
		id   int64
		name string
	}

	var allRows []crtRow
	for rows.Next() {
		var r crtRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan CastRoleTypeRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CastRoleTypeRef: %w", err)
	}

	log.Printf("migrateCastRoleTypeRef: read %d rows from Cast.\"CastRoleTypeRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx cast_role_type_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO cast_role_type_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert cast_role_type_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert cast_role_type_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit cast_role_type_ref: %w", err)
	}

	return nil
}

// migrateAwardEventRef migrates Awards."AwardEventRef" -> award_event_ref.
func migrateAwardEventRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "AwardEventID", "AwardEventName"
		FROM "Awards"."AwardEventRef"
		ORDER BY "AwardEventID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query AwardEventRef: %w", err)
	}
	defer rows.Close()

	type aeRow struct {
		id   int64
		name string
	}

	var allRows []aeRow
	for rows.Next() {
		var r aeRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan AwardEventRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate AwardEventRef: %w", err)
	}

	log.Printf("migrateAwardEventRef: read %d rows from Awards.\"AwardEventRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx award_event_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO award_event_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert award_event_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert award_event_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit award_event_ref: %w", err)
	}

	return nil
}

// migrateAwardNominationTypeRef migrates Awards."AwardNominationTypeRef" -> award_nomination_type_ref.
func migrateAwardNominationTypeRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "AwardNominationTypeID", "AwardNominationTypeName"
		FROM "Awards"."AwardNominationTypeRef"
		ORDER BY "AwardNominationTypeID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query AwardNominationTypeRef: %w", err)
	}
	defer rows.Close()

	type antRow struct {
		id   int64
		name string
	}

	var allRows []antRow
	for rows.Next() {
		var r antRow
		if err := rows.Scan(&r.id, &r.name); err != nil {
			return fmt.Errorf("scan AwardNominationTypeRef row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate AwardNominationTypeRef: %w", err)
	}

	log.Printf("migrateAwardNominationTypeRef: read %d rows from Awards.\"AwardNominationTypeRef\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx award_nomination_type_ref: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO award_nomination_type_ref (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert award_nomination_type_ref: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.name); err != nil {
			return fmt.Errorf("insert award_nomination_type_ref id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit award_nomination_type_ref: %w", err)
	}

	return nil
}

// migrateCertificateCountry migrates Certificates."CertificateCountry" -> certificate_country.
func migrateCertificateCountry(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	const srcQuery = `
		SELECT "CertificateCountryID", "CertificateID", "CountryID"
		FROM "Certificates"."CertificateCountry"
		ORDER BY "CertificateCountryID"
	`

	rows, err := oldDB.QueryContext(ctx, srcQuery)
	if err != nil {
		return fmt.Errorf("query CertificateCountry: %w", err)
	}
	defer rows.Close()

	type ccRow struct {
		id           int64
		certificateID int64
		countryID    int64
	}

	var allRows []ccRow
	for rows.Next() {
		var r ccRow
		if err := rows.Scan(&r.id, &r.certificateID, &r.countryID); err != nil {
			return fmt.Errorf("scan CertificateCountry row: %w", err)
		}
		allRows = append(allRows, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CertificateCountry: %w", err)
	}

	log.Printf("migrateCertificateCountry: read %d rows from Certificates.\"CertificateCountry\"", len(allRows))

	if dryRun {
		return nil
	}

	tx, err := newDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx certificate_country: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `
		INSERT INTO certificate_country (id, certificate_id, country_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET certificate_id = EXCLUDED.certificate_id,
		    country_id = EXCLUDED.country_id
	`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("prepare insert certificate_country: %w", err)
	}
	defer stmt.Close()

	for _, r := range allRows {
		if _, err := stmt.ExecContext(ctx, r.id, r.certificateID, r.countryID); err != nil {
			return fmt.Errorf("insert certificate_country id=%d: %w", r.id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit certificate_country: %w", err)
	}

	return nil
}
