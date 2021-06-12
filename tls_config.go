package lever

import "github.com/go-lever/lever/config"

const (
	defaultHTTPSPort = "443"
	defaultCertFile = ""
	defaultKeyFile = ""
)

type tlsConfig struct {
	port string
	certFile string
	KeyFile string
	domain string
	email string
}

func newTLSConfig() *tlsConfig {
	return &tlsConfig{
		port:     defaultHTTPSPort,
		certFile: defaultCertFile,
		KeyFile:  defaultKeyFile,
		domain:   config.TLSDomain(),
		email:    config.TLSEmail(),
	}
}