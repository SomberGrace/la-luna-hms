package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type authContextKey struct{}

// AuthenticatedUser contains the verified JWT claims used for authorization checks.
type AuthenticatedUser struct {
	UserID   string
	Username string
	Role     string
}

func parseAuthClaims(secret string, r *http.Request) (*AuthenticatedUser, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("missing authorization")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userID, _ := claims["user_id"].(string)
	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)
	if userID == "" || username == "" || role == "" {
		return nil, fmt.Errorf("missing required claims")
	}

	return &AuthenticatedUser{UserID: userID, Username: username, Role: strings.ToUpper(role)}, nil
}

func withAuthenticatedUser(r *http.Request, user *AuthenticatedUser) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), authContextKey{}, user))
}

func UserFromRequest(r *http.Request) *AuthenticatedUser {
	user, _ := r.Context().Value(authContextKey{}).(*AuthenticatedUser)
	return user
}

func RequireRole(secret string, allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[strings.ToUpper(role)] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next(w, r)
				return
			}

			user, err := parseAuthClaims(secret, r)
			if err != nil {
				http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
				return
			}

			if _, ok := allowed[user.Role]; !ok {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}

			next(w, withAuthenticatedUser(r, user))
		}
	}
}
