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

func (s *TokenService) Verify(ctx context.Context, tokenStr string) (payload *jwt.Payload, faked bool, expired bool, err error) {
	token, err := jwt.Parse(tokenStr, s.secret)
	if err != nil {
		return nil, false, false, errors.Wrap(err, "parse token err")
	}

	payload = token.Payload()

	faked, err = token.Faked()
	if err != nil {
		return payload, false, false, errors.Wrap(err, "faked token err")
	}

	if faked {
		return payload, true, false, nil
	}

	expired = token.Expired()
	if expired {
		return payload, false, true, nil
	}

	return payload, false, false, nil
}
