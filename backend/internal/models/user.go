package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a VPN user
type User struct {
	ID          string    `json:"id" bson:"id"`
	Username    string    `json:"username" bson:"username"`
	Email       string    `json:"email" bson:"email"`
	Status      string    `json:"status" bson:"status"` // active, inactive, suspended
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	LastLogin   time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`
	PodName     string    `json:"pod_name,omitempty" bson:"pod_name,omitempty"`
	PodIP       string    `json:"pod_ip,omitempty" bson:"pod_ip,omitempty"`
	PublicKey   string    `json:"public_key,omitempty" bson:"public_key,omitempty"`
	PrivateKey  string    `json:"private_key,omitempty" bson:"private_key,omitempty"`
	ConfigData  string    `json:"config_data,omitempty" bson:"config_data,omitempty"`
	DataUsage   int64     `json:"data_usage" bson:"data_usage"` // bytes
	ConnectionCount int   `json:"connection_count" bson:"connection_count"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Status   string `json:"status,omitempty"`
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers       int   `json:"total_users"`
	ActiveUsers      int   `json:"active_users"`
	InactiveUsers    int   `json:"inactive_users"`
	SuspendedUsers   int   `json:"suspended_users"`
	TotalDataUsage   int64 `json:"total_data_usage"`
	TotalConnections int   `json:"total_connections"`
}

// NewUser creates a new user with default values
func NewUser(username, email string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     email,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
		DataUsage: 0,
		ConnectionCount: 0,
	}
}

// IsActive returns true if the user is active
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// UpdateLastLogin updates the last login time
func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
	u.UpdatedAt = time.Now()
}

// IncrementConnectionCount increments the connection count
func (u *User) IncrementConnectionCount() {
	u.ConnectionCount++
	u.UpdatedAt = time.Now()
}

// AddDataUsage adds data usage in bytes
func (u *User) AddDataUsage(bytes int64) {
	u.DataUsage += bytes
	u.UpdatedAt = time.Now()
}
