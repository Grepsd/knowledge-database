package main

import (
	article2 "github.com/grepsd/knowledge-database/internal/application/command/article"
	article3 "github.com/grepsd/knowledge-database/internal/application/handler/article"
	"github.com/grepsd/knowledge-database/internal/application/query"
	"github.com/grepsd/knowledge-database/internal/domain/model/entity"
	"github.com/grepsd/knowledge-database/internal/infrastructure/persistence/sql/repository"
	"log"
	"time"
)

func main() {
	repo := repository.NewArticle()

	u, _ := entity.NewArticleURL("http://go.com")
	//a, _ := entity.CreateArticle(entity.NewArticleTitle("test r"), u)

	command := article2.NewSaveArticle(u, entity.NewArticleTitle("test cmd"), entity.NewArticleSavedDateTime(time.Now()))
	handler := article3.NewSaveArticle(repo)
	_, err := handler.Handle(command)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Println("failure to save article : " + err.Error())
	}

	query := query.NewFindAll(DisplayArticles)
	queryHandler := article3.NewFindAll(repo)
	queryHandler.Handle(query)
	return
	articles, err := repo.FindAll()
	if err != nil {
		log.Fatal(err)
	}
	var article *entity.Article
	for _, article = range articles {
		log.Printf("article title %s url %s\n", article.Title(), article.Url())
	}
}

func DisplayArticles(articles []*entity.Article, err error) {
	if err != nil {
		log.Fatal("query failed, catched in result handler : " + err.Error())
	}
	var article *entity.Article
	for _, article = range articles {
		log.Printf("article title %s url %s\n", article.Title(), article.Url())
	}
}