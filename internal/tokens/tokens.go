package tokens

import (
	"fmt"
	"simbirGo/internal/entities"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateNewJwt(user entities.User) (string, error) {
	op := "usecase.token.GenerateNewJwt()"
	key := []byte("boba")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":      user.Id,
			"isAdmin": user.IsAdmin,
		})

	strToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign jwt: %w", op, err)
	}
	return strToken, nil
}

func ParseToken(tokenString string) (entities.Token, error) {
	op := "tokens.ParseToken()"
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method: %v", op, t.Header["alg"])
		}
		return []byte("boba"), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return entities.Token{
			Id:      uint(claims["id"].(float64)),
			IsAdmin: claims["isAdmin"].(bool),
		}, nil
	}

	return entities.Token{}, fmt.Errorf("%s: failed to parse token: %w", op, err)
}
