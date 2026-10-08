package web

import "net/http"

func NewRouter(handler *Handler) *http.ServeMux {
	if handler == nil {
		panic("web.NewRouter: *Handler cannot be nil")
	}
	router := http.NewServeMux()
	router.HandleFunc(connectPattern, handler.Connect)
	return router
}
