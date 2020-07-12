package http

import (
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
	http.HandleFunc("/articles", s.articles.routeCollection())
	http.HandleFunc("/articles/", s.articles.routeItem())
	log.Fatal(http.ListenAndServe(":8888", nil))
}
