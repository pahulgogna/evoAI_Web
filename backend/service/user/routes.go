package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pahulgogna/evoAI_Web/backend/auth"
	"github.com/pahulgogna/evoAI_Web/backend/config"
	"github.com/pahulgogna/evoAI_Web/backend/types"
	"github.com/pahulgogna/evoAI_Web/backend/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(userStore types.UserStore) *Handler {
	return &Handler{
		store: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/{id}", h.getUser).Methods(http.MethodGet)
	router.HandleFunc("/user/register", h.createNewUser).Methods(http.MethodPost)
	router.HandleFunc("/user{id}", h.removeUser).Methods(http.MethodDelete)
	router.HandleFunc("/user/login", h.loginUser).Methods(http.MethodPost)
}

// TODO: protect
func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	if userId == "" {
		utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Errorf("id not found"))
		return
	}

	user, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, types.ResponseUser{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	})
}

// TODO: protect
func (h *Handler) createNewUser(w http.ResponseWriter, r *http.Request) {

	var payload types.RegisterUser
	if err := utils.ParseBodyJSON(r, &payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}
	log.Println(err.Error())

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.store.CreateUser(&types.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusCreated, nil)
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUser
	if err := utils.ParseBodyJSON(r, &payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid credentials"))
		return
	}

	if err := auth.ComparePasswords(user.Password, payload.Password); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid credentials"))
		return
	}

	secret := config.Envs.JWTSecret

	tokenString, err := auth.CreateJWT([]byte(secret), user.ID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, map[string]string{"token": tokenString})
}

// TODO: protect
func (h *Handler) removeUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	if userId == "" {
		utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Errorf("id not found"))
		return
	}

	err := h.store.RemoveUser(userId)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJsonResponse(w, http.StatusNoContent, nil)
}
