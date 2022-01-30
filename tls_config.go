package lever

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/go-lever/lever/config"
)

const (
	defaultHTTPSPort = "443"
	defaultCertFile  = "cert.pem"
	defaultKeyFile   = "key.pem"
)

type tlsConfig struct {
	port     string
	certFile string
	KeyFile  string
	domain   string
	email    string
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

func (tls *tlsConfig) generateDevCertificate() {

	localhost := []string{"localhost", "127.0.0.1"}

	if _, err := os.Stat(tls.certFile); os.IsNotExist(err) {
		err := tls.generateCertificate(localhost)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(tls.KeyFile); os.IsNotExist(err) {
		err := tls.generateCertificate(localhost)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (tls *tlsConfig) generateCertificate(hosts []string) error {

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key : %s", err.Error())
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %s", err.Error())
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Localhost Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	//writes the certificate file
	certOut, err := os.Create(tls.certFile)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", tls.certFile, err.Error())
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", tls.certFile, err.Error())
	}

	if err := certOut.Close(); err != nil {
		return fmt.Errorf("error closing %s : %s", tls.certFile, err.Error())
	}

	//writes the key file
	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", tls.KeyFile, err.Error())
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("unable to marshal private key: %s", err.Error())
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed write data to %s: %s", tls.KeyFile, err.Error())
	}

	if err := keyOut.Close(); err != nil {
		return fmt.Errorf("error closing %s: %s", tls.KeyFile, err.Error())
	}

	return nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}
