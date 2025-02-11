package util

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log/slog"
	"os"
)

// TlsConfig creates a tls config for secure connection
// The communication uses mTLS for secure connection
// So a client certificate, client key and CA certificate are required
func TlsConfig(ctx context.Context, secure bool, clientCertPath, clientKeyPath, caCertPath string) (*tls.Config, error) {
	if secure == false {
		return nil, nil
	}
	// Load client cert and key
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		slog.ErrorContext(ctx, "error while loading client cert, terminating", slog.String("error", err.Error()))
		return nil, err
	}
	// Load CA cert
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		slog.ErrorContext(ctx, "error while loading server cert, terminating", slog.String("error", err.Error()))
		return nil, err
	}
	// Create a cert pool and add the embedded root and sub CA certs
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create tls.Config
	certConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return certConfig, nil
}
