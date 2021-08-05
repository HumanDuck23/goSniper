package main

type SnipeConfig struct {
	USERNAME string
	TOKEN string
	DROPTIME string
	OFFSET int
}

type AuthConfig struct {
	EMAIL string
	PASSWORD string
	TOKEN string
}

type DropTime struct {
	UNIX int
}
