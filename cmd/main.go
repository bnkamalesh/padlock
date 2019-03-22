package main

import (
	"github.com/bnkamalesh/padlock/api/http"
	"github.com/bnkamalesh/webgo"
)

func main() {
	router := webgo.NewRouter(&webgo.Config{
		Port: "8080",
	}, http.Routes())
	router.Start()
}
