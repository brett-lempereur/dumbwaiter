# Dumbwaiter

A simple command-line tool for recording HTTP requests.

## Installation

To install dumbwaiter, ensure you have a recent version of Go installed and run
the command:

```bash
go get github.com/brett-lempereur/dumbwaiter
```

This will install the `dumbwaiter` tool into your `GOPATH`.

## Usage

To run it, ensure that `$GOPATH/bin` is in your path and run:

```bash
dumbwaiter -v --address=:8080 --status=201 output.json
```

This will launch a server on port 8080 that will respond to any request on any
path with a HTTP 201 status code, capture the body of the request in a file,
and echo the body to the screen.

For further usage details, run:

```bash
$ dumbwaiter --help
usage: dumbwaiter [<flags>] <path>

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose          verbose output
  -a, --address=":8080"  server address
  -s, --status=200       response status code

Args:
  <path>  target path
```
