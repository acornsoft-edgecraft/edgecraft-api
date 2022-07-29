// Package common - Defines the functions for Utiltity
package common

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"strconv"
	"time"
)

// ===== [ Constants and Variables ] =====

const (
	// ContextJWTKey - The key for the jwt context value
	ContextJWTKey ContextKey = "jwt"

	// Random string generation
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var (
	// Random source by nano time
	randSrc = rand.NewSource(time.Now().UnixNano())
)

// ===== [ Types ] =====

// ContextKey - Implements type for context key
type ContextKey string

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// CreateRandString - Creates an random string with the size of n
func CreateRandString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// ParseJWT - Parses and validates a token using the HMAC signing method
func ParseJWT(secretKey string, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("jwt validation failed")
}

// CreateJWT - Creates, signs, and encodes a JWT token using the HMAC signing method
func CreateJWT(secretKey string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetID - Returns the id from claims of jwt token
func GetID(claims jwt.MapClaims) (string, error) {
	id, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf(strconv.Itoa(CodeInvalidToken))
	}
	return id, nil
}
