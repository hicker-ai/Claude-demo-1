package domain

import (
	"time"

	"github.com/google/uuid"
)

// Group represents a user group in the system.
type Group struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Children    []*Group   `json:"children,omitempty"`
	Users       []*User    `json:"users,omitempty"`
}

// CreateGroupInput holds input for creating a new group.
type CreateGroupInput struct {
	Name        string
	Description string
	ParentID    *uuid.UUID
}

// UpdateGroupInput holds input for updating an existing group.
type UpdateGroupInput struct {
	Name        *string
	Description *string
	ParentID    *uuid.UUID
}
