package parser

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Address represents a network address in the form host:port.
type Address struct {
	// Host is the hostname or IP address part of the address.
	Host string

	// Port is the numeric TCP/UDP port.
	Port int
}

// BaseURL represents a base URL, split into protocol, domain and port parts.
type BaseURL struct {
	// Protocol is the scheme part of the URL (e.g. "http", "https").
	Protocol string

	// Domain is the hostname part of the URL (e.g. "example.com").
	Domain string

	// Port is the port number as a string (e.g. "8080").
	Port string
}

// String returns the string representation of Address in the format host:port.
func (a Address) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

// Set parses a string of the form "host:port" and assigns values to Address.
// Returns an error if the string is invalid or the port cannot be parsed.
func (a *Address) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("некорректный формат, ожидается host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return fmt.Errorf("ошибка парсинга порта: %v", err)
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

// String returns the string representation of BaseURL in the format protocol://domain:port.
func (a BaseURL) String() string {
	return a.Protocol + "://" + a.Domain + ":" + a.Port
}

// Set parses a URL string and assigns its components to BaseURL.
// The expected format is protocol://domain:port (e.g., http://example.com:8080).
// Returns an error if the string cannot be parsed or does not contain a port.
func (a *BaseURL) Set(value string) error {
	parsedURL, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("ошибка парсинга URL: %v", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("некорректный URL, ожидается формат protocol://domain:port")
	}

	// Разделяем домен и порт
	hostParts := strings.Split(parsedURL.Host, ":")
	a.Protocol = parsedURL.Scheme
	a.Domain = hostParts[0]

	if len(hostParts) > 1 {
		a.Port = hostParts[1]
	} else {
		return errors.New("некорректный URL, должен содержать порт (например, http://example.com:8080)")
	}

	return nil
}
