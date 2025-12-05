package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	handlers *HTTPHeandlers
	server   *http.Server
}

func NewHTTPServer(handlers *HTTPHeandlers) *HTTPServer {
	return &HTTPServer{
		handlers: handlers,
	}
}

func (h *HTTPServer) Run() error {
	router := mux.NewRouter()

	router.Path("/miners").Methods(http.MethodPost).HandlerFunc(h.handlers.HandleCreateNewMiner)
	router.Path("/miners").Methods(http.MethodGet).HandlerFunc(h.handlers.HandlerGetMiner)
	router.Path("/miners/salaries").Methods(http.MethodGet).HandlerFunc(h.handlers.HandleGetMinerSalaries)
	router.Path("/equipment").Methods(http.MethodPost).HandlerFunc(h.handlers.HandleByeEquipment)
	router.Path("/equipment").Methods(http.MethodGet).HandlerFunc(h.handlers.HandleCheckEquipment)
	router.Path("/equipment/price").Methods(http.MethodGet).HandlerFunc(h.handlers.HandleGetEquipmentPrice)
	router.Path("/company").Methods(http.MethodGet).HandlerFunc(h.handlers.HandleGetCompanyStatistics)
	router.HandleFunc("/company/complete", h.handlers.HandleCompleateGame).Methods("POST")

	server := http.Server{
		Addr:    ":9091",
		Handler: router,
	}
	h.server = &server

	h.handlers.SetCloseServerFunc(server.Close)

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return nil
}
