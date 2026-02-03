package llm

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pahulgogna/evoAI_Web/backend/auth"
	"github.com/pahulgogna/evoAI_Web/backend/types"
	"github.com/pahulgogna/evoAI_Web/backend/utils"
)

type Handler struct {
	chatStore    types.ChatStore
	llmInterface types.LLMInterface
}

func NewHandler(chatStore types.ChatStore, llmInterface types.LLMInterface) *Handler {
	return &Handler{
		chatStore:    chatStore,
		llmInterface: llmInterface,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/llm/stream/{model}/{chatId}", auth.WithJWTAuth(h.streamResponse)).Methods(http.MethodGet)
}

func (h *Handler) streamResponse(w http.ResponseWriter, r *http.Request) {

	wildcards := mux.Vars(r)

	chatId := wildcards["chatId"]
	model := wildcards["model"]

	if chatId == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("chatId not found"))
		return
	}

	user := auth.GetUserFromContext(r.Context())

	messages, err := h.chatStore.GetAllChatMessages(chatId, user.Id)
	if err != nil {
		if err.Error() == "chat not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	clientDisconnected := r.Context().Done()

	utils.SetEventStreamHeaders(w)

	responseChannel, err := h.llmInterface.StreamMessage(r.Context(), messages, model)
	if err != nil {
		return
	}

	llmResponse := strings.Builder{}

Loop:
	for {
		select {
		case <-clientDisconnected:
			break Loop
		case chunck, ok := <-responseChannel:
			if !ok {
				break Loop
			}
			err := utils.WriteToEventStream(w, "chunk", chunck)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			llmResponse.WriteString(chunck)
		}
	}

	if llmResponse.Len() == 0 {
		return
	}

	id64, err := strconv.ParseInt(chatId, 10, 16)
	if err != nil {
		// invalid chatId, cannot store message
		return
	}
	_, err = h.chatStore.StoreMessage(user.Id, &types.StoreMessage{
		ChatId:        int16(id64),
		Role:          types.RoleAssistant,
		Content:       llmResponse.String(),
		CreateNewChat: false,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
