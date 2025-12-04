// cmd/migrate-old-db/main.go
//
// Small migration tool to copy reference data from the old "mediadb"
// (LabVIEW-era schema) into the new "movies3db" schema defined in
// db/new/schema.sql.
//
// Usage example:
//
//   go run ./cmd/migrate-old-db \
//     -old "host=127.0.0.1 user=postgres password=YOURPASS dbname=mediadb sslmode=disable" \
//     -new "host=127.0.0.1 user=postgres password=YOURPASS dbname=movies3db sslmode=disable"
//
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	oldDSN := flag.String("old", "", "Postgres DSN for OLD database (mediadb)")
	newDSN := flag.String("new", "", "Postgres DSN for NEW database (movies3db)")
	flag.Parse()

	if *oldDSN == "" || *newDSN == "" {
		log.Fatalf("both -old and -new DSNs are required")
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

	start := time.Now()
	log.Println("=== Starting reference data migration ===")

	steps := []struct {
		name string
		fn   func(context.Context, *sql.DB, *sql.DB) error
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
		{"award_nomination_type_ref", migrateAwardNomTypeRef},
		// certificate_country depends on both country_ref + certificate_ref
		{"certificate_country", migrateCertificateCountry},
	}

	for _, step := range steps {
		log.Printf("--- Migrating %s ---", step.name)
		stepStart := time.Now()
		if err := step.fn(ctx, oldDB, newDB); err != nil {
			log.Fatalf("migration step %s failed: %v", step.name, err)
		}
		log.Printf("--- Done %s in %s ---", step.name, time.Since(stepStart).Truncate(time.Millisecond))
	}

	log.Printf("=== All reference migrations completed in %s ===", time.Since(start).Truncate(time.Millisecond))
}

/*
 * Reference migrations
 */

func migrateCountryRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."CountryRef" (CountryID, CountryName, CountryCode)
	// New: country_ref (id, name, iso2_code, iso3_code)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "CountryName", "CountryCode"
		FROM "References"."CountryRef"
		ORDER BY "CountryID"
	`)
	if err != nil {
		return fmt.Errorf("query old CountryRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO country_ref (name, iso2_code)
		VALUES ($1, NULLIF($2, ''))
		ON CONFLICT (name) DO UPDATE
		SET iso2_code = COALESCE(EXCLUDED.iso2_code, country_ref.iso2_code)
	`)
	if err != nil {
		return fmt.Errorf("prepare insert country_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name, code sql.NullString
		if err := rows.Scan(&name, &code); err != nil {
			return fmt.Errorf("scan CountryRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name.String, code.String)
		if err != nil {
			return fmt.Errorf("insert country_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CountryRef: %w", err)
	}

	log.Printf("migrateCountryRef: inserted/updated %d rows", count)
	return nil
}

func migrateLanguageRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."LanguageRef" (LanguageID, LanguageName, LanguageCode)
	// New: language_ref (id, name, iso_code)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "LanguageName", "LanguageCode"
		FROM "References"."LanguageRef"
		ORDER BY "LanguageID"
	`)
	if err != nil {
		return fmt.Errorf("query old LanguageRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO language_ref (name, iso_code)
		VALUES ($1, $2)
		ON CONFLICT (iso_code) DO UPDATE
		SET name = EXCLUDED.name
	`)
	if err != nil {
		return fmt.Errorf("prepare insert language_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name, code string
		if err := rows.Scan(&name, &code); err != nil {
			return fmt.Errorf("scan LanguageRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name, code)
		if err != nil {
			return fmt.Errorf("insert language_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate LanguageRef: %w", err)
	}

	log.Printf("migrateLanguageRef: inserted/updated %d rows", count)
	return nil
}

func migrateGenreRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."GenreRef" (GenreID, GenreName)
	// New: genre_ref (id, name)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "GenreName"
		FROM "References"."GenreRef"
		ORDER BY "GenreID"
	`)
	if err != nil {
		return fmt.Errorf("query old GenreRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO genre_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert genre_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan GenreRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert genre_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate GenreRef: %w", err)
	}

	log.Printf("migrateGenreRef: inserted %d rows", count)
	return nil
}

func migrateCertificateRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."CertificateRef" (CertificateID, CertificateName)
	// New: certificate_ref (id, name, description)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "CertificateName"
		FROM "References"."CertificateRef"
		ORDER BY "CertificateID"
	`)
	if err != nil {
		return fmt.Errorf("query old CertificateRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO certificate_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert certificate_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan CertificateRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert certificate_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CertificateRef: %w", err)
	}

	log.Printf("migrateCertificateRef: inserted %d rows", count)
	return nil
}

func migrateTitleTypeRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."TitleTypeRef" (TypeID, TypeName)
	// New: title_type_ref (id, name, is_series)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "TypeName"
		FROM "References"."TitleTypeRef"
		ORDER BY "TypeID"
	`)
	if err != nil {
		return fmt.Errorf("query old TitleTypeRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO title_type_ref (name, is_series)
		VALUES ($1, FALSE)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert title_type_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan TitleTypeRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert title_type_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate TitleTypeRef: %w", err)
	}

	log.Printf("migrateTitleTypeRef: inserted %d rows", count)
	return nil
}

func migrateConnectionTypeRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."ConnectionTypeRef" (ConnectionTypeID, ConnectionTypeDescription)
	// New: connection_type_ref (id, name)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "ConnectionTypeDescription"
		FROM "References"."ConnectionTypeRef"
		ORDER BY "ConnectionTypeID"
	`)
	if err != nil {
		return fmt.Errorf("query old ConnectionTypeRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO connection_type_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert connection_type_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan ConnectionTypeRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert connection_type_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate ConnectionTypeRef: %w", err)
	}

	log.Printf("migrateConnectionTypeRef: inserted %d rows", count)
	return nil
}

func migrateParentalGuideRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."ParentGuideRef" (ParentGuideID, ParentGuideDescription)
	// New: parental_guide_category_ref (id, name)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "ParentGuideDescription"
		FROM "References"."ParentGuideRef"
		ORDER BY "ParentGuideID"
	`)
	if err != nil {
		return fmt.Errorf("query old ParentGuideRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO parental_guide_category_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert parental_guide_category_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan ParentGuideRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert parental_guide_category_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate ParentGuideRef: %w", err)
	}

	log.Printf("migrateParentalGuideRef: inserted %d rows", count)
	return nil
}

func migrateQualityRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."QualityRef" (QualityID, QualityName)
	// New: quality_ref (id, name)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "QualityName"
		FROM "References"."QualityRef"
		ORDER BY "QualityID"
	`)
	if err != nil {
		return fmt.Errorf("query old QualityRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO quality_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert quality_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan QualityRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert quality_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate QualityRef: %w", err)
	}

	log.Printf("migrateQualityRef: inserted %d rows", count)
	return nil
}

func migrateDisplayRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."DisplayRef" (DisplayID, DisplayType)
	// New: display_ref (id, name)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "DisplayType"
		FROM "References"."DisplayRef"
		ORDER BY "DisplayID"
	`)
	if err != nil {
		return fmt.Errorf("query old DisplayRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO display_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert display_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan DisplayRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert display_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate DisplayRef: %w", err)
	}

	log.Printf("migrateDisplayRef: inserted %d rows", count)
	return nil
}

func migrateCastRoleTypeRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."CastTypeRef" (CastTypeID, CastTypeDescription)
	// New: cast_role_type_ref (id, name)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "CastTypeDescription"
		FROM "References"."CastTypeRef"
		ORDER BY "CastTypeID"
	`)
	if err != nil {
		return fmt.Errorf("query old CastTypeRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO cast_role_type_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert cast_role_type_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan CastTypeRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert cast_role_type_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CastTypeRef: %w", err)
	}

	log.Printf("migrateCastRoleTypeRef: inserted %d rows", count)
	return nil
}

func migrateAwardEventRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."AwardEventRef" (EventID, EventName)
	// New: award_event_ref (id, name)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "EventName"
		FROM "References"."AwardEventRef"
		ORDER BY "EventID"
	`)
	if err != nil {
		return fmt.Errorf("query old AwardEventRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO award_event_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert award_event_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan AwardEventRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert award_event_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate AwardEventRef: %w", err)
	}

	log.Printf("migrateAwardEventRef: inserted %d rows", count)
	return nil
}

func migrateAwardNomTypeRef(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."AwardNominationTypeRef" (NominationTypeID, NominationType)
	// New: award_nomination_type_ref (id, name)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT "NominationType"
		FROM "References"."AwardNominationTypeRef"
		ORDER BY "NominationTypeID"
	`)
	if err != nil {
		return fmt.Errorf("query old AwardNominationTypeRef: %w", err)
	}
	defer rows.Close()

	stmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO award_nomination_type_ref (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare insert award_nomination_type_ref: %w", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan AwardNominationTypeRef: %w", err)
		}
		_, err := stmt.ExecContext(ctx, name)
		if err != nil {
			return fmt.Errorf("insert award_nomination_type_ref: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate AwardNominationTypeRef: %w", err)
	}

	log.Printf("migrateAwardNomTypeRef: inserted %d rows", count)
	return nil
}

func migrateCertificateCountry(ctx context.Context, oldDB, newDB *sql.DB) error {
	// Old: "References"."CertificateCountryRef"
	//   (CountryID, CertificateID, Age)
	// plus lookups from CountryRef + CertificateRef.
	//
	// New: certificate_country (country_id, certificate_id, min_age)
	rows, err := oldDB.QueryContext(ctx, `
		SELECT
			cc."CountryID",
			cc."CertificateID",
			cc."Age",
			c."CountryName",
			cr."CertificateName"
		FROM "References"."CertificateCountryRef" cc
		JOIN "References"."CountryRef" c
			ON cc."CountryID" = c."CountryID"
		JOIN "References"."CertificateRef" cr
			ON cc."CertificateID" = cr."CertificateID"
		ORDER BY cc."CountryID", cc."CertificateID"
	`)
	if err != nil {
		return fmt.Errorf("query old CertificateCountryRef: %w", err)
	}
	defer rows.Close()

	// Prepare lookup statements on new DB
	getCountryIDStmt, err := newDB.PrepareContext(ctx, `
		SELECT id FROM country_ref WHERE name = $1
	`)
	if err != nil {
		return fmt.Errorf("prepare get country_ref id: %w", err)
	}
	defer getCountryIDStmt.Close()

	getCertIDStmt, err := newDB.PrepareContext(ctx, `
		SELECT id FROM certificate_ref WHERE name = $1
	`)
	if err != nil {
		return fmt.Errorf("prepare get certificate_ref id: %w", err)
	}
	defer getCertIDStmt.Close()

	insertStmt, err := newDB.PrepareContext(ctx, `
		INSERT INTO certificate_country (country_id, certificate_id, min_age)
		VALUES ($1, $2, $3)
		ON CONFLICT (country_id, certificate_id) DO UPDATE
		SET min_age = EXCLUDED.min_age
	`)
	if err != nil {
		return fmt.Errorf("prepare insert certificate_country: %w", err)
	}
	defer insertStmt.Close()

	var (
		skipped int
		count   int
	)

	for rows.Next() {
		var (
			oldCountryID     int
			oldCertID        int
			age              sql.NullInt64
			countryName      string
			certificateName  string
		)
		if err := rows.Scan(&oldCountryID, &oldCertID, &age, &countryName, &certificateName); err != nil {
			return fmt.Errorf("scan CertificateCountryRef: %w", err)
		}

		var newCountryID int
		if err := getCountryIDStmt.QueryRowContext(ctx, countryName).Scan(&newCountryID); err != nil {
			if err == sql.ErrNoRows {
				log.Printf("migrateCertificateCountry: WARNING: no new country_ref for %q (old CountryID=%d), skipping", countryName, oldCountryID)
				skipped++
				continue
			}
			return fmt.Errorf("lookup country_ref for %q: %w", countryName, err)
		}

		var newCertID int
		if err := getCertIDStmt.QueryRowContext(ctx, certificateName).Scan(&newCertID); err != nil {
			if err == sql.ErrNoRows {
				log.Printf("migrateCertificateCountry: WARNING: no new certificate_ref for %q (old CertificateID=%d), skipping", certificateName, oldCertID)
				skipped++
				continue
			}
			return fmt.Errorf("lookup certificate_ref for %q: %w", certificateName, err)
		}

		var minAge *int64
		if age.Valid {
			minAge = &age.Int64
		}

		_, err := insertStmt.ExecContext(ctx, newCountryID, newCertID, minAge)
		if err != nil {
			return fmt.Errorf("insert certificate_country: %w", err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate CertificateCountryRef: %w", err)
	}

	log.Printf("migrateCertificateCountry: inserted/updated %d rows (skipped %d)", count, skipped)
	return nil
}
