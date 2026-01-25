package chat

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pahulgogna/evoAI_Web/backend/auth"
	"github.com/pahulgogna/evoAI_Web/backend/types"
	"github.com/pahulgogna/evoAI_Web/backend/utils"
)

type Handler struct {
	store types.ChatStore
}

func NewHandler(store types.ChatStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/chat/new", auth.WithJWTAuth(h.newChat)).Methods(http.MethodGet)
}

func (h *Handler) newChat(w http.ResponseWriter, r *http.Request) {

	user := auth.GetUserFromContext(r.Context())

	chatId, err := h.store.NewChat(user.Id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, map[string]int32{"chatId": chatId})
}

