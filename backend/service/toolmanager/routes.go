package toolmanager

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pahulgogna/evoAI_Web/backend/types"
	"github.com/pahulgogna/evoAI_Web/backend/utils"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tools/create", h.createTool).Methods("POST")
	router.HandleFunc("/tools", h.getTools).Methods("GET")
	router.HandleFunc("/tools/*id", h.getTool).Methods("GET")
	router.HandleFunc("/tools/*id", h.updateTool).Methods("PATCH")
	router.HandleFunc("/tools/delete/*id", h.deleteTool).Methods("DELETE")
}

func (h *Handler) createTool(w http.ResponseWriter, r *http.Request) {

	var payload types.Tool

	if err := utils.ParseBodyJSON(r, &payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

}

func (h *Handler) getTools(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getTool(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteTool(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) updateTool(w http.ResponseWriter, r *http.Request) {

}
