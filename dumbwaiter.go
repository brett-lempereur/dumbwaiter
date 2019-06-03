// A simple command-line tool for recording HTTP requests.
package main

import (
	"context"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"time"
)

// Command-line arguments.
var (
	verbose = kingpin.Flag("verbose", "verbose output").Short('v').Default("true").Bool()
	address = kingpin.Flag("address", "server address").Short('a').Default(":8080").String()
	status = kingpin.Flag("status", "response status code").Default("200").Short('s').Int()
	target = kingpin.Arg("path", "target path").Required().String()
)

// Launches the command-line application and web server.
func main() {
	kingpin.Parse()

	// Build the server.
	handler := NewHandler(*target, *verbose, *status)
	server := http.Server {
		Addr: *address,
		Handler: handler,
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Launch the server in a separate routine and wait for shutdown.
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: %s", err)
		}
	}()
	select {
		case err := <-handler.Shutdown():
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error: %s", err)
				os.Exit(1)
			}
			err = server.Shutdown(context.Background())
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error: %s", err)
				os.Exit(2)
			}

	}
}
