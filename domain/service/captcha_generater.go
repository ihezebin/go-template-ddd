package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/component/sms"
	"github.com/ihezebin/olympus/random"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type CaptchaGenerater interface {
	Generate(ctx context.Context, key string) (bool, string, error)
	Verify(ctx context.Context, key, captcha string) (bool, error)
}

type smsCaptchaGenerater struct {
	redisCli       redis.UniversalClient
	smsCli         sms.SmsClient
	timeout        time.Duration
	frequencyLimit time.Duration
}

func NewSmsCaptchaGenerater(redisCli redis.UniversalClient, smsCli sms.SmsClient, timeout time.Duration, frequencyLimit time.Duration) CaptchaGenerater {
	return &smsCaptchaGenerater{
		redisCli:       cache.RedisCacheClient(),
		smsCli:         sms.ClientTencent(),
		timeout:        timeout,
		frequencyLimit: frequencyLimit,
	}
}

func (s *smsCaptchaGenerater) key(key string) string {
	return fmt.Sprintf("captcha:sms:%s", key)
}

func (s *smsCaptchaGenerater) Generate(ctx context.Context, key string) (bool, string, error) {
	redisKey := s.key(key)
	phone := key
	// 校验手机号是否合法
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !re.MatchString(phone) {
		return false, "", errors.New("invalid phone number")
	}

	// 校验频率是否ok
	freqOk, err := s.isFrequencyOk(ctx, key)
	if err != nil {
		return false, "", errors.Wrap(err, "check frequency failed")
	}
	if !freqOk {
		return false, "", nil
	}

	captcha := random.DigitString(6)
	err = s.smsCli.SendCapatchMessage(ctx, key, captcha)
	if err != nil {
		return false, captcha, errors.Wrap(err, "send sms captcha failed")
	}

	err = s.redisCli.SetNX(ctx, redisKey, captcha, s.timeout).Err()
	if err != nil {
		return false, captcha, errors.Wrap(err, "set captcha to redis failed")
	}
	return true, captcha, nil
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

	// 验证成功后，删除redis中的验证码
	err = s.redisCli.Del(ctx, redisKey).Err()
	if err != nil {
		return false, errors.Wrap(err, "del captcha from redis failed")
	}
	// 验证成功后，删除频率限制
	err = s.delFrequency(ctx, key)
	if err != nil {
		return false, errors.Wrap(err, "del frequency from redis failed")
	}
	return true, nil
}

func (s *smsCaptchaGenerater) frequencyKey(key string) string {
	return fmt.Sprintf("captcha:sms:frequency:%s", key)
}

func (s *smsCaptchaGenerater) isFrequencyOk(ctx context.Context, key string) (bool, error) {
	// 频率限制
	freqKey := s.frequencyKey(key)
	ok, err := s.redisCli.SetNX(ctx, freqKey, 0, s.frequencyLimit).Result()
	if err != nil {
		return false, errors.Wrap(err, "set frequency to redis failed")
	}
	return ok, nil
}

func (s *smsCaptchaGenerater) delFrequency(ctx context.Context, key string) error {
	return s.redisCli.Del(ctx, s.frequencyKey(key)).Err()
}
