package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pahulgogna/evoAI_Web/backend/config"
	"github.com/pahulgogna/evoAI_Web/backend/types"
	"github.com/pahulgogna/evoAI_Web/backend/utils"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		string(UserKey) : userID,
		"expiredAt" : time.Now().Add(expiration).Unix(),
	})

	return token.SignedString(secret)
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)

		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v\n", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims[string(UserKey)].(string)

		u, err := store.GetUserById(userId)
		if err != nil {
			log.Printf("failed to get user by id: %v\n", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx	= context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}


func getTokenFromRequest(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func validateToken(authToken string) (*jwt.Token, error) {
	return jwt.Parse(authToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}


func permissionDenied(w http.ResponseWriter) {
	utils.WriteErrorResponse(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}


func GetUserIDFromContext(ctx context.Context) string {
	userId, ok := ctx.Value(UserKey).(string)
	if !ok {
		return ""
	}

	return userId
}