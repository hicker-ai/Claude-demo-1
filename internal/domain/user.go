package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the account status of a user.
type UserStatus string

const (
	UserStatusEnabled  UserStatus = "enabled"
	UserStatusDisabled UserStatus = "disabled"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	DisplayName  string     `json:"display_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Phone        string     `json:"phone"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Groups       []*Group   `json:"groups,omitempty"`
}

// CreateUserInput holds input for creating a new user.
type CreateUserInput struct {
	Username    string
	DisplayName string
	Email       string
	Password    string
	Phone       string
}

// UpdateUserInput holds input for updating an existing user.
type UpdateUserInput struct {
	DisplayName *string
	Email       *string
	Phone       *string
}

// ListUsersInput holds parameters for listing users.
type ListUsersInput struct {
	Page     int
	PageSize int
	Search   string
}

// ListResult holds a paginated list of items.
type ListResult[T any] struct {
	Items    []*T `json:"items"`
	Total    int  `json:"total"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
}
