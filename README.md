# go-http-default-request-header-transporter [![.github/workflows/check.yml](https://github.com/moznion/go-http-default-request-header-transporter/actions/workflows/check.yml/badge.svg)](https://github.com/moznion/go-http-default-request-header-transporter/actions/workflows/check.yml) [![codecov](https://codecov.io/gh/moznion/go-http-default-request-header-transporter/branch/main/graph/badge.svg?token=AFQ94U61WY)](https://codecov.io/gh/moznion/go-http-default-request-header-transporter) [![Go Reference](https://pkg.go.dev/badge/github.com/moznion/go-http-default-request-header-transporter#section-readme.svg)](https://pkg.go.dev/github.com/moznion/go-http-default-request-header-transporter#section-readme)

A utility [http.Transport](https://pkg.go.dev/net/http#Transport) to inject the given default request header.

## Synopsis

```go
import (
	"http"
	"time"

	"github.com/moznion/go-http-default-request-header-transporter"
)

func main() {
	defaultHeader := http.Header{}
	defaultHeader.Set("user-agent", "custom-UA/0.0.1")
	defaultHeader.Add("x-test", "foo")
	defaultHeader.Add("x-test", "bar")

	hc := &http.Client{}
	hc.Timeout = 3 * time.Second
	hc.Transport = transporter.NewDefaultRequestHeaderTransporter(hc.Transport, defaultHeader)
	resp, err := hc.Get(httpURL) // <= this request header has the values of `defaultHeader`
}
```

Please also refer to the [examples_test.go](./examples_test.go).

## Documentations

[![Go Reference](https://pkg.go.dev/badge/github.com/moznion/go-http-default-request-header-transporter#section-readme.svg)](https://pkg.go.dev/github.com/moznion/go-http-default-request-header-transporter#section-readme)

## Author

moznion (<moznion@mail.moznion.net>)

