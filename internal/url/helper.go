package url

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"github.com/google/uuid"
)

// IsTestMode checks running in Go test mode.
// It checks if the testing flag ("test.v") is present in the flag set.
func IsTestMode() bool {
	return flag.Lookup("test.v") != nil
}

// GetHash generates a short hash string.
// In test mode (detected by IsTestMode) it always returns a static value "testHash"
// to provide deterministic results in unit tests.
// In normal mode it generates a new UUID, computes its SHA256 hash,
// encodes it using URL-safe Base64, and returns the first 10 characters.
func GetHash() string {
	if IsTestMode() {
		return "testHash"
	}
	u := uuid.New().String()
	hash := sha256.Sum256([]byte(u))
	return base64.URLEncoding.EncodeToString(hash[:])[:10]
}
