package web

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

var signingKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// generateJwt is a mock function to generate a JWT token.
func generateJwt(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// JwtValidator is middleware function to validate JWT.
func JwtValidator(ctx context.Context, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("JwtValidator called")
		// 从请求头中获取 Authorization 字段
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 检查 Authorization 字段格式是否为 "Bearer <token>"
		const prefix = "Bearer "
		if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// 提取 JWT
		tokenStr := authHeader[len(prefix):]
		var username string

		// 解析和验证 JWT
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		}, jwt.WithLeeway(5+time.Second))

		if err != nil || !token.Valid {
			errMsg := "Invalid JWT token, " + err.Error()
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		} else if claims, ok := token.Claims.(*Claims); ok {
			username = claims.Username
			if username == "" {
				http.Error(w, "Invalid JWT token, unknown claims", http.StatusUnauthorized)
				return
			}

		} else {
			http.Error(w, "Failed to extract claims from JWT", http.StatusUnauthorized)
			return
		}

		// JWT 验证通过，调用下一个处理函数
		newCtx := context.WithValue(ctx, "username", username)
		next(w, r.WithContext(newCtx))
	}
}
