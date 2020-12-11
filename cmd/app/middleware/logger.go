package middleware

import (
	"log"
	"net/http"
)

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request){
		log.Printf("START: %s %s", request.Method, request.URL.Path)
		handler.ServeHTTP(writer,request)
		log.Printf("FINISH: %s %s", request.Method, request.URL.Path)
	})
}