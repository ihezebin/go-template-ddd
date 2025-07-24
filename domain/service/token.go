package service

import (
	"context"
	"time"

	"github.com/ihezebin/jwt"
	"github.com/pkg/errors"
)

type TokenService struct {
	secret  string
	expired time.Duration
}

func NewTokenService(secret string, expired time.Duration) *TokenService {
	return &TokenService{secret: secret, expired: expired}
}

func (s *TokenService) Signed(ctx context.Context, owner string) (string, error) {
	token := jwt.Default(jwt.WithOwner(owner), jwt.WithExpire(s.expired))
	tokenStr, err := token.Signed(s.secret)
	if err != nil {
		return "", errors.Wrap(err, "generate token err")
	}

	// 保存 token 到 redis
	// tokenKey := fmt.Sprintf(constant.TokenRedisKeyFormat, owner)
	// err = cache.RedisCacheClient().Set(ctx, tokenKey, tokenStr, s.expired).Err()
	// if err != nil {
	// 	return "", errors.Wrapf(err, "save token to redis err, key: %s, value: %s", tokenKey, tokenStr)
	// }

	return tokenStr, nil
}
