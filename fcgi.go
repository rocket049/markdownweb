package main

import (
	"net"
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

	err = fcgi.Serve(listener, handler)
	return err
}
