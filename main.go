package main

import (
	"context"
	"net/http"
	"os"

	"github.com/michaeldabbott/standardise/pkg/server"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	svr := server.NewFactory().Create()
	svr.Router.Get("/new", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("howdy fucker"))
	})

	if err := svr.Serve(context.Background()); err != nil {
		return err
	}
	return nil
}
