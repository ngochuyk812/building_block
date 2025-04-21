package auth_context

import (
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	_, err := GenerateJWT(&ClaimModel{
		IdSite: "231",
	}, "key", 2*time.Minute)
	if err != nil {
		t.Error(err.Error())
	}
}
