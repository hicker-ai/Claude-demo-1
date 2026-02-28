package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qinzj/claude-demo/internal/service"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authSvc}
}

// LoginReq is the login request DTO.
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user and returns a JWT token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	token, u, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		Error(c, http.StatusUnauthorized, "authentication failed")
		return
	}

	OK(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":           u.ID,
			"username":     u.Username,
			"display_name": u.DisplayName,
		},
	})
}

// Logout handles user logout (stateless JWT - just returns success).
func (h *AuthHandler) Logout(c *gin.Context) {
	OK(c, nil)
}
