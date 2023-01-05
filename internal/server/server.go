package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetNewServer(ip string, port uint16, router *httprouter.Router) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: router,
	}
}

func GetNewSSLServer(ip string, port uint16, router *httprouter.Router) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: router,
	}
}
