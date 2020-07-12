package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/article"
	"io/ioutil"
	"net/http"
	"strings"
)

type articleHTTPHandler struct {
	helpers    Helpers
	repository article.ReadWriteRepositoryer
}

func NewArticleHTTPHandler(s Helpers, r article.ReadWriteRepositoryer) articleHTTPHandler {
	return articleHTTPHandler{helpers: s, repository: r}
}

func (a *articleHTTPHandler) Articles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			a.Create(w, r)
			break
		case http.MethodGet:
			a.ListArticles(w, r)
			break
		}
	}
}

func (a *articleHTTPHandler) Article() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			a.GetArticleById(w, r)
			break
		case http.MethodPut:
			a.PutArticle(w, r)
			break
		}
	}
}

func (a *articleHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic("failed to read body : " + err.Error())
	}

	art := new(article.Article)
	err = json.Unmarshal(payload, art)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to json.decode payload: %w", err))
		return
	}
	art.ID = uuid.New()
	slug, err := article.GenerateSlugFromTitle(art.Title)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to generate article slug from title: %w", err))
		return
	}
	art.Slug = slug

	err = a.repository.Create(*art)
	if err != nil {
		if strings.Contains(err.Error(), " duplicate key value violates") {
			w.WriteHeader(http.StatusConflict)
			a.helpers.writeResponse(w, fmt.Errorf("failed to create article: %w", err).Error())
			return
		}
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to create article: %w", err))
		return
	}

	a.helpers.respondWithJSON(w, http.StatusCreated, art)
}

func (a *articleHTTPHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := a.repository.GetAll()
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to list articles: %w", err))
		return
	}
	a.helpers.respondWithJSON(w, http.StatusOK, articles)
}

func (a *articleHTTPHandler) GetArticleById(w http.ResponseWriter, r *http.Request) {
	var art article.Article
	key := a.helpers.getLastSegmentFromURI(r)
	if id, err := uuid.Parse(key); err == nil && id.String() != "00000000-0000-0000-0000-000000000000" {
		art, err = a.repository.GetOneById(id)
		if err != nil {
			a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to retrieve article by id: %w", err))
			return
		}
	} else {
		title, err := article.GenerateSlugFromTitle(key)
		if err != nil {
			a.helpers.writeErrorResponse(w, r, err)
			return
		}
		art, err = a.repository.GetOneBySlug(string(title))
		if err != nil {
			a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to get article by slug: %w", err))
			return
		}

	}

	a.helpers.respondWithJSON(w, http.StatusOK, art)
}

func (a *articleHTTPHandler) PutArticle(w http.ResponseWriter, r *http.Request) {
	art := new(article.Article)
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf(" : %w", err))
		return
	}
	err = json.Unmarshal(payload, art)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to json.decode payload : %w", err))
		return
	}
	uriIdentifier := a.helpers.getLastSegmentFromURI(r)
	if id, err := uuid.FromBytes([]byte(uriIdentifier)); err != nil && id.String() != "00000000-0000-0000-0000-000000000000" {
		art.ID = id
		a.updateArticleByID(w, r, art)
		return
	}
	slug, err := article.GenerateSlugFromTitle(uriIdentifier)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("uriIdentifier : %w", err))
		return
	}
	art.Slug = slug

	registeredArticle, err := a.repository.GetOneBySlug(slug)
	if err != nil {
		if errors.As(err, fmt.Errorf("article not found")) {
			art.ID = uuid.New()
			err = a.repository.Create(*art)
			if err != nil {
				a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to create article: %w", err))
				return
			}
			newArticle, err := a.repository.GetOneById(art.ID)
			if err != nil {
				a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to retrieve newly created article: %w", err))
				return
			}
			err = a.helpers.respondWithJSON(w, http.StatusCreated, newArticle)
			if err != nil {
				a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to write response : %w", err))
				return
			}
			return
		}
	}

	art.ID = registeredArticle.ID
	err = a.repository.Update(*art)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("update product failed : %w", err))
		return
	}
	newArticle, err := a.repository.GetOneById(art.ID)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to retrieve newly created article: %w", err))
		return
	}
	err = a.helpers.respondWithJSON(w, http.StatusOK, newArticle)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to write response : %w", err))
		return
	}
}

func (a *articleHTTPHandler) updateArticleByID(w http.ResponseWriter, r *http.Request, art *article.Article) {
	_, err := a.repository.GetOneById(art.ID)
	if err != nil {
		if errors.As(err, errors.New("article not found")) {
			http.NotFound(w, r)
			return
		}
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to load article : %w", err))
		return
	}
	slug, err := article.GenerateSlugFromTitle(art.Title)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to generate slug : %w", err))
		return
	}
	art.Slug = slug
	err = a.repository.Update(*art)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to update article : %w", err))
		return
	}
	err = a.helpers.respondWithJSON(w, http.StatusOK, art)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to write response : %w", err))
		return
	}
	return
}
