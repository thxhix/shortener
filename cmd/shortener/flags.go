package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ServerHostParams struct {
	Host string
	Port int
}

type BaseURL struct {
	Protocol string
	Domain   string
	Port     string
}

var hostParams ServerHostParams
var destinationURL BaseURL

func (a ServerHostParams) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *ServerHostParams) Set(s string) error {
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

func (a BaseURL) String() string {
	return a.Protocol + "://" + a.Domain + ":" + a.Port
}

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

func parseFlags() {
	flag.Var(&hostParams, "a", "Net address host:port (например, localhost:8080)")
	flag.Var(&destinationURL, "d", "Destination URL (например, http://example.com:8080)")

	flag.Parse()
}
