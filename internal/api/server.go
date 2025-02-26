package api

import (
	db "payment-system/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	Router *gin.Engine
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Публичные маршруты
	router.POST("/api/auth", server.handleLogin)
	router.GET("/info", server.handleGetInfo)
	router.GET("/buy/:item", server.handleBuyItem)

	server.Router = router
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
