package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Column represents a column in a table.
type Column struct {
	Name string
	Type string
	IsPK bool
}

// Table represents a table with schema/name and columns.
type Table struct {
	Schema  string
	Name    string
	Columns []Column
}

// Relationship represents a foreign-key style relationship between tables.
type Relationship struct {
	Parent string // "schema.table" that is referenced
	Child  string // "schema.table" that has the FK
	Label  string // e.g. "FK"
}

// GenerateNewERD is the entry point for the "new" schema ERD.
func GenerateNewERD(sqlPath, outPath string) error {
	return GenerateERD(sqlPath, outPath)
}

// GenerateERD is the common generator used by both old and new paths.
func GenerateERD(sqlPath, outPath string) error {
	data, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("reading sql file: %w", err)
	}

	tables, rels := parseSQLSchema(string(data))
	mermaid := buildMermaidERD(tables, rels)

	if err := os.WriteFile(outPath, []byte(mermaid), 0o644); err != nil {
		return fmt.Errorf("writing mermaid file: %w", err)
	}

	return nil
}

// ----------------------------
// SQL parsing
// ----------------------------

// parseSQLSchema parses a PostgreSQL schema SQL into tables and relationships.
// It is designed to handle both:
//
//   - pg_dump-style old schema.sql (quoted identifiers, multi-line
//     ALTER TABLE ... FOREIGN KEY ...)
//   - hand-written new schema.sql (simple CREATE TABLE with inline REFERENCES)
//
// It only looks at:
//
//   - CREATE TABLE ... ( ... )
//   - CONSTRAINT ... PRIMARY KEY (...)
//   - inline "REFERENCES other_schema.other_table(...)"
//   - ALTER TABLE ... ADD CONSTRAINT ... FOREIGN KEY (...) REFERENCES ...
func parseSQLSchema(sql string) (map[string]*Table, []Relationship) {
	tables := make(map[string]*Table)
	var rels []Relationship

	// Regexes with (?is) => case-insensitive, dot matches newline.
	createTableRe := regexp.MustCompile(`(?is)^CREATE\s+TABLE\s+(.+?)\s*\((.*)\)\s*$`)
	alterFKRe := regexp.MustCompile(
		`(?is)^ALTER\s+TABLE\s+(?:ONLY\s+)?(.+?)\s+ADD\s+CONSTRAINT\s+.+?FOREIGN\s+KEY\s*\(([^)]+)\)\s+REFERENCES\s+(.+?)\s*\(([^)]+)\)`)

	statements := strings.Split(sql, ";")
	for _, raw := range statements {
		sFull := strings.TrimSpace(raw)
		if sFull == "" {
			continue
		}
		if isCommentOnly(sFull) {
			continue
		}

		upper := strings.ToUpper(sFull)

		// ---- CREATE TABLE ----
		if idx := strings.Index(upper, "CREATE TABLE"); idx != -1 {
			s := strings.TrimSpace(sFull[idx:])
			if m := createTableRe.FindStringSubmatch(s); m != nil {
				tblNameRaw := strings.TrimSpace(m[1])
				body := m[2]

				schema, name := splitQualified(tblNameRaw)
				key := fmt.Sprintf("%s.%s", schema, name)

				table := &Table{
					Schema:  schema,
					Name:    name,
					Columns: nil,
				}

				parts := splitCommaTopLevel(body)
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if part == "" {
						continue
					}
					upPart := strings.ToUpper(part)

					// Table-level CONSTRAINT (PK/FK)
					if strings.HasPrefix(upPart, "CONSTRAINT") {
						// PRIMARY KEY (col1, col2, ...)
						if strings.Contains(upPart, "PRIMARY KEY") {
							if colsText := firstParenGroup(part); colsText != "" {
								colNames := strings.Split(colsText, ",")
								for _, cn := range colNames {
									nameTrim := strings.Trim(strings.TrimSpace(cn), `"`)
									for i, c := range table.Columns {
										if c.Name == nameTrim {
											table.Columns[i].IsPK = true
										}
									}
								}
							}
						} else if strings.Contains(upPart, "FOREIGN KEY") {
							// Table-level FK: CONSTRAINT ... FOREIGN KEY (...) REFERENCES other_schema.other_table(...)
							if tgtName := extractReferencedTable(part); tgtName != "" {
								tgtSchema, tgtTable := splitQualified(tgtName)
								parent := fmt.Sprintf("%s.%s", tgtSchema, tgtTable)
								child := fmt.Sprintf("%s.%s", schema, name)
								rels = append(rels, Relationship{
									Parent: parent,
									Child:  child,
									Label:  "FK",
								})
							}
						}
						continue
					}

					// Column definition
					colName, colType, ok := parseColumnDef(part)
					if !ok {
						continue
					}
					isPK := strings.Contains(upPart, "PRIMARY KEY")

					table.Columns = append(table.Columns, Column{
						Name: colName,
						Type: colType,
						IsPK: isPK,
					})

					// Inline FK: col ... REFERENCES other_schema.other_table(...)
					if strings.Contains(upPart, "REFERENCES") {
						if tgtName := extractReferencedTable(part); tgtName != "" {
							tgtSchema, tgtTable := splitQualified(tgtName)
							parent := fmt.Sprintf("%s.%s", tgtSchema, tgtTable)
							child := fmt.Sprintf("%s.%s", schema, name)
							rels = append(rels, Relationship{
								Parent: parent,
								Child:  child,
								Label:  "FK",
							})
						}
					}
				}

				tables[key] = table
				continue
			}
		}

		// ---- ALTER TABLE ... ADD CONSTRAINT ... FOREIGN KEY ----
		if idx := strings.Index(upper, "ALTER TABLE"); idx != -1 {
			s := strings.TrimSpace(sFull[idx:])
			if m := alterFKRe.FindStringSubmatch(s); m != nil {
				srcNameRaw := strings.TrimSpace(m[1])
				tgtNameRaw := strings.TrimSpace(m[3])

				srcSchema, srcTable := splitQualified(srcNameRaw)
				tgtSchema, tgtTable := splitQualified(tgtNameRaw)

				parent := fmt.Sprintf("%s.%s", tgtSchema, tgtTable)
				child := fmt.Sprintf("%s.%s", srcSchema, srcTable)

				rels = append(rels, Relationship{
					Parent: parent,
					Child:  child,
					Label:  "FK",
				})
				continue
			}
		}
	}

	return tables, rels
}

// isCommentOnly reports whether a statement is only "--" comments and blank lines.
func isCommentOnly(s string) bool {
	for _, line := range strings.Split(s, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(trimmed, "--") {
			return false
		}
	}
	return true
}

// splitQualified splits a possibly schema-qualified identifier into schema and table.
//
// Examples:
//   "public.CastTable"             -> "public", "CastTable"
//   "\"Lines\".\"AwardTitleLine\"" -> "Lines", "AwardTitleLine"
//   "title"                        -> "public", "title"
func splitQualified(name string) (schema, table string) {
	name = strings.TrimSpace(name)

	if strings.Contains(name, ".") {
		parts := strings.SplitN(name, ".", 2)
		schema = strings.Trim(parts[0], `" `)
		table = strings.Trim(parts[1], `" `)
		if schema == "" {
			schema = "public"
		}
		if table == "" {
			table = "X"
		}
		return schema, table
	}

	// No explicit schema -> assume public
	return "public", strings.Trim(name, `" `)
}

// splitCommaTopLevel splits a CREATE TABLE body by top-level commas (not inside parentheses).
func splitCommaTopLevel(body string) []string {
	var parts []string
	var buf strings.Builder
	depth := 0

	for _, r := range body {
		switch r {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				part := strings.TrimSpace(buf.String())
				if part != "" {
					parts = append(parts, part)
				}
				buf.Reset()
				continue
			}
		}
		buf.WriteRune(r)
	}

	if tail := strings.TrimSpace(buf.String()); tail != "" {
		parts = append(parts, tail)
	}

	return parts
}

// firstParenGroup returns the first (...) group content in s, or "" if none.
func firstParenGroup(s string) string {
	start := strings.Index(s, "(")
	if start == -1 {
		return ""
	}
	end := strings.Index(s[start+1:], ")")
	if end == -1 {
		return ""
	}
	return s[start+1 : start+1+end]
}

// extractReferencedTable finds the identifier after "REFERENCES".
func extractReferencedTable(s string) string {
	up := strings.ToUpper(s)
	idx := strings.Index(up, "REFERENCES")
	if idx == -1 {
		return ""
	}
	rest := strings.TrimSpace(s[idx+len("REFERENCES"):])
	// Up to first '('
	if p := strings.Index(rest, "("); p != -1 {
		rest = rest[:p]
	}
	return strings.TrimSpace(rest)
}

// parseColumnDef extracts column name and type from a column definition line.
func parseColumnDef(part string) (colName, colType string, ok bool) {
	part = strings.TrimSpace(part)
	if part == "" {
		return "", "", false
	}

	// Remove trailing comma if present (should already be handled, but safe).
	if strings.HasSuffix(part, ",") {
		part = strings.TrimSpace(strings.TrimSuffix(part, ","))
	}

	// We only care about: name TYPE ...
	// name can be quoted or unquoted.
	// TYPE is taken as the first token after the name (e.g. "integer", "varchar").
	re := regexp.MustCompile(`^"?(?P<name>[A-Za-z_][A-Za-z0-9_]*)"?\s+(?P<type>[A-Za-z0-9_]+)`)
	m := re.FindStringSubmatch(part)
	if m == nil {
		return "", "", false
	}

	colName = m[1]
	colType = strings.ToUpper(m[2])
	return colName, colType, true
}

// ----------------------------
// Mermaid ERD generation
// ----------------------------

// buildMermaidERD builds a Mermaid erDiagram string from tables and relationships.
func buildMermaidERD(tables map[string]*Table, rels []Relationship) string {
	var lines []string
	lines = append(lines, "erDiagram")

	// Sort tables by schema.name for stable output
	keys := make([]string, 0, len(tables))
	for k := range tables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		tbl := tables[key]
		entityName := mermaidEntityName(tbl.Schema, tbl.Name)

		lines = append(lines, fmt.Sprintf("  %s {", entityName))
		if len(tbl.Columns) == 0 {
			lines = append(lines, "    STRING id")
		} else {
			for _, col := range tbl.Columns {
				colType := col.Type
				if colType == "" {
					colType = "STRING"
				}
				colNameSafe := mermaidSafe(col.Name)
				pkSuffix := ""
				if col.IsPK {
					pkSuffix = " PK"
				}
				lines = append(lines, fmt.Sprintf("    %s %s%s", colType, colNameSafe, pkSuffix))
			}
		}
		lines = append(lines, "  }")
		lines = append(lines, "")
	}

	// De-duplicate relationships
	seen := make(map[string]struct{})
	for _, r := range rels {
		parentEntity := mermaidEntityNameFromQualified(r.Parent)
		childEntity := mermaidEntityNameFromQualified(r.Child)
		label := r.Label
		if label == "" {
			label = "FK"
		}
		key := parentEntity + "->" + childEntity + ":" + label
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		// Parent ||--o{ Child : FK
		lines = append(lines, fmt.Sprintf("  %s ||--o{ %s : %s", parentEntity, childEntity, label))
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

func mermaidEntityName(schema, table string) string {
	return mermaidSafe(schema + "_" + table)
}

func mermaidEntityNameFromQualified(qualified string) string {
	schema, table := splitQualified(qualified)
	return mermaidEntityName(schema, table)
}

// mermaidSafe converts an identifier into something acceptable for Mermaid ER diagrams.
// It:
//   - replaces non-alphanumeric / non-underscore characters with '_'
//   - ensures the first character is a letter or '_'
func mermaidSafe(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "X"
	}

	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}

	safe := b.String()
	if safe == "" {
		return "X"
	}
	// Mermaid identifiers must start with letter or underscore
	first := safe[0]
	if !((first >= 'a' && first <= 'z') ||
		(first >= 'A' && first <= 'Z') ||
		first == '_') {
		safe = "T_" + safe
	}
	return safe
}
