package usecase

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"github.com/google/uuid"
)

func IsTestMode() bool {
	return flag.Lookup("test.v") != nil
}

func GetHash() string {
	if IsTestMode() {
		return "testHash"
	}
	u := uuid.New().String()
	hash := sha256.Sum256([]byte(u))
	return base64.URLEncoding.EncodeToString(hash[:])[:10]
}
