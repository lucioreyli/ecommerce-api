package auth

import (
	"store-api/config"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte(config.Envs.JWTSecret)
	t.Run("Should create JWT", func(t *testing.T) {
		token, err := CreateJWT(secret, 1)
		if err != nil {
			t.Errorf("Error creating JWT: %V", err)
		}

		if token == "" {
			t.Error("Expected token to be not empty")
		}
	})
}
