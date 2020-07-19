package http

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/tag"
	"io/ioutil"
	"net/http"
)

type tagHTTPHandler struct {
	helpers    *helpers
	repository tag.ReadWriteRepositoryer
}

func NewTagHTTPHandler(s *helpers, r tag.ReadWriteRepositoryer) tagHTTPHandler {
	return tagHTTPHandler{helpers: s, repository: r}
}

func (h *tagHTTPHandler) routeCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.getAllTags(w, r)
			break
		case http.MethodPost:
			h.createTag(w, r)
			break
		default:
			http.NotFound(w, r)
		}
	}
}

func (h *tagHTTPHandler) routeItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			http.NotFound(w, r)
		}
	}
}

func (h *tagHTTPHandler) createTag(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		h.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to read payload : %w", err))
		return
	}
	var t tag.Tag
	err = json.Unmarshal(payload, &t)
	if err != nil {
		h.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to decode payload : %w", err))
		return
	}
	t.ID = uuid.New()

	err = h.repository.Create(t)
	if err != nil {
		h.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to save tag : %w", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (h *tagHTTPHandler) getAllTags(w http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()
	hasCategoriesParameter := parameters.Get("has_articles")
	var hasCategoriesFilter bool
	if hasCategoriesParameter != "" {
		if hasCategoriesParameter == "true" {
			hasCategoriesFilter = true
		}
	}

	tags, err := h.repository.GetAll(hasCategoriesFilter)
	if err != nil {
		h.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to load tags : %w", err))
		return
	}
	err = h.helpers.respondWithJSON(w, http.StatusOK, tags)
	if err != nil {
		h.helpers.writeErrorResponse(w, r, fmt.Errorf("failed to write result : %w", err))
		return
	}
}
