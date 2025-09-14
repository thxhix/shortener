package meta

import (
	"fmt"
)

func PrintMeta(ver string, date string, commit string) {
	fmt.Printf("Build version: %s\n", prepareValue(ver))
	fmt.Printf("Build date: %s\n", prepareValue(date))
	fmt.Printf("Build commit: %s\n", prepareValue(commit))
}

func prepareValue(v string) string {
	if v == "" {
		return "N/A"
	}
	return v
}
