package refreshtoken

import (
	"context"
	"fmt"
	"time"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/redis"
	rdb "github.com/redis/go-redis/v9"
)

type store struct {
	*redis.Redis
}

func New(r *redis.Redis) core.RefreshTokenStore {
	return &store{r}
}

func (s *store) GetRefreshToken(ctx context.Context, tokenID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userID, err := s.Redis.Get(ctx, tokenID).Int()
	if err == rdb.Nil {
		logger.Log().Debug(ctx, "refresh token does not exists")
		return 0, core.ErrRefreshTokenNotValid
	} else if err != nil {
		logger.Log().Error(ctx, "failed to get refresh token: %w", err)
		return 0, err
	}

	return userID, nil
}

func (s *store) DeleteRefreshToken(ctx context.Context, prevTokenID string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	fmt.Println(prevTokenID)

	deleted, err := s.Redis.Del(ctx, prevTokenID).Result()
	if deleted == 0 {
		logger.Log().Warn(ctx, "no keys were deleted. Key may not exist: %s", prevTokenID)
	} else if err != nil {
		logger.Log().Error(ctx, "failed to delete refresh token: %w", err)
		return err
	}

	return nil
}

func (s *store) SetRefreshToken(ctx context.Context, userID int, tokenID string, expiresIn time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := s.Redis.Set(ctx, tokenID, userID, expiresIn).Err(); err != nil {
		logger.Log().Error(ctx, "failed to set refresh token: %w", err)
		return err
	}

	return nil
}
