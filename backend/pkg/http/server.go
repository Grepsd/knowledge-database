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
	http.HandleFunc("/articles", s.articles.Articles())
	http.HandleFunc("/articles/", s.articles.Article())
	log.Fatal(http.ListenAndServe(":8888", nil))
}
