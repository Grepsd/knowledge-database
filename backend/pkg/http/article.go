package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/article"
	"github.com/grepsd/knowledge-database/pkg/sql"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type articleHTTPHandler struct {
	helpers    *helpers
	repository article.ReadWriteRepositoryer
}

func NewArticleHTTPHandler(s *helpers, r article.ReadWriteRepositoryer) articleHTTPHandler {
	return articleHTTPHandler{helpers: s, repository: r}
}

func (a *articleHTTPHandler) routeCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			a.create(w, r)
			break
		case http.MethodGet:
			a.listArticles(w, r)
			break
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
			break
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (a *articleHTTPHandler) routeItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		re, err := regexp.Compile(`^/articles/([a-f0-9-]+)?/?([a-z]+)?`)
		if err != nil {
			panic("could not compile regexp : " + err.Error())
		}
		matches := re.FindAllStringSubmatch(r.RequestURI, 10)
		fmt.Printf("%#v\n", matches)
		if len(matches[0]) == 3 {
			if matches[0][2] == "tags" {
				if r.Method == http.MethodPost {
					a.assignTagToArticle(w, r.WithContext(context.WithValue(r.Context(), "articleID", matches[0][1])))
				}
			}
			return
		}

		switch r.Method {
		case http.MethodGet:
			a.getArticleById(w, r)
			break
		case http.MethodPut:
			a.putArticle(w, r)
			break
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
			break
		case http.MethodDelete:
			a.delete(w, r)
			break
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (a *articleHTTPHandler) create(w http.ResponseWriter, r *http.Request) {
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

func (a *articleHTTPHandler) listArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := a.repository.GetAll()
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to list articles: %w", err))
		return
	}
	a.helpers.respondWithJSON(w, http.StatusOK, articles)
}

func (a *articleHTTPHandler) getArticleById(w http.ResponseWriter, r *http.Request) {
	var art article.Article
	key := a.helpers.getLastSegmentFromURI(r)
	if id, err := uuid.Parse(key); err == nil && id.String() != "00000000-0000-0000-0000-000000000000" {
		art, err = a.repository.GetOneById(id)
		if err != nil {
			a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to retrieve article by ID: %w", err))
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

func (a *articleHTTPHandler) putArticle(w http.ResponseWriter, r *http.Request) {
	if authorized, ok := r.Context().Value("auth").(bool); ok {
		if !authorized {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
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
		art, err := a.repository.GetOneById(id)
		if err != nil {
			a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to load updated article: %w", err))
			return
		}
		a.helpers.respondWithJSON(w, http.StatusOK, art)
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

func (a *articleHTTPHandler) delete(w http.ResponseWriter, r *http.Request) {
	strID := a.helpers.getLastSegmentFromURI(r)
	id, err := uuid.Parse(strID)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to parse resource ID : %w", err))
		return
	}
	err = a.repository.DeleteById(id)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to delete article : %w", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *articleHTTPHandler) extractsRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to read request body : %w", err))
			return
		}
		defer r.Body.Close()
		var url struct {
			URL string `json "url"`
		}
		err = json.Unmarshal(body, &url)
		if err != nil {
			a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to extract payload : %w", err))
			return
		}

		resp, err := http.Get(url.URL)
		if err != nil {
			a.helpers.writeErrorResponse(w, r, fmt.Errorf("cannot connect to URL : %w", err))
			return
		}
		defer resp.Body.Close()

		title, found := GetHtmlTitle(resp.Body)
		if !found {
			a.helpers.writeErrorResponse(w, r, fmt.Errorf("no title found"))
			return
		}

		output := struct {
			URL   string `json "url"`
			Title string `json "title"`
		}{url.URL, title}
		err = a.helpers.respondWithJSON(w, http.StatusOK, output)
		if err != nil {
			panic("failed to write data to response")
		}
		return
	}
	http.NotFound(w, r)
}

func (a *articleHTTPHandler) assignTagToArticle(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to read payload : %w", err))
		return
	}
	type tt struct {
		ID string `json "ID"`
	}
	var t tt
	err = json.Unmarshal(payload, &t)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to decode payload : %w", err))
		return
	}

	tagID, err := uuid.Parse(t.ID)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to validate tag ID format : %w", err))
		return
	}

	var artID string
	artID, ok := r.Context().Value("articleID").(string)
	if !ok {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to validate article ID : %w", err))
		return
	}

	articleID, err := uuid.Parse(artID)
	if err != nil {
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to validate article ID format : %w", err))
		return
	}

	err = a.repository.AssignTagToArticle(articleID, tagID)
	if err != nil {
		if errors.Is(err, sql.ErrDuplicateKey) {
			a.helpers.respondWithJSON(w, http.StatusConflict, nil)
			return
		}
		a.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to save article to tag link : %w", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
