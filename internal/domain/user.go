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
	ID           uuid.UUID
	Username     string
	DisplayName  string
	Email        string
	PasswordHash string
	Phone        string
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Groups       []*Group
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
	Items    []*T
	Total    int
	Page     int
	PageSize int
}
