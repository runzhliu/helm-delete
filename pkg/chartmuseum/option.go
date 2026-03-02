package chartmuseum

// Option is a function that configures a Client.
type Option func(*options)

type options struct {
	url              string
	username         string
	password         string
	accessToken      string
	authHeader       string
	contextPath      string
	timeout          int64
	caFile           string
	certFile         string
	keyFile          string
	insecureSkipVerify bool
}

// URL sets the ChartMuseum server URL.
func URL(url string) Option {
	return func(opts *options) {
		opts.url = url
	}
}

// Username sets the basic auth username.
func Username(username string) Option {
	return func(opts *options) {
		opts.username = username
	}
}

// Password sets the basic auth password.
func Password(password string) Option {
	return func(opts *options) {
		opts.password = password
	}
}

// AccessToken sets the bearer token for authentication.
func AccessToken(accessToken string) Option {
	return func(opts *options) {
		opts.accessToken = accessToken
	}
}

// AuthHeader sets a custom header name to use for the access token.
func AuthHeader(authHeader string) Option {
	return func(opts *options) {
		opts.authHeader = authHeader
	}
}

// ContextPath sets a URL prefix for ChartMuseum (e.g. when behind a reverse proxy).
func ContextPath(contextPath string) Option {
	return func(opts *options) {
		opts.contextPath = contextPath
	}
}

// Timeout sets the request timeout in seconds.
func Timeout(timeout int64) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// CAFile sets the path to a CA certificate bundle.
func CAFile(caFile string) Option {
	return func(opts *options) {
		opts.caFile = caFile
	}
}

// CertFile sets the path to a TLS client certificate.
func CertFile(certFile string) Option {
	return func(opts *options) {
		opts.certFile = certFile
	}
}

// KeyFile sets the path to a TLS client private key.
func KeyFile(keyFile string) Option {
	return func(opts *options) {
		opts.keyFile = keyFile
	}
}

// InsecureSkipVerify disables TLS certificate verification.
func InsecureSkipVerify(insecureSkipVerify bool) Option {
	return func(opts *options) {
		opts.insecureSkipVerify = insecureSkipVerify
	}
}
