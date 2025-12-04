package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	sqlPath := flag.String("sql", "", "Path to input schema.sql")
	outPath := flag.String("out", "", "Path to output .mmd ER diagram")
	flag.Parse()

	if *sqlPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: gen_erd -sql path/to/schema.sql -out path/to/schema.mmd")
		os.Exit(1)
	}

	absSQL, err := filepath.Abs(*sqlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving sql path: %v\n", err)
		os.Exit(1)
	}

	absOut, err := filepath.Abs(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving out path: %v\n", err)
		os.Exit(1)
	}

	// Route via "old" or "new" helper, based on path.
	// Both use the same core generator, but we keep separate entrypoints.
	var genErr error
	lower := strings.ToLower(absSQL)

	if strings.Contains(lower, string(os.PathSeparator)+"old"+string(os.PathSeparator)) {
		genErr = GenerateOldERD(absSQL, absOut)
	} else if strings.Contains(lower, string(os.PathSeparator)+"new"+string(os.PathSeparator)) {
		genErr = GenerateNewERD(absSQL, absOut)
	} else {
		// Fallback: generic
		genErr = GenerateERD(absSQL, absOut)
	}

	if genErr != nil {
		fmt.Fprintf(os.Stderr, "error generating ERD: %v\n", genErr)
		os.Exit(1)
	}

	fmt.Printf("ERD generated: %s\n", absOut)
}
