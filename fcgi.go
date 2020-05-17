package main

import (
	"net"
	"net/http"
	"net/http/fcgi"

	iris "github.com/kataras/iris/v12"
)

func runFcgi(handler *iris.Application, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	handler.Build()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.URL.String())
		handler.ServeHTTP(w, r)
	})
	err = fcgi.Serve(listener, nil)
	return err
}
