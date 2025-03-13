package main

import (
	http "github.com/thxhix/shortener/internal/app/server"
)

func main() {
	server := http.NewServer()
	err := server.StartPooling()
	if err != nil {
		panic(err)
	}
}
