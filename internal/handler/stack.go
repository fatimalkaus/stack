package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/fatimalkaus/stack/internal/stack"
	"github.com/gorilla/mux"
)

// InitRoutes initializes stack routes.
func InitRoutes(r *mux.Router, s stack.Stack) {
	r.Handle("/", Push(s)).Methods("POST")
	r.Handle("/", Pop(s)).Methods("DELETE")
	r.Handle("/", Peek(s)).Methods("GET")
}

// PushRequest represents push request.
type PushRequest struct {
	Data any `json:"data"`
}

// Push pushes data to stack.
func Push(stack stack.Stack) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PushRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			invalidBody(w)
			return
		}

		if err := stack.Push(r.Context(), req.Data); err != nil {
			slog.Error("failed to push", "error", err.Error())
			internal(w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

// PopResponse represents pop response.
type PopResponse struct {
	Data any `json:"data"`
}

// Pop pops data from stack.
func Pop(s stack.Stack) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, err := s.Pop(r.Context())
		if errors.Is(err, stack.ErrIsEmpty) {
			noItems(w)
			return
		}
		if err != nil {
			slog.Error("failed to pop", "error", err.Error())
			internal(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PopResponse{
			Data: val,
		})
	})
}

// Peek peeks data from stack.
func Peek(s stack.Stack) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, err := s.Peek(r.Context())
		if errors.Is(err, stack.ErrIsEmpty) {
			noItems(w)
			return
		}
		if err != nil {
			slog.Error("failed to peek", "error", err.Error())
			internal(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PopResponse{
			Data: val,
		})
	})
}

// errorResponse represents error response.
type errorResponse struct {
	Error string `json:"error"`
}

func noItems(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(errorResponse{
		Error: "no items",
	})
	if err != nil {
		slog.Error("failed to encode", "error", err.Error())
	}
}

func invalidBody(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(errorResponse{
		Error: "invalid body",
	})
	if err != nil {
		slog.Error("failed to encode", "error", err.Error())
	}
}

func internal(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(errorResponse{
		Error: "internal server error",
	})
	if err != nil {
		slog.Error("failed to encode", "error", err.Error())
	}
}
