package chat

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
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
	router.HandleFunc("/chat/message", auth.WithJWTAuth(h.storeMessage)).Methods(http.MethodPost)
	router.HandleFunc("/chat/all", auth.WithJWTAuth(h.getAllChats)).Methods(http.MethodGet)
	router.HandleFunc("/chat/{id}", auth.WithJWTAuth(h.getAllChatMessages)).Methods(http.MethodGet)
	router.HandleFunc("/chat/m/{messageId}", auth.WithJWTAuth(h.deleteMessage)).Methods(http.MethodDelete)
	router.HandleFunc("/chat/{chatId}", auth.WithJWTAuth(h.deleteChat)).Methods(http.MethodDelete)
}

func (h *Handler) storeMessage(w http.ResponseWriter, r *http.Request) {

	user := auth.GetUserFromContext(r.Context())

	var payload types.StoreMessage
	if err := utils.ParseBodyJSON(r, &payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	if _, ok := types.ValidMessageSenders[payload.Role]; !ok {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid message sender"))
		return
	}

	if payload.CreateNewChat {
		createdChatAndMessage, err := h.store.NewChatWithMessage(user.Id, &payload)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		utils.WriteJsonResponse(w, http.StatusCreated, createdChatAndMessage)
		return
	}

	messageId, err := h.store.StoreMessage(user.Id, &payload)
	if err != nil {
		if err.Error() == "chat not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJsonResponse(w, http.StatusCreated, types.NewChatWithMessageResponse{
		ChatId:    payload.ChatId,
		MessageId: messageId,
	})
}

func (h *Handler) getAllChats(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())

	chats, err := h.store.GetAllChats(user.Id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, chats)
}

func (h *Handler) getAllChatMessages(w http.ResponseWriter, r *http.Request) {

	user := auth.GetUserFromContext(r.Context())

	wildcards := mux.Vars(r)
	chatId := wildcards["id"]
	if chatId == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid chat id"))
		return
	}

	chats, err := h.store.GetAllChatMessages(chatId, user.Id)
	if err != nil {
		if err.Error() == "chat not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, chats)
}

func (h *Handler) deleteMessage(w http.ResponseWriter, r *http.Request) {
	
	wildcards := mux.Vars(r)

	messageId := wildcards["messageId"]

	if messageId == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("messageId not present"))
		return
	}

	user := auth.GetUserFromContext(r.Context())

	err := h.store.DeleteMessage(user.Id, messageId)
	if err != nil {
		if err.Error() == "message not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJsonResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) deleteChat(w http.ResponseWriter, r *http.Request) {

	wildcards := mux.Vars(r)

	chatId := wildcards["chatId"]

	if chatId == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("chatId not found"))
		return
	}

	user := auth.GetUserFromContext(r.Context())

	err := h.store.DeleteChat(user.Id, chatId)
	if err != nil {
		if err.Error() == "chat not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJsonResponse(w, http.StatusNoContent, nil)
}
