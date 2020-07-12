package main

import (
	"github.com/grepsd/knowledge-database/pkg/http"
	"github.com/grepsd/knowledge-database/pkg/sql"
)

func main() {
	db := sql.NewDB("user=postgres password=tpassword database=knowledge-database sslmode=disable")
	articleRepository := sql.NewArticleRepository(&db)
	httpHelpers := http.NewHelpers()
	articleHandler := http.NewArticleHTTPHandler(httpHelpers, articleRepository)
	s := http.NewServer(articleHandler)
	s.Init()
}
