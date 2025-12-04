package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// runReferenceMigration migrates all lookup / reference tables that do not depend
// on person/title rows.
func runReferenceMigration(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	log.Printf("=== Starting reference data migration [DRY-RUN=%v] ===", dryRun)
	start := time.Now()

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

	for _, step := range steps {
		log.Printf("--- Migrating %s ---", step.name)
		s := time.Now()
		if err := step.fn(ctx, oldDB, newDB, dryRun); err != nil {
			return fmt.Errorf("migration step %s failed: %w", step.name, err)
		}
		log.Printf("--- Done %s in %s ---", step.name, time.Since(s).Truncate(time.Millisecond))
	}

	log.Printf("=== All reference migrations completed in %s [DRY-RUN=%v] ===",
		time.Since(start).Truncate(time.Millisecond), dryRun)
	return nil
}

//
// country_ref
//

func migrateCountryRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."CountryRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old CountryRef: %w", err)
	}

	if dryRun {
		log.Printf("migrateCountryRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "CountryName", "CountryCode"
		FROM "References"."CountryRef"
		ORDER BY "CountryID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query CountryRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx country_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO country_ref (name, iso2_code, iso3_code)
		VALUES ($1, $2, $3)
		ON CONFLICT (name) DO UPDATE
		SET
			iso2_code = EXCLUDED.iso2_code,
			iso3_code = EXCLUDED.iso3_code
	`)
	if err != nil {
		return fmt.Errorf("prepare country_ref insert: %w", err)
	}
	defer stmt.Close()

	var (
		processed   int64
		nonISOCount int64
	)

	for rows.Next() {
		var (
			name string
			code sql.NullString
		)
		if err := rows.Scan(&name, &code); err != nil {
			return fmt.Errorf("scan CountryRef row: %w", err)
		}

		var iso2, iso3 *string
		if code.Valid {
			c := strings.TrimSpace(strings.ToLower(code.String))
			switch len(c) {
			case 2:
				up := strings.ToUpper(c)
				iso2 = &up
			case 3:
				up := strings.ToUpper(c)
				iso3 = &up
			default:
				nonISOCount++
				log.Printf("WARN: country_ref name=%q has non-ISO code %q (len=%d); inserting with NULL iso2/iso3",
					name, c, len(c))
			}
		}

		if _, errExec := stmt.ExecContext(ctx, name, iso2, iso3); errExec != nil {
			return fmt.Errorf("insert country_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CountryRef: %w", err)
	}

	if nonISOCount > 0 {
		log.Printf("migrateCountryRef: %d rows had non-ISO codes; inserted with NULL iso2/iso3", nonISOCount)
	}
	log.Printf("migrateCountryRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit country_ref: %w", errCommit)
	}
	return nil
}

//
// language_ref
//

func migrateLanguageRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."LanguageRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old LanguageRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateLanguageRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "LanguageName", "LanguageCode"
		FROM "References"."LanguageRef"
		ORDER BY "LanguageID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query LanguageRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx language_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO language_ref (name, iso_code)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE
		SET iso_code = EXCLUDED.iso_code
	`)
	if err != nil {
		return fmt.Errorf("prepare language_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name, code string
		if err := rows.Scan(&name, &code); err != nil {
			return fmt.Errorf("scan LanguageRef row: %w", err)
		}
		code = strings.TrimSpace(strings.ToLower(code))
		if code == "" {
			return fmt.Errorf("language_ref name=%q has empty code", name)
		}

		if _, errExec := stmt.ExecContext(ctx, name, code); errExec != nil {
			return fmt.Errorf("insert language_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate LanguageRef: %w", err)
	}

	log.Printf("migrateLanguageRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit language_ref: %w", errCommit)
	}
	return nil
}

//
// genre_ref
//

func migrateGenreRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."GenreRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old GenreRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateGenreRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "GenreName"
		FROM "References"."GenreRef"
		ORDER BY "GenreID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query GenreRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx genre_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO genre_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare genre_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan GenreRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert genre_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate GenreRef: %w", err)
	}

	log.Printf("migrateGenreRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit genre_ref: %w", errCommit)
	}
	return nil
}

//
// certificate_ref
//

func migrateCertificateRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."CertificateRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old CertificateRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateCertificateRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "CertificateName"
		FROM "References"."CertificateRef"
		ORDER BY "CertificateID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query CertificateRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx certificate_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO certificate_ref (name, description)
		VALUES ($1, NULL)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare certificate_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan CertificateRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert certificate_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CertificateRef: %w", err)
	}

	log.Printf("migrateCertificateRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit certificate_ref: %w", errCommit)
	}
	return nil
}

//
// title_type_ref
//

func migrateTitleTypeRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."TitleTypeRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old TitleTypeRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateTitleTypeRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "TypeName"
		FROM "References"."TitleTypeRef"
		ORDER BY "TypeID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query TitleTypeRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx title_type_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO title_type_ref (name, is_series)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE
		SET is_series = (title_type_ref.is_series OR EXCLUDED.is_series)
	`)
	if err != nil {
		return fmt.Errorf("prepare title_type_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan TitleTypeRef row: %w", err)
		}
		lower := strings.ToLower(strings.TrimSpace(name))
		isSeries := strings.Contains(lower, "series") || lower == "episode"

		if _, errExec := stmt.ExecContext(ctx, name, isSeries); errExec != nil {
			return fmt.Errorf("insert title_type_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate TitleTypeRef: %w", err)
	}

	log.Printf("migrateTitleTypeRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit title_type_ref: %w", errCommit)
	}
	return nil
}

//
// connection_type_ref
//

func migrateConnectionTypeRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."ConnectionTypeRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old ConnectionTypeRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateConnectionTypeRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "ConnectionTypeDescription"
		FROM "References"."ConnectionTypeRef"
		ORDER BY "ConnectionTypeID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query ConnectionTypeRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx connection_type_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO connection_type_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare connection_type_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan ConnectionTypeRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert connection_type_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate ConnectionTypeRef: %w", err)
	}

	log.Printf("migrateConnectionTypeRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit connection_type_ref: %w", errCommit)
	}
	return nil
}

//
// parental_guide_category_ref
//

func migrateParentalGuideRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."ParentGuideRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old ParentGuideRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateParentalGuideRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "ParentGuideDescription"
		FROM "References"."ParentGuideRef"
		ORDER BY "ParentGuideID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query ParentGuideRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx parental_guide_category_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO parental_guide_category_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare parental_guide_category_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan ParentGuideRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert parental_guide_category_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate ParentGuideRef: %w", err)
	}

	log.Printf("migrateParentalGuideRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit parental_guide_category_ref: %w", errCommit)
	}
	return nil
}

//
// quality_ref
//

func migrateQualityRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."QualityRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old QualityRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateQualityRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "QualityName"
		FROM "References"."QualityRef"
		ORDER BY "QualityID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query QualityRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx quality_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO quality_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare quality_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan QualityRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert quality_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate QualityRef: %w", err)
	}

	log.Printf("migrateQualityRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit quality_ref: %w", errCommit)
	}
	return nil
}

//
// display_ref
//

func migrateDisplayRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."DisplayRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old DisplayRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateDisplayRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "DisplayType"
		FROM "References"."DisplayRef"
		ORDER BY "DisplayID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query DisplayRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx display_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO display_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare display_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan DisplayRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert display_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate DisplayRef: %w", err)
	}

	log.Printf("migrateDisplayRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit display_ref: %w", errCommit)
	}
	return nil
}

//
// cast_role_type_ref
//

func migrateCastRoleTypeRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."CastTypeRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old CastTypeRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateCastRoleTypeRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "CastTypeDescription"
		FROM "References"."CastTypeRef"
		ORDER BY "CastTypeID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query CastTypeRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx cast_role_type_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO cast_role_type_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare cast_role_type_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan CastTypeRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert cast_role_type_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CastTypeRef: %w", err)
	}

	log.Printf("migrateCastRoleTypeRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit cast_role_type_ref: %w", errCommit)
	}
	return nil
}

//
// award_event_ref
//

func migrateAwardEventRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."AwardEventRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old AwardEventRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateAwardEventRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "EventName"
		FROM "References"."AwardEventRef"
		ORDER BY "EventID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query AwardEventRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx award_event_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO award_event_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare award_event_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan AwardEventRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert award_event_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate AwardEventRef: %w", err)
	}

	log.Printf("migrateAwardEventRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit award_event_ref: %w", errCommit)
	}
	return nil
}

//
// award_nomination_type_ref
//

func migrateAwardNominationTypeRef(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."AwardNominationTypeRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old AwardNominationTypeRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateAwardNominationTypeRef [DRY-RUN]: would process %d rows", total)
		return nil
	}

	rows, err := oldDB.QueryContext(ctx, `
		SELECT "NominationType"
		FROM "References"."AwardNominationTypeRef"
		ORDER BY "NominationTypeID" ASC
	`)
	if err != nil {
		return fmt.Errorf("query AwardNominationTypeRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx award_nomination_type_ref: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO award_nomination_type_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare award_nomination_type_ref insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan AwardNominationTypeRef row: %w", err)
		}
		if _, errExec := stmt.ExecContext(ctx, name); errExec != nil {
			return fmt.Errorf("insert award_nomination_type_ref name=%q: %w", name, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate AwardNominationTypeRef: %w", err)
	}

	log.Printf("migrateAwardNominationTypeRef: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit award_nomination_type_ref: %w", errCommit)
	}
	return nil
}

//
// certificate_country
//

func migrateCertificateCountry(ctx context.Context, oldDB, newDB *sql.DB, dryRun bool) error {
	var total int64
	if err := oldDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "References"."CertificateCountryRef"`).Scan(&total); err != nil {
		return fmt.Errorf("count old CertificateCountryRef: %w", err)
	}
	if dryRun {
		log.Printf("migrateCertificateCountry [DRY-RUN]: would process %d rows", total)
		return nil
	}

	// We assume that IDs in old reference tables map 1:1 to rows we migrated already
	// into certificate_ref and country_ref (since they only had one row per ID).
	// So we can safely reuse them here.
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "CountryID", "CertificateID", "Age"
		FROM "References"."CertificateCountryRef"
		ORDER BY "CountryID", "CertificateID"
	`)
	if err != nil {
		return fmt.Errorf("query CertificateCountryRef: %w", err)
	}
	defer rows.Close()

	tx, err := newDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx certificate_country: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO certificate_country (country_id, certificate_id, min_age)
		VALUES ($1, $2, $3)
		ON CONFLICT (country_id, certificate_id) DO UPDATE
		SET min_age = EXCLUDED.min_age
	`)
	if err != nil {
		return fmt.Errorf("prepare certificate_country insert: %w", err)
	}
	defer stmt.Close()

	var processed int64
	for rows.Next() {
		var (
			countryID     int32
			certificateID int32
			age           sql.NullInt32
		)
		if err := rows.Scan(&countryID, &certificateID, &age); err != nil {
			return fmt.Errorf("scan CertificateCountryRef row: %w", err)
		}

		var minAge *int32
		if age.Valid {
			v := age.Int32
			minAge = &v
		}

		if _, errExec := stmt.ExecContext(ctx, countryID, certificateID, minAge); errExec != nil {
			return fmt.Errorf("insert certificate_country (%d,%d): %w", countryID, certificateID, errExec)
		}
		processed++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CertificateCountryRef: %w", err)
	}

	log.Printf("migrateCertificateCountry: %d rows processed", processed)
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("commit certificate_country: %w", errCommit)
	}
	return nil
}
