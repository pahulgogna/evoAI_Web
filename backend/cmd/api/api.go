package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pahulgogna/evoAI_Web/backend/service/chat"
	"github.com/pahulgogna/evoAI_Web/backend/service/toolmanager"
	"github.com/pahulgogna/evoAI_Web/backend/service/user"
)

type ApiServer struct {
	addr string
	db   *sqlx.DB
}

func NewServer(addr string, db *sqlx.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	toolStore := toolmanager.NewStore(s.db)
	toolsHandler := toolmanager.NewHandler(toolStore)
	toolsHandler.RegisterRoutes(subrouter)

	chatStore := chat.NewStore(s.db)
	chatHandler := chat.NewHandler(chatStore)
	chatHandler.RegisterRoutes(subrouter)

	log.Println("server started on:", s.addr)
	return http.ListenAndServe(s.addr, router)
}
