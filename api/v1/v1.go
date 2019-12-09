package v1

import (
	"net/http"

	"github.com/VolticFroogo/cryptopad-server/api/v1/pad"
	"github.com/gorilla/mux"
)

const (
	urlPrefix = "/api/v1/"
)

// Handle adds the v1 API endpoints.
func Handle(r *mux.Router) {
	r.Handle(urlPrefix+"pad", http.HandlerFunc(pad.Get)).Methods(http.MethodGet)
	r.Handle(urlPrefix+"pad", http.HandlerFunc(pad.Put)).Methods(http.MethodPut)
	r.Handle(urlPrefix+"pad", http.HandlerFunc(pad.Delete)).Methods(http.MethodDelete)
}
