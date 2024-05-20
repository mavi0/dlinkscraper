package router

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/reiver/go-telnet"
)

// Router represents a telnet router connection.
type Router struct {
	connection *telnet.Conn
}

// NewRouter creates a new Router instance and establishes a connection to the specified address.
func NewRouter(address string) (*Router, error) {
	connection, err := telnet.DialTo(address)
	if err != nil {
		return nil, err
	}
	return &Router{connection: connection}, nil
}

// Close closes the telnet connection.
func (r *Router) Close() {
	r.connection.Close()
}

// Expect reads data from the connection until the specified string is found. It
// returns the received data, the length of the received data, and any error
// that occurred.
func (r *Router) Expect(s string) (string, int, error) {
	buffer := make([]byte, 4096)
	for {
		singleBuffer := make([]byte, 1)
		n, err := r.connection.Read(singleBuffer)
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			break
		}
		buffer = append(buffer, singleBuffer...)
		if strings.Contains(string(buffer), s) {
			return string(bytes.Trim(buffer, "\x00")), len(buffer), nil
		}
	}
	return "", 0, nil
}

// ReadLine reads a line of data from the connection. It returns the received
// line, the length of the received line, and any error that occurred.
func (r *Router) ReadLine() (string, int, error) {
	return r.Expect("\n")
}

// ExpectPrompt reads data from the connection until the router prompt is found.
// It returns an error if the prompt is not found.
func (r *Router) ExpectPrompt() error {
	_, _, err := r.Expect("~ # ")
	return err
}

// Write writes the specified string to the connection.
// It returns the number of bytes written and any error that occurred.
func (r *Router) Write(s string) (int, error) {
	return r.connection.Write([]byte(s))
}

// WriteCommand writes the specified command to the connection after expecting a
// prompt. It returns an error if the prompt is not found or if there was an
// error writing the command.
func (r *Router) WriteCommand(s string) error {
	if err := r.ExpectPrompt(); err != nil {
		return fmt.Errorf("command expected a prompt before running: %w", err)
	}
	_, err := r.connection.Write([]byte(s))
	return err
}
