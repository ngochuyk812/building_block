package auth_context

import (
	"log"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	req := &ClaimModel{
		IdSite: "231",
	}

	token, err := GenerateJWT(req, "key", 2*time.Minute)
	if err != nil {
		t.Error(err.Error())
	}
	ca, err := VerifyJWT(token, "key")
	if err != nil {
		t.Error(err.Error())
	}
	log.Printf("%+v", ca)
	if ca.IdSite != req.IdSite {
		t.Error("err verify token ")

	}
}
