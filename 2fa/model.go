package main

import (
	"math/rand"
	"time"
)

type Entry struct {
	ID       string    `json:"id"`
	Name     string    `json:"name,omitempty"`
	Seed     string    `json:"seed,omitempty"`
	Order    int       `json:"order,omitempty"`
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
