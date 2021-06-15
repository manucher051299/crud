package middleware

import (
	"log"
	"net/http"
)

func Basic(auth func(login, pass string) bool) func(handler http.Handler) http.Handler {

	return func(handler http.Handler) http.Handler {

		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			username, password, ok := request.BasicAuth()

			if !ok {
				log.Print("Cannot parse password and login")
				http.Error(writer, http.StatusText(401), 401)
				return
			}
			if !auth(username, password) {
				http.Error(writer, http.StatusText(401), 401)
			}

			handler.ServeHTTP(writer, request)

		})
	}

}
