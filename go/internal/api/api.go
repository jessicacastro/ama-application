package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/jessicacastro/ama-application/go/internal/store/pgstore"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type apiHandler struct {
	query       *pgstore.Queries
	router      *chi.Mux
	upgrader    websocket.Upgrader
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	mu          *sync.Mutex
}

func (apiHandler apiHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	apiHandler.router.ServeHTTP(responseWriter, request)
}

func NewAPIHandler(pgStoreQuery *pgstore.Queries) http.Handler {
	apiHandler := &apiHandler{
		query: pgStoreQuery,
	}

	chiRouter := chi.NewRouter()

	apiHandler.router = chiRouter

	return apiHandler
}
