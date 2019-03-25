package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// The maximum amount of memory to use when parsing multipart requests.
const MultipartMemoryLimit = 32 * 1024 * 1024

// A request handler.
type RequestHandler struct {
	target string
	echo bool
	status int
	shutdown chan error
}

// Constructs a new request handler.
func NewHandler(target string, echo bool, status int) *RequestHandler {
	return &RequestHandler {
		target: target,
		echo: echo,
		status: status,
		shutdown: make(chan error),
	}
}

// Handles a request by echoing its contents and storing them to disk.
func (rh *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	if merr := r.ParseMultipartForm(MultipartMemoryLimit); merr == nil {
		err = rh.HandleMultipart(r)
	} else {
		err = rh.HandleRequest(r)
	}
	w.WriteHeader(rh.status)

	// Determine if the request was successfully handled.
	if err != nil {
		rh.shutdown <- err
	} else {
		rh.shutdown <- nil
	}
}

// Handles a multipart request by attempting to store its contents on disk.
func (rh *RequestHandler) HandleMultipart(r *http.Request) error {
	writer, err := os.Create(rh.target)
	if err != nil {
		return err
	}
	defer writer.Close()
	archive := zip.NewWriter(writer)
	defer archive.Close()

	// Copy files from the request into the archive.
	for _, headers := range r.MultipartForm.File {
		for _, header := range headers {
			if rh.echo {
				fmt.Printf("received file: %s\n", header.Filename)
			}
			output, err := archive.Create(header.Filename)
			if err != nil {
				return err
			}
			input, err := header.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(output, input)
			input.Close()
			if err != nil {
				return err
			}
		}
	}

	// Create a file to hold the submitted form data.
	for key, values := range r.MultipartForm.Value {
		output, err := archive.Create("form-data.txt")
		if err != nil {
			return err
		}

		message := fmt.Sprintf("%s = %s\n", key, values)
		if rh.echo {
			fmt.Printf("received form data: %s", message)
		}
		_, err = output.Write([]byte(message))
		if err != nil {
			return err
		}
	}
	return nil
}

// Handles a request with a single part.
func (rh *RequestHandler) HandleRequest(r *http.Request) error {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(content)
	if rh.echo {
		if _, err := io.Copy(os.Stdout, reader); err != nil {
			return err
		}
	}
	if rh.target != "" {
		_, err := reader.Seek(0, io.SeekStart)
		if err != nil {
			panic("critical error: could not reset byte reader")
		}
		if err := ioutil.WriteFile(rh.target, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

// The shutdown channel.
func (rh *RequestHandler) Shutdown() chan error {
	return rh.shutdown
}
