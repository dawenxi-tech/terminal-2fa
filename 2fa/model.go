package main

import "time"

type Entry struct {
	ID       string    `json:"id"`
	Name     string    `json:"name,omitempty"`
	Seed     string    `json:"seed,omitempty"`
	Order    int       `json:"order,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
}
