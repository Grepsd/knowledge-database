package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/article"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	articleRepository article.ReadWriteRepositoryer
}

func NewServer(articleRepository article.ReadWriteRepositoryer) *Server {
	return &Server{articleRepository: articleRepository}
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
		case http.MethodPut:
			s.PutArticle(w, r)
			break
		}
	}
}

func (s *Server) CreateArticle(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic("failed to read body : " + err.Error())
	}

	art := new(article.Article)
	err = json.Unmarshal(payload, art)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to json.decode payload: %w", err))
		return
	}
	art.ID = uuid.New()
	slug, err := article.GenerateSlugFromTitle(art.Title)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to generate article slug from title: %w", err))
		return
	}
	art.Slug = slug

	err = s.articleRepository.Create(*art)
	if err != nil {
		if strings.Contains(err.Error(), " duplicate key value violates") {
			w.WriteHeader(http.StatusConflict)
			s.writeResponse(w, fmt.Errorf("failed to create article: %w", err).Error())
			return
		}
		s.writeErrorResponse(w, r, fmt.Errorf("failed to create article: %w", err))
		return
	}

	s.respondWithJSON(w, http.StatusCreated, art)
}

func (s *Server) ListArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := s.articleRepository.GetAll()
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to list articles: %w", err))
		return
	}
	s.respondWithJSON(w, http.StatusOK, articles)
}

func (s *Server) GetArticleById(w http.ResponseWriter, r *http.Request) {
	var art article.Article
	key := s.getLastSegmentFromURI(r)
	if id, err := uuid.Parse(key); err == nil && id.String() != "00000000-0000-0000-0000-000000000000" {
		art, err = s.articleRepository.GetOneById(id)
		if err != nil {
			s.writeErrorResponse(w, r, fmt.Errorf("failed to retrieve article by id: %w", err))
			return
		}
	} else {
		title, err := article.GenerateSlugFromTitle(key)
		if err != nil {
			s.writeErrorResponse(w, r, err)
			return
		}
		art, err = s.articleRepository.GetOneBySlug(string(title))
		if err != nil {
			s.writeErrorResponse(w, r, fmt.Errorf("failed to get article by slug: %w", err))
			return
		}

	}

	s.respondWithJSON(w, http.StatusOK, art)
}

func (s *Server) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Print(err)
}

func (s *Server) writeResponse(w http.ResponseWriter, data string) error {
	_, err := fmt.Fprint(w, data)
	if err != nil {
		err = fmt.Errorf("failed to write to responseWriter : %w", err)
	}
	return err
}

func (s *Server) writeJsonContentType(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func (s *Server) PutArticle(w http.ResponseWriter, r *http.Request) {
	art := new(article.Article)
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf(" : %w", err))
		return
	}
	err = json.Unmarshal(payload, art)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to json.decode payload : %w", err))
		return
	}
	uriIdentifier := s.getLastSegmentFromURI(r)
	if id, err := uuid.FromBytes([]byte(uriIdentifier)); err != nil && id.String() != "00000000-0000-0000-0000-000000000000" {
		art.ID = id
		s.updateArticleByID(w, r, art)
		return
	}
	slug, err := article.GenerateSlugFromTitle(uriIdentifier)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("uriIdentifier : %w", err))
		return
	}
	art.Slug = slug

	registeredArticle, err := s.articleRepository.GetOneBySlug(slug)
	if err != nil {
		if errors.As(err, fmt.Errorf("article not found")) {
			art.ID = uuid.New()
			err = s.articleRepository.Create(*art)
			if err != nil {
				s.writeErrorResponse(w, r, fmt.Errorf("failed to create article: %w", err))
				return
			}
			newArticle, err := s.articleRepository.GetOneById(art.ID)
			if err != nil {
				s.writeErrorResponse(w, r, fmt.Errorf("failed to retrieve newly created article: %w", err))
				return
			}
			err = s.respondWithJSON(w, http.StatusCreated, newArticle)
			if err != nil {
				s.writeErrorResponse(w, r, fmt.Errorf("failed to write response : %w", err))
				return
			}
			return
		}
	}

	art.ID = registeredArticle.ID
	err = s.articleRepository.Update(*art)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("update product failed : %w", err))
		return
	}
	newArticle, err := s.articleRepository.GetOneById(art.ID)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to retrieve newly created article: %w", err))
		return
	}
	err = s.respondWithJSON(w, http.StatusOK, newArticle)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to write response : %w", err))
		return
	}
}

func (s *Server) getLastSegmentFromURI(r *http.Request) string {
	uri := r.RequestURI
	parts := strings.Split(uri, "/")
	return parts[len(parts)-1]
}

func (s *Server) writeJson(w http.ResponseWriter, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshall: %w", err)
	}
	s.writeJsonContentType(w)
	_, err = w.Write(payload)
	if err != nil {
		panic("failed to write response: " + err.Error())
	}
	return err
}

func (s *Server) updateArticleByID(w http.ResponseWriter, r *http.Request, art *article.Article) {
	_, err := s.articleRepository.GetOneById(art.ID)
	if err != nil {
		if errors.As(err, errors.New("article not found")) {
			http.NotFound(w, r)
			return
		}
		s.writeErrorResponse(w, r, fmt.Errorf("failed to load article : %w", err))
		return
	}
	slug, err := article.GenerateSlugFromTitle(art.Title)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to generate slug : %w", err))
		return
	}
	art.Slug = slug
	err = s.articleRepository.Update(*art)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to update article : %w", err))
		return
	}
	err = s.respondWithJSON(w, http.StatusOK, art)
	if err != nil {
		s.writeErrorResponse(w, r, fmt.Errorf("failed to write response : %w", err))
		return
	}
	return
}

func (s *Server) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	err := s.writeJson(w, data)
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}
	return err
}
