package meta

import (
	"fmt"
)

// PrintMeta prints build metadata (version, date, commit) to stdout.
// If any of the values are empty, "N/A" will be printed instead.
func PrintMeta(ver string, date string, commit string) {
	fmt.Printf("Build version: %s\n", prepareValue(ver))
	fmt.Printf("Build date: %s\n", prepareValue(date))
	fmt.Printf("Build commit: %s\n", prepareValue(commit))
}

// prepareValue returns v, or "N/A" if v is an empty string.
func prepareValue(v string) string {
	if v == "" {
		return "N/A"
	}
	return v
}
