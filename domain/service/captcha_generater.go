package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/component/sms"
	"github.com/ihezebin/olympus/random"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type CaptchaGenerater interface {
	Generate(ctx context.Context, key string) error
	Verify(ctx context.Context, key, captcha string) (bool, error)
}

type smsCaptchaGenerater struct {
	redisCli redis.UniversalClient
	smsCli   *sms.SmsClient
	timeout  time.Duration
}

func NewSmsCaptchaGenerater(timeout time.Duration) CaptchaGenerater {
	return &smsCaptchaGenerater{
		redisCli: cache.RedisCacheClient(),
		smsCli:   sms.Client(),
		timeout:  timeout,
	}
}

func (s *smsCaptchaGenerater) key(key string) string {
	return fmt.Sprintf("captcha:sms:%s", key)
}

func (s *smsCaptchaGenerater) Generate(ctx context.Context, key string) error {
	captcha := random.DigitString(6)
	redisKey := s.key(key)
	phone := key
	// 校验手机号是否合法
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number")
	}

	err := s.smsCli.SendCapatchMessage(ctx, key, captcha, strconv.Itoa(int(s.timeout.Minutes())))
	if err != nil {
		return errors.Wrap(err, "send sms captcha failed")
	}

	err = s.redisCli.Set(ctx, redisKey, captcha, s.timeout).Err()
	if err != nil {
		return errors.Wrap(err, "set captcha to redis failed")
	}
	return nil
}

func (s *smsCaptchaGenerater) Verify(ctx context.Context, key, captcha string) (bool, error) {
	redisKey := s.key(key)
	redisCaptcha, err := s.redisCli.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, errors.Wrap(err, "get captcha from redis failed")
	}

	if redisCaptcha != captcha {
		return false, nil
	}
	return true, nil
}
