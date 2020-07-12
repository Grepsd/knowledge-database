package main

import (
	"github.com/grepsd/knowledge-database/pkg"
	"github.com/grepsd/knowledge-database/pkg/sql"
)

func main() {
	db := sql.NewDB("user=postgres password=tpassword database=knowledge-database sslmode=disable")
	articleRepository := sql.NewArticleRepository(db)
	s := pkg.NewServer(articleRepository)
	s.Init()
}
