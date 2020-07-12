package http

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	articles articleHTTPHandler
}

func NewServer(articles articleHTTPHandler) *Server {
	return &Server{articles: articles}
}

func (s *Server) Init() {
	http.HandleFunc("/articles", auth(s.articles.routeCollection()))
	http.HandleFunc("/articles/", auth(s.articles.routeItem()))
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, password, ok := r.BasicAuth()
		if ok && user == "grepsd" && password == "test" {
			ctx = context.WithValue(ctx, "auth", true)
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="Knowledge Base"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			ctx = context.WithValue(ctx, "auth", false)
			return
		}
		requestWithContext := r.WithContext(ctx)
		next(w, requestWithContext)
	}
}
