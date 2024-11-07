package main

import "time"

type Entry struct {
	Name     string    `json:"name,omitempty"`
	Seed     string    `json:"seed,omitempty"`
	Desc     string    `json:"desc,omitempty"`
	Order    int       `json:"order,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
}
