package auth

import (
	"context"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/golang-jwt/jwt"
)

func validToken(ctx context.Context, tokenString string, secret string) (*int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Log().Error(ctx, "unexpected signing method")
			return nil, core.ErrUnauthorized
		}

		return []byte(secret), nil
	})
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, core.ErrUnauthorized
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(float64)
		if !ok {
			return nil, core.ErrUnauthorized
		}

		idInt := int(id)
		return &idInt, nil
	}

	return nil, core.ErrUnauthorized
}

func getUserIDFromContext(ctx context.Context) (int, error) {
	id, ok := ctx.Value(userIDContextKey).(int)
	if !ok {
		logger.Log().Debug(ctx, "user id is not provided")
		return 0, core.ErrUnauthorized
	}

	return id, nil
}
