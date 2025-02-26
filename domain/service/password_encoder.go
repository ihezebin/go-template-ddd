package service

import (
	"crypto/md5"
	"encoding/hex"
	"unicode"

	"github.com/pkg/errors"
)

type PasswordEncoder interface {
	Encode(plainPwd string, salt string) (string, error)
	Verify(plainPwd, salt, encodedPwd string) (bool, error)
	Strength(plainPwd string) int
}

type md5WithSaltPasswordEncoder struct {
}

func NewMd5WithSaltPasswordEncoder() PasswordEncoder {
	return &md5WithSaltPasswordEncoder{}
}

func (e *md5WithSaltPasswordEncoder) Strength(plainPwd string) int {
	var score int
	password := plainPwd
	// 长度评分
	if len(password) >= 8 {
		score += 1
	}
	// 字符分类统计
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasLower {
		score += 1
	}
	if hasUpper {
		score += 1
	}
	if hasDigit {
		score += 1
	}
	if hasSpecial {
		score += 1
	}

	return score
}

func (e *md5WithSaltPasswordEncoder) Encode(plainPwd string, salt string) (string, error) {
	m5 := md5.New()
	_, err := m5.Write([]byte(plainPwd))
	if err != nil {
		return "", errors.Wrapf(err, "md5 write plain pwd error")
	}
	_, err = m5.Write([]byte(salt))
	if err != nil {
		return "", errors.Wrapf(err, "md5 write salt error")
	}
	return hex.EncodeToString(m5.Sum(nil)), nil
}

func (e *md5WithSaltPasswordEncoder) Verify(plainPwd, encodedPwd string, salt string) (bool, error) {
	encodedPwd, err := e.Encode(plainPwd, salt)
	if err != nil {
		return false, errors.Wrapf(err, "encode plain pwd error")
	}
	return encodedPwd == encodedPwd, nil
}
