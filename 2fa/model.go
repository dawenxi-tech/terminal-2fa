package main

import (
	"github.com/xlzd/gotp"
	"math/rand"
	"time"
)

type Entry struct {
	ID       string    `json:"id"`
	Name     string    `json:"name,omitempty"`
	Secret   string    `json:"secret,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
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
