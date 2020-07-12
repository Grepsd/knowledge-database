package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type helpers struct {
}

func NewHelpers() *helpers {
	return &helpers{}
}

func (h *helpers) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	h.writeJsonContentType(w)
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}
	err := h.writeJson(w, data)
	return err
}

func (h *helpers) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Print(err)
}

func (h *helpers) writeResponse(w http.ResponseWriter, data string) error {
	_, err := fmt.Fprint(w, data)
	if err != nil {
		err = fmt.Errorf("failed to write to responseWriter : %w", err)
	}
	return err
}

func (h *helpers) writeJsonContentType(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func (h *helpers) getLastSegmentFromURI(r *http.Request) string {
	uri := r.RequestURI
	parts := strings.Split(uri, "/")
	return parts[len(parts)-1]
}

func (h *helpers) writeJson(w http.ResponseWriter, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshall: %w", err)
	}
	h.writeJsonContentType(w)
	_, err = w.Write(payload)
	if err != nil {
		panic("failed to write response: " + err.Error())
	}
	return err
}

func (h *helpers) isUserAuthenticated(r *http.Request) bool {
	return r.Context().Value("auth").(bool)
}