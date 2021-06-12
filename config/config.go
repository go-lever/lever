package config

import "os"

func DevMode() bool {
	if os.Getenv("LEVER_MODE") == "prod" {
		return false
	}
	return true
}

func TLSDomain() string {
	return os.Getenv("LEVER_DOMAIN")
}

func TLSEmail() string {
	return os.Getenv("LEVEL_EMAIL")
}