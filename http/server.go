package http

import (
	"github.com/gorilla/mux"
)

func NewRouter(h *HTTPHandlers) *mux.Router {
	router := mux.NewRouter()

	router.Path("/company").Methods("GET").HandlerFunc(h.HandleGetCompanyInfo)
	router.Path("/company").Methods("POST").HandlerFunc(h.HandleFinishGame)
	router.Path("/miners").Methods("GET").HandlerFunc(h.HandleGetMiners)
	router.Path("/miners/{class}").Methods("GET").HandlerFunc(h.HandleGetMinersByClass)
	router.Path("/miners").Methods("POST").HandlerFunc(h.HandleHireMiner)
	router.Path("/equipment").Methods("GET").HandlerFunc(h.HandleGetEquipment)
	router.Path("/equipment/{name}").Methods("POST").HandlerFunc(h.HandleBuyEquipment)

	return router
}
