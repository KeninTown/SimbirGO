package tokens

import (
	"fmt"
	"simbirGo/internal/entities"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateNewJwt(user entities.User) (string, error) {
	op := "usecase.token.GenerateNewJwt()"
	key := []byte("byte")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":      user.Id,
			"isAdmin": user.IsAdmin,
		})

	sToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign jwt: %w", op, err)
	}
	return sToken, nil
}
