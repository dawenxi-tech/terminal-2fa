package main

import (
	"encoding/json"
	"github.com/xlzd/gotp"
	"math/rand"
	"time"
)

type Entry struct {
	ID       string    `json:"id"`
	Name     string    `json:"name,omitempty"`
	Secret   *Secret   `json:"secret,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
}

type Secret string

func newSecret(msg string) *Secret {
	s := Secret(msg)
	return &s
}

func (s *Secret) Val() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

func (s *Secret) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte(""), nil
	}
	msg, err := encrypt(string(*s))
	if err != nil {
		return nil, err
	}
	return json.Marshal(msg)
}

func (s *Secret) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	msg, err := decrypt(str)
	if err != nil {
		return err
	}
	*s = Secret(msg)
	return nil
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func newId() string {
	const litters = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	chars := [16]byte{}
	for i := range chars {
		chars[i] = litters[random.Intn(len(litters))]
	}
	return string(chars[:])
}

func isValidTOTPCode(code string) (valid bool) {
	defer func() {
		if x := recover(); x != nil {
			valid = false
		}
	}()
	gotp.NewDefaultTOTP(code).At(time.Now().UnixNano())
	valid = true
	return
}
