package entity

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/ihezebin/web-template-ddd/component/constant"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Test struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name"`
	Password string             `json:"password"`
}

// 实体应该具备一些自己的特殊能力

// MD5Password 密码加密
func (t *Test) MD5Password() error {
	m5 := md5.New()
	_, err := m5.Write([]byte(t.Password))
	if err != nil {
		return err
	}
	_, err = m5.Write([]byte(constant.MD5_Salt))
	if err != nil {
		return err
	}

	t.Password = hex.EncodeToString(m5.Sum(nil))

	return nil
}

// String 序列化为字符串
func (t *Test) String() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
