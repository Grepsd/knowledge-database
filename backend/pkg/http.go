package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/article"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Server struct {
	db *DB
}

func NewServer(db *DB) *Server {
	return &Server{db}
}

func (s *Server) Init() {
	http.HandleFunc("/articles", s.Articles())
	http.HandleFunc("/articles/", s.Article())
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func (s *Server) Articles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateArticle(w, r)
			break
		case http.MethodGet:
			s.ListArticles(w, r)
			break
		}
	}
}

func (s *Server) Article() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.GetArticleById(w, r)
			break
		}
	}
}

func (s *Server) CreateArticle(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatal("failed to read body : " + err.Error())
	}

	art := article.Article{}
	err = json.Unmarshal(payload, &art)
	if err != nil {
		log.Fatal("failed to unmarshall : " + err.Error())
	}
	art.ID = uuid.New()

	err = s.db.CreateArticle(&art)
	if err != nil {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "failed to create article : "+err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(art)
	if err != nil {
		log.Fatal("failed to marshall : " + err.Error())
	}
	_, err = fmt.Fprintf(w, string(resp))
}

func (s *Server) ListArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := s.db.ListArticles()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to create article : "+err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(articles)
	if err != nil {
		log.Fatal("failed to marshall : " + err.Error())
	}
	_, err = fmt.Fprintf(w, string(resp))
}

func (s *Server) GetArticleById(w http.ResponseWriter, r *http.Request) {
	var art article.Article
	uri := r.RequestURI
	fmt.Println(uri)
	parts := strings.Split(uri, "/")
	key := parts[len(parts)-1]
	if id, err := uuid.FromBytes([]byte(key)); err != nil && id.String() != "00000000-0000-0000-0000-000000000000" {
		art, err = s.db.GetOneById(id)
	} else {
		re, err := regexp.Compile("[^a-zA-Z0-9-_]+")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "failed to compile regexp : "+err.Error())
			return
		}
		res := bytes.ToLower(re.ReplaceAll([]byte(key), []byte("_")))

		art, err = s.db.GetOneBySlug(string(res))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "failed to get article : "+err.Error())
			return
		}

	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(art)
	if err != nil {
		log.Fatal("failed to marshall : " + err.Error())
	}
	_, err = fmt.Fprintf(w, string(resp))
}
