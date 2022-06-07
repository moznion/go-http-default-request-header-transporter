package transporter

import "net/http"

// DefaultRequestHeaderTransporter is a utility http.Transport to inject the default request header.
type DefaultRequestHeaderTransporter struct {
	originalTransporter http.RoundTripper
	header              http.Header
}

// NewDefaultRequestHeaderTransporter makes a new instance of DefaultRequestHeaderTransporter.
// Created DefaultRequestHeaderTransporter injects given header value into every requests' header.
// If given originalTransporter is nil, it uses http.DefaultTransport instead.
func NewDefaultRequestHeaderTransporter(originalTransporter http.RoundTripper, header http.Header) *DefaultRequestHeaderTransporter {
	if originalTransporter == nil {
		originalTransporter = http.DefaultTransport
	}

	return &DefaultRequestHeaderTransporter{
		originalTransporter: originalTransporter,
		header:              header,
	}
}

// RoundTrip inherits the original transporter's RoundTrip with injecting the default header values.
func (t *DefaultRequestHeaderTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, values := range t.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return t.originalTransporter.RoundTrip(req)
}
