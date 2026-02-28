package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qinzj/claude-demo/internal/config"
)

// LDAPConfigHandler handles LDAP configuration endpoints.
type LDAPConfigHandler struct {
	cfg *config.LDAPConfig
}

// NewLDAPConfigHandler creates a new LDAPConfigHandler.
func NewLDAPConfigHandler(cfg *config.LDAPConfig) *LDAPConfigHandler {
	return &LDAPConfigHandler{cfg: cfg}
}

// GetConfig returns the current LDAP configuration.
func (h *LDAPConfigHandler) GetConfig(c *gin.Context) {
	OK(c, gin.H{
		"base_dn": h.cfg.BaseDN,
		"mode":    h.cfg.Mode,
		"port":    h.cfg.Port,
	})
}

// UpdateConfigReq is the request DTO for updating LDAP config.
type UpdateConfigReq struct {
	BaseDN string `json:"base_dn" binding:"required"`
	Mode   string `json:"mode" binding:"required,oneof=openldap activedirectory"`
	Port   int    `json:"port" binding:"required,min=1,max=65535"`
}

// UpdateConfig updates the LDAP configuration.
func (h *LDAPConfigHandler) UpdateConfig(c *gin.Context) {
	var req UpdateConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	h.cfg.BaseDN = req.BaseDN
	h.cfg.Mode = req.Mode
	h.cfg.Port = req.Port
	OK(c, nil)
}

// GetStatus returns the LDAP server status.
func (h *LDAPConfigHandler) GetStatus(c *gin.Context) {
	OK(c, gin.H{
		"running": true,
		"port":    h.cfg.Port,
		"mode":    h.cfg.Mode,
	})
}
