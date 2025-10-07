package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"uid"`
	jwt.RegisteredClaims
}
// Auto-generated swagger comments for jwtSecret
// @Summary Auto-generated summary for jwtSecret
// @Description Auto-generated description for jwtSecret — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func jwtSecret() []byte {
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "dev-secret-change-me"
	}
	return []byte(sec)
}

func GenerateToken(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret())
}
// Auto-generated swagger comments for parseToken
// @Summary Auto-generated summary for parseToken
// @Description Auto-generated description for parseToken — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) { return jwtSecret(), nil })
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
// Auto-generated swagger comments for JWTMiddleware
// @Summary Auto-generated summary for JWTMiddleware
// @Description Auto-generated description for JWTMiddleware — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if authz == "" || !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tok := strings.TrimSpace(authz[len("Bearer "):])
		claims, err := parseToken(tok)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
// Auto-generated swagger comments for UpgradeWithJWT
// @Summary Auto-generated summary for UpgradeWithJWT
// @Description Auto-generated description for UpgradeWithJWT — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func UpgradeWithJWT(fn func(*gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			JWTMiddleware()(c)
			if c.IsAborted() {
				return
			}
		} else {
			tok := c.Query("token")
			if tok == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
				return
			}
			claims, err := parseToken(tok)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}
			c.Set("userID", claims.UserID)
		}
		fn(c)
	}
}
