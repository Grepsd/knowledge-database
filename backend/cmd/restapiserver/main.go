package main

import (
	"github.com/grepsd/knowledge-database/pkg"
	"github.com/grepsd/knowledge-database/pkg/http"
	"github.com/grepsd/knowledge-database/pkg/sql"
)

func main() {
	db := sql.NewDB("host=db user=postgres password=tpassword database=knowledge-database sslmode=disable")
	articleRepository := sql.NewArticleRepository(&db)
	tagRepository := sql.NewTagRepository(&db)
	httpHelpers := http.NewHelpers()
	articleHandler := http.NewArticleHTTPHandler(httpHelpers, articleRepository)
	tagHandler := http.NewTagHTTPHandler(httpHelpers, tagRepository)
	s := http.NewServer(articleHandler, tagHandler)
	metrics := pkg.NewMetrics()
	s.RegisterHandler("/metrics", metrics.GetMetrics())

	s.Init(metrics)
}
