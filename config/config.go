package config

import "os"

func DevMode() bool {
	if os.Getenv("LEVER_MODE") == "prod" {
		return false
	}
	return true
}

func LocalMode() bool {
	if os.Getenv("LEVER_LOCAL") == "true" {
		return true
	}
	return false
}

func TLSDomain() string {
	return os.Getenv("LEVER_DOMAIN")
}

func TLSEmail() string {
	return os.Getenv("LEVEL_EMAIL")
}
