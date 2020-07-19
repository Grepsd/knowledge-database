package http

import (
	"context"
	"github.com/grepsd/knowledge-database/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"time"
)

type Server struct {
	articles articleHTTPHandler
	tags     tagHTTPHandler
}

func NewServer(articles articleHTTPHandler, tags tagHTTPHandler) *Server {
	return &Server{articles: articles, tags: tags}
}

func (s *Server) Init(m *pkg.Metrics) {
	http.HandleFunc("/articles", requestCount(auth(s.articles.routeCollection()), m.RequestCounter(), m.RequestDuration()))
	http.HandleFunc("/articles/", requestCount(auth(s.articles.routeItem()), m.RequestCounter(), m.RequestDuration()))
	http.HandleFunc("/tags", requestCount(auth(s.tags.routeCollection()), m.RequestCounter(), m.RequestDuration()))
	http.HandleFunc("/tags/", requestCount(auth(s.tags.routeItem()), m.RequestCounter(), m.RequestDuration()))
	http.HandleFunc("/extracts", requestCount(auth(s.articles.extractsRoutes), m.RequestCounter(), m.RequestDuration()))
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func (s *Server) RegisterHandler(path string, handler http.Handler) {
	http.Handle(path, handler)
}

func requestCount(next http.HandlerFunc, counter prometheus.Counter, requestDuration prometheus.Histogram) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		counter.Inc()
		t := time.Now()
		next(w, r)
		duration := time.Now().Sub(t)
		microseconds := float64(duration.Microseconds()) / 1000
		requestDuration.Observe(float64(duration.Milliseconds()) + microseconds)
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
		}
		ctx := r.Context()
		//user, password, ok := r.BasicAuth()
		//if ok && user == "grepsd" && password == "test" {
		ctx = context.WithValue(ctx, "auth", true)
		//} else {
		//	w.Header().Set("WWW-Authenticate", `Basic realm="Knowledge Base"`)
		//	w.WriteHeader(401)
		//	w.Write([]byte("Unauthorised.\n"))
		//	ctx = context.WithValue(ctx, "auth", false)
		//	return
		//}
		requestWithContext := r.WithContext(ctx)
		next(w, requestWithContext)
	}
}
