package transporter

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetHeader(t *testing.T) {
	errCh := make(chan error)
	listenerPortCh := make(chan int)
	reqHeaderCh := make(chan http.Header, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqHeaderCh <- r.Header
		w.WriteHeader(200)
	})

	server := &http.Server{
		Handler: mux,
	}

	go func(server *http.Server) {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			errCh <- err
		}

		listenerPortCh <- listener.Addr().(*net.TCPAddr).Port

		err = server.Serve(listener)
		if err != nil {
			log.Printf("server shutdown: %s", err)
		}
	}(server)
	defer func() {
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Printf("[error] failed to shutdown the http server: %s", err)
		}
	}()

	ctx, canceler := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer func() {
		canceler()
	}()

	var port int
	select {
	case port = <-listenerPortCh:
	case err := <-errCh:
		t.Fatal(err)
	case <-ctx.Done():
		t.Fatal("it takes too long time to establish the http server")
	}

	httpURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	customUA := "custom-UA/0.0.1"
	testHeaderKey := "x-test"

	hc := &http.Client{}
	hc.Timeout = 3 * time.Second
	resp, err := hc.Get(httpURL)
	assert.NoError(t, err)
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	header := <-reqHeaderCh
	assert.NotEqual(t, customUA, header.Get("user-agent"))
	assert.Empty(t, header.Values(testHeaderKey))

	defaultHeader := http.Header{}
	defaultHeader.Set("user-agent", customUA)
	defaultHeader.Add(testHeaderKey, "foo")
	defaultHeader.Add(testHeaderKey, "bar")

	// set custom default request header
	hc.Transport = NewDefaultRequestHeaderTransporter(hc.Transport, defaultHeader)
	resp, err = hc.Get(httpURL)
	assert.NoError(t, err)
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	header = <-reqHeaderCh
	assert.Equal(t, customUA, header.Get("user-agent"))
	assert.Equal(t, 1, len(header.Values("user-agent")))
	assert.EqualValues(t, []string{"foo", "bar"}, header.Values(testHeaderKey))

	// nested transporters
	hc.Transport = NewDefaultRequestHeaderTransporter(hc.Transport, defaultHeader)
	resp, err = hc.Get(httpURL)
	assert.NoError(t, err)
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	header = <-reqHeaderCh
	assert.Equal(t, customUA, header.Get("user-agent"))
	assert.Equal(t, 1, len(header.Values("user-agent")))
	assert.EqualValues(t, []string{"foo", "bar", "foo", "bar"}, header.Values(testHeaderKey))
}
