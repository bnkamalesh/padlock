package main

import (
	"github.com/bnkamalesh/webgo"
)

func main() {
	router := webgo.NewRouter(&webgo.Config{
		Port: "8080",
	}, nil)
	router.Start()
}
