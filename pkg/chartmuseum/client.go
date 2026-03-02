package chartmuseum

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Client is an HTTP client for communicating with a ChartMuseum server.
type Client struct {
	*http.Client
	opts options
}

// NewClient creates a new ChartMuseum client with the given options.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{}
	for _, opt := range opts {
		opt(&c.opts)
	}

	transport, err := newTransport(c.opts)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(c.opts.timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	c.Client = &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	return c, nil
}

func newTransport(opts options) (*http.Transport, error) {
	tlsCfg, err := newClientTLS(opts)
	if err != nil {
		return nil, err
	}
	return &http.Transport{
		TLSClientConfig: tlsCfg,
		Proxy:           http.ProxyFromEnvironment,
	}, nil
}

func newClientTLS(opts options) (*tls.Config, error) {
	cfg := &tls.Config{
		InsecureSkipVerify: opts.insecureSkipVerify, //nolint:gosec
	}

	if opts.certFile != "" && opts.keyFile != "" {
		cert, err := tls.LoadX509KeyPair(opts.certFile, opts.keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client cert/key (%s, %s): %w", opts.certFile, opts.keyFile, err)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}

	if opts.caFile != "" {
		ca, err := os.ReadFile(opts.caFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file %s: %w", opts.caFile, err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return nil, fmt.Errorf("failed to parse CA certificate from %s", opts.caFile)
		}
		cfg.RootCAs = pool
	}

	return cfg, nil
}
