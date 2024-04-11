package entity

import (
	"crypto/md5"
	"encoding"
	"encoding/hex"
	"encoding/json"
	"regexp"

	"gorm.io/gorm/schema"
)

type Example struct {
	Id       string `json:"id" bson:"id" gorm:"id"`
	Username string `json:"username" bson:"username" gorm:"column:username"`
	Password string `json:"password" bson:"password" gorm:"column:password"`
	Email    string `json:"email" bson:"email" gorm:"column:email"`
	Salt     string `json:"salt,omitempty" bson:"salt" gorm:"column:salt"`
	//DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" bson:"deleted_at" gorm:"column:deleted_at"`
}

func (e *Example) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *Example) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

var _ schema.Tabler = (*Example)(nil)
var _ encoding.BinaryMarshaler = (*Example)(nil)
var _ encoding.BinaryUnmarshaler = (*Example)(nil)

func (e *Example) TableName() string {
	return "example"
}

func (e *Example) Sensitive() *Example {
	temp := *e
	temp.Password = ""
	temp.Salt = ""
	return &temp
}

func (e *Example) MD5PasswordWithSalt() string {
	m5 := md5.New()
	_, err := m5.Write([]byte(e.Password))
	if err != nil {
		return ""
	}
	m5.Write([]byte(e.Salt))
	if err != nil {
		return e.Password
	}
	return hex.EncodeToString(m5.Sum(nil))
}

func (e *Example) CheckPasswordMatch(password string) bool {
	temp := &Example{
		Password: password,
		Salt:     e.Salt,
	}
	if e.Password == temp.MD5PasswordWithSalt() {
		return true
	}
	return false
}

func (e *Example) ValidateUsernameRule() bool {
	if regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{5,19}$`).MatchString(e.Username) {
		return true
	}
	return false
}

func (e *Example) ValidatePasswordRule() bool {
	if regexp.MustCompile(`^.{6,100}$`).MatchString(e.Password) {
		return true
	}
	return false
}

func (e *Example) ValidateEmailRule() bool {
	if regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(e.Email) {
		return true
	}
	return false
}
