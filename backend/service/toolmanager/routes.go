package toolmanager

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pahulgogna/evoAI_Web/backend/auth"
	"github.com/pahulgogna/evoAI_Web/backend/types"
	"github.com/pahulgogna/evoAI_Web/backend/utils"
)

type Handler struct {
	store types.ToolStore
}

func NewHandler(store types.ToolStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tools", auth.WithJWTAuth(h.createTool)).Methods("POST")
	router.HandleFunc("/tools", auth.WithJWTAuth(h.getTools)).Methods("GET")
	router.HandleFunc("/tools/{id}", auth.WithJWTAuth(h.getTool)).Methods("GET")
	router.HandleFunc("/tools/{id}", auth.WithJWTAuth(h.updateTool)).Methods("PATCH")
	router.HandleFunc("/tools/{id}", auth.WithJWTAuth(h.deleteTool)).Methods("DELETE")
}

func (h *Handler) createTool(w http.ResponseWriter, r *http.Request) {

	var payload types.CreateTool
	if err := utils.ParseBodyJSON(r, &payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	user := auth.GetUserFromContext(r.Context())

	err := h.store.CreateToolByUser(&payload, user.Id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusCreated, map[string]string{"message": "tool created successfully"})
}

func (h *Handler) getTools(w http.ResponseWriter, r *http.Request) {

	user := auth.GetUserFromContext(r.Context())

	tools, err := h.store.GetToolsByUserId(user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Errorf("no tools found"))
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, tools)
}

func (h *Handler) getTool(w http.ResponseWriter, r *http.Request) {

	wildcards := mux.Vars(r)
	toolId := wildcards["id"]

	user := auth.GetUserFromContext(r.Context())

	tool, err := h.store.GetToolByIdAndUserId(toolId, user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Errorf("tool not found"))
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, tool)
}

func (h *Handler) deleteTool(w http.ResponseWriter, r *http.Request) {
	wildcards := mux.Vars(r)
	toolId := wildcards["id"]

	user := auth.GetUserFromContext(r.Context())

	err := h.store.DeleteTool(toolId, user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Errorf("tool not found"))
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, nil)
}

func (h *Handler) updateTool(w http.ResponseWriter, r *http.Request) {
	wildcards := mux.Vars(r)
	toolId := wildcards["id"]

	user := auth.GetUserFromContext(r.Context())

	var payload types.UpdateTool
	if err := utils.ParseBodyJSON(r, &payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.store.UpdateTool(&payload, toolId, user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Errorf("tool not found"))
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, map[string]string{"message": "tool updated successfully"})
}
