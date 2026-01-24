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

func CreateJWT(secret []byte, user *types.User) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		string(UserKey): fmt.Sprintf("%d", user.ID),
		"expiredAt":     time.Now().Add(expiration).Unix(),
		"admin":         user.IsAdmin,
	})

	return token.SignedString(secret)
}

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("error extracting claims from jwt token")
			permissionDenied(w)
			return
		}

		userIdVal, ok := claims[string(UserKey)]
		if !ok {
			log.Println("userID claim not found")
			permissionDenied(w)
			return
		}

		var userId string
		switch v := userIdVal.(type) {
		case string:
			userId = v
		case float64:
			userId = fmt.Sprintf("%.0f", v)
		default:
			log.Printf("unexpected type for userID claim: %T\n", v)
			permissionDenied(w)
			return
		}

		isAdminVal, ok := claims["admin"]
		var isAdmin bool = false
		if ok {
			if b, ok := isAdminVal.(bool); ok {
				isAdmin = b
			}
		}

		// u, err := store.GetUserById(userId)
		// if err != nil {
		// 	log.Printf("failed to get user by id: %v\n", err)
		// 	permissionDenied(w)
		// 	return
		// }

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userId)
		ctx = context.WithValue(ctx, "admin", isAdmin)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func WithAdminAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		user := GetUserFromContext(r.Context())
		if user == nil {
			log.Println("jwtUser not present in context")
			permissionDenied(w)
			return
		}

		if !user.IsAdmin {
			log.Println("unauthorized request")
			permissionDenied(w)
			return
		}

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

func GetUserFromContext(ctx context.Context) *types.JWTUser {
	userId, ok := ctx.Value(UserKey).(string)
	if !ok {
		return nil
	}
	isAdmin, ok := ctx.Value("admin").(bool)
	if !ok {
		isAdmin = false
	}


	return &types.JWTUser{
		Id: userId,
		IsAdmin: isAdmin,
	}
}
