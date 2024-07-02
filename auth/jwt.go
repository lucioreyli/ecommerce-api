package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"store-api/config"
	"store-api/types"
	"store-api/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTEXpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   strconv.Itoa(userID),
		"expireAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWT(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Failed to validate token %v", err)
			permissionDenied(w)
			return
		}
		if !token.Valid {
			log.Printf("Invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)
		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("Failed to get user %v %d", err, userID)
			permissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("Failed to get user %v", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" {
		return tokenAuth
	}
	return ""
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
}

func GetUserIDFromContext(c context.Context) *int {
	userID, ok := c.Value(UserKey).(int)
	if !ok {
		err := 0
		return &err
	}
	return &userID
}
