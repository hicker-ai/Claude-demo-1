package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/qinzj/claude-demo/internal/domain"
	"github.com/qinzj/claude-demo/internal/service"
)

// UserHandler handles user CRUD endpoints.
type UserHandler struct {
	userService  *service.UserService
	groupService *service.GroupService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userSvc *service.UserService, groupSvc *service.GroupService) *UserHandler {
	return &UserHandler{userService: userSvc, groupService: groupSvc}
}

// CreateUserReq is the request DTO for creating a user.
type CreateUserReq struct {
	Username    string `json:"username" binding:"required,max=64"`
	DisplayName string `json:"display_name" binding:"required,max=128"`
	Email       string `json:"email" binding:"required,email,max=255"`
	Password    string `json:"password" binding:"required,min=8"`
	Phone       string `json:"phone,omitempty" binding:"omitempty,max=32"`
}

// UpdateUserReq is the request DTO for updating a user.
type UpdateUserReq struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,max=128"`
	Email       *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
	Phone       *string `json:"phone,omitempty" binding:"omitempty,max=32"`
}

// ChangePasswordReq is the request DTO for changing password.
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// SetStatusReq is the request DTO for setting user status.
type SetStatusReq struct {
	Status string `json:"status" binding:"required,oneof=enabled disabled"`
}

// Create creates a new user.
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	u, err := h.userService.CreateUser(c.Request.Context(), domain.CreateUserInput{
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Password:    req.Password,
		Phone:       req.Phone,
	})
	if err != nil {
		Error(c, http.StatusInternalServerError, "failed to create user: "+err.Error())
		return
	}
	OK(c, u)
}

// Get retrieves a user by ID.
func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	u, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusNotFound, "user not found")
		return
	}
	OK(c, u)
}

// List lists users with pagination and search.
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	result, err := h.userService.ListUsers(c.Request.Context(), domain.ListUsersInput{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})
	if err != nil {
		Error(c, http.StatusInternalServerError, "failed to list users")
		return
	}
	OK(c, result)
}

// Update updates a user.
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	u, err := h.userService.UpdateUser(c.Request.Context(), id, domain.UpdateUserInput{
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Phone:       req.Phone,
	})
	if err != nil {
		Error(c, http.StatusInternalServerError, "failed to update user")
		return
	}
	OK(c, u)
}

// Delete deletes a user.
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		Error(c, http.StatusInternalServerError, "failed to delete user")
		return
	}
	OK(c, nil)
}

// ChangePassword changes a user's password.
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), id, req.OldPassword, req.NewPassword); err != nil {
		Error(c, http.StatusBadRequest, "failed to change password: "+err.Error())
		return
	}
	OK(c, nil)
}

// SetStatus enables or disables a user.
func (h *UserHandler) SetStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req SetStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	if err := h.userService.SetUserStatus(c.Request.Context(), id, domain.UserStatus(req.Status)); err != nil {
		Error(c, http.StatusInternalServerError, "failed to set status")
		return
	}
	OK(c, nil)
}

// GetGroups returns the groups a user belongs to.
func (h *UserHandler) GetGroups(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	groups, err := h.groupService.GetUserGroups(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusInternalServerError, "failed to get user groups")
		return
	}
	OK(c, groups)
}
