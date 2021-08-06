package main

import "sync"

type SnipeConfig struct {
	USERNAME string
	TOKEN    string
	DROPTIME int
	OFFSET   int
}

type AuthConfig struct {
	EMAIL    string
	PASSWORD string
	TOKEN    string
}

type DropTime struct {
	UNIX int `json:"unix"`
}

type snipeFunc func(string, string, *sync.WaitGroup)
