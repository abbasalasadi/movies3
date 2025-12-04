package main

// GenerateOldERD generates an ERD from the "old" LabVIEW-era schema.
// Right now it just delegates to the generic generator, but keeping this
// function separate lets us customize behavior later if we need to.
func GenerateOldERD(sqlPath, outPath string) error {
	return GenerateERD(sqlPath, outPath)
}
