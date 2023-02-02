package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if err := errMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func errMain() error {
	err := loadConfig("config.yaml")
	if err != nil {
		return err
	}

	err = initRuleGlobCache()
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:    Config.Listen,
		Handler: reverseProxy{},
	}
	return server.ListenAndServe()
}
