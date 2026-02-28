package domain

import (
	"time"

	"github.com/google/uuid"
)

// Group represents a user group in the system.
type Group struct {
	ID          uuid.UUID
	Name        string
	Description string
	ParentID    *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Children    []*Group
	Users       []*User
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
