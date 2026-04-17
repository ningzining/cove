package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerate(t *testing.T) {
	claims := &CustomMapClaims{
		Provider: "test",
		UserID:   "019d95b9add77b51bd74343017291392",
		Phone:    "13800000001",
		Nickname: "普通用户",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "cove",
			Subject:   "user",
			Audience:  []string{"cove"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	tokenString, err := Generate(claims, "cove-secret-key-change-in-production")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tokenString: %s", tokenString)
}
