package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/domain"
	"github.com/grepsd/knowledge-database/pkg/infrastructure/persistence/pgsql"
	_ "github.com/lib/pq"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"database/sql"
)

func main() {
	//queryDispatcher := cqrs.NewSimpleQueryDispatcher()
	//commandDispatcher := cqrs.NewSimpleCommandDispatcher()
	connStr := "user=postgres password=tpassword dbname=knowledge-database sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	//driver, err := postgres.WithInstance(db, &postgres.Config{})
	//if err != nil {
	//	log.Fatal("driver : " + err.Error())
	//}
	//m, err := migrate.NewWithDatabaseInstance(
	//	"file://config/sql/migrations/",
	//	"postgres", driver)
	//if err != nil {
	//	log.Fatal("migrate : " + err.Error())
	//}
	//err = m.Up()
	//if err != nil {
	//	log.Fatal("up : " + err.Error())
	//}

	d := pgsql.New(db)
	a, err := d.CreateArticle(context.Background(), pgsql.CreateArticleParams{
		ID:      uuid.New(),
		Title:   "test",
		URL:     "http://go.com",
		SavedOn: time.Now(),
		ReadOn:  sql.NullTime{},
	})
	if err != nil {
		log.Fatal("cannot create article :" + err.Error())
	}
	c := func(a pgsql.Article) domain.Article {
		return domain.NewArticle(a.ID, a.Title, a.URL, a.SavedOn, a.ReadOn.Time, []domain.Tag{})
	}
	fmt.Printf("article : %+v\n", a)
	fmt.Printf("-- article : %+v\n", c(a))

	articles, err := d.GetAllArticles(context.Background())
	if err != nil {
		log.Fatal("cannot list articles :" + err.Error())
	}
	b, err := json.Marshal(articles)
	if err != nil {
		log.Fatal("cannot marshal articles :" + err.Error())
	}
	fmt.Println(string(b))
}
