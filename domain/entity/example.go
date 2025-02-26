package entity

import (
	"encoding"
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

var _ encoding.BinaryUnmarshaler = (*Example)(nil)

// UnmarshalBinary Reids Scan use
func (e *Example) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

var _ encoding.BinaryMarshaler = (*Example)(nil)

// MarshalBinary Reids Set use
func (e *Example) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

var _ schema.Tabler = (*Example)(nil)

// TableName Gorm use
func (e *Example) TableName() string {
	return "example"
}

func (e *Example) Sensitive() *Example {
	temp := *e
	temp.Password = ""
	temp.Salt = ""
	return &temp
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
