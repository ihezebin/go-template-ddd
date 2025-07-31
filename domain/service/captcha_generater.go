package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/ihezebin/go-template-ddd/component/constant"
	"github.com/ihezebin/go-template-ddd/component/sms"
	"github.com/ihezebin/olympus/logger"
	"github.com/ihezebin/olympus/random"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type CaptchaGenerater interface {
	Generate(ctx context.Context, key string, usage constant.CaptchaUsageType) (bool, string, int, error)
	Verify(ctx context.Context, key string, usage constant.CaptchaUsageType, captcha string) (bool, error)
}

type smsCaptchaGenerater struct {
	redisCli       redis.UniversalClient
	smsCli         sms.SmsClient
	timeout        time.Duration
	frequencyLimit time.Duration
}

var globalCaptchaGenerater CaptchaGenerater

func GetCaptchaGenerater() CaptchaGenerater {
	return globalCaptchaGenerater
}

func SetCaptchaGenerater(generater CaptchaGenerater) {
	globalCaptchaGenerater = generater
}

func NewSmsCaptchaGenerater(redisCli redis.UniversalClient, smsCli sms.SmsClient, timeout time.Duration, frequencyLimit time.Duration) CaptchaGenerater {
	return &smsCaptchaGenerater{
		redisCli:       redisCli,
		smsCli:         smsCli,
		timeout:        timeout,
		frequencyLimit: frequencyLimit,
	}
}

func (s *smsCaptchaGenerater) key(key string, usage constant.CaptchaUsageType) string {
	return fmt.Sprintf("captcha:sms:%s:%s", key, usage.String())
}

func (s *smsCaptchaGenerater) Generate(ctx context.Context, key string, usage constant.CaptchaUsageType) (bool, string, int, error) {
	redisKey := s.key(key, usage)
	phone := key
	// 校验手机号是否合法
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !re.MatchString(phone) {
		return false, "", 0, errors.New("invalid phone number")
	}

	// 校验频率是否ok
	freqOk, freq, err := s.isFrequencyOk(ctx, key)
	if err != nil {
		return false, "", 0, errors.Wrap(err, "check frequency failed")
	}
	if !freqOk {
		return false, "", freq, nil
	}

	captcha := random.DigitString(6)
	go func() {
		nCtx := context.Background()
		err = s.smsCli.SendCapatchMessage(nCtx, key, captcha)
		if err != nil {
			logger.Errorf(nCtx, "send sms captcha failed, phone: %s, captcha: %s, err: %v", key, captcha, err)
		}
	}()
	err = s.redisCli.SetNX(ctx, redisKey, captcha, s.timeout).Err()
	if err != nil {
		return false, captcha, 0, errors.Wrap(err, "set captcha to redis failed")
	}
	return true, captcha, 0, nil
}

func (s *smsCaptchaGenerater) Verify(ctx context.Context, key string, usage constant.CaptchaUsageType, captcha string) (bool, error) {
	redisKey := s.key(key, usage)
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

func (s *smsCaptchaGenerater) isFrequencyOk(ctx context.Context, key string) (bool, int, error) {
	// 频率限制
	freqKey := s.frequencyKey(key)
	ok, err := s.redisCli.SetNX(ctx, freqKey, 0, s.frequencyLimit).Result()
	if err != nil {
		return false, 0, errors.Wrap(err, "set frequency to redis failed")
	}
	if !ok { // 频率限制未通过，则记录限制下的触发次数并增加限制时间
		freq, err := s.redisCli.Get(ctx, freqKey).Int()
		if err != nil {
			return false, 0, errors.Wrap(err, "get frequency from redis failed")
		}
		err = s.redisCli.Expire(ctx, freqKey, s.frequencyLimit+time.Duration(freq*int(s.timeout))).Err()
		if err != nil {
			return false, 0, errors.Wrap(err, "expire frequency failed")
		}
		err = s.redisCli.Incr(ctx, freqKey).Err()
		if err != nil {
			return false, 0, errors.Wrap(err, "incr frequency failed")
		}
		return false, freq, nil
	}
	return true, 0, nil
}

func (s *smsCaptchaGenerater) delFrequency(ctx context.Context, key string) error {
	return s.redisCli.Del(ctx, s.frequencyKey(key)).Err()
}
