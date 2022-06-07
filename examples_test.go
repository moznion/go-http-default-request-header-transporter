package transporter

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func ExampleDefaultRequestHeaderTransporter_RoundTrip() {
	listenerPortCh := make(chan int)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// output custom default request header values
		for _, h := range r.Header.Values("user-agent") {
			fmt.Printf("%s\n", h)
		}
		for _, h := range r.Header.Values("x-test") {
			fmt.Printf("%s\n", h)
		}

		w.WriteHeader(200)
	})

	server := &http.Server{
		Handler: mux,
	}

	go func(server *http.Server) {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
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

	port := <-listenerPortCh
	httpURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	defaultHeader := http.Header{}
	defaultHeader.Set("user-agent", "custom-UA/0.0.1")
	defaultHeader.Add("x-test", "foo")
	defaultHeader.Add("x-test", "bar")

	hc := &http.Client{}
	hc.Timeout = 3 * time.Second
	hc.Transport = NewDefaultRequestHeaderTransporter(hc.Transport, defaultHeader)
	resp, err := hc.Get(httpURL)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	// Output:
	// custom-UA/0.0.1
	// foo
	// bar
}
