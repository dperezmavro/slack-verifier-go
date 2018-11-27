package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dperezmavro/slack-verifier-go/middleware"
)

// simplest handler
func healthCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
}

func healthCheckHandlerSimple(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok2")

}

type handler func(http.Handler) http.Handler

func addMiddleware(h http.Handler, middlwr ...handler) http.Handler {
	for _, mid := range middlwr {
		h = mid(h)
	}
	return h
}

func main() {
	slackKey := []byte(os.Getenv(""))
	http.Handle("/path", middleware.VerifySlackMessage(slackKey)(http.HandlerFunc(healthCheckHandlerSimple)))
	http.Handle("/path2",
		addMiddleware(
			healthCheckHandler(),
			middleware.VerifySlackMessage(slackKey),
		),
	)

	http.ListenAndServe(":8080", nil)
}
