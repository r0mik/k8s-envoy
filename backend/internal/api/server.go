package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"vpnaas-backend/internal/k8s"
	"vpnaas-backend/internal/metrics"
	"vpnaas-backend/internal/models"
)

// Server represents the API server
type Server struct {
	vpnManager *k8s.VPNManager
	users      map[string]*models.User // In-memory storage for demo
}

// NewServer creates a new API server
func NewServer(vpnManager *k8s.VPNManager) *Server {
	return &Server{
		vpnManager: vpnManager,
		users:      make(map[string]*models.User),
	}
}

// ListUsers returns all users
func (s *Server) ListUsers(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.RecordAPIRequestDuration("GET", "/users", time.Since(start).Seconds())
	}()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	// Update metrics
	s.updateUserMetrics()

	metrics.RecordAPIRequest("GET", "/users", "200")
	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": len(users),
	})
}

// CreateUser creates a new user
func (s *Server) CreateUser(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.RecordAPIRequestDuration("POST", "/users", time.Since(start).Seconds())
	}()

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.RecordAPIRequest("POST", "/users", "400")
		metrics.RecordError("validation", "api")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	for _, user := range s.users {
		if user.Username == req.Username || user.Email == req.Email {
			metrics.RecordAPIRequest("POST", "/users", "409")
			metrics.RecordError("duplicate_user", "api")
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
	}

	// Create new user
	user := models.NewUser(req.Username, req.Email)

	// Create VPN pod
	ctx := context.Background()
	if err := s.vpnManager.CreateUserVPN(ctx, user); err != nil {
		logrus.Errorf("Failed to create VPN for user %s: %v", user.Username, err)
		metrics.RecordAPIRequest("POST", "/users", "500")
		metrics.RecordError("vpn_creation", "api")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create VPN"})
		return
	}

	// Store user
	s.users[user.ID] = user

	// Update metrics
	s.updateUserMetrics()
	metrics.IncrementConnections()

	metrics.RecordAPIRequest("POST", "/users", "201")
	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"message": "User created successfully",
	})
}

// GetUser returns a specific user
func (s *Server) GetUser(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.RecordAPIRequestDuration("GET", "/users/:id", time.Since(start).Seconds())
	}()

	userID := c.Param("id")
	user, exists := s.users[userID]
	if !exists {
		metrics.RecordAPIRequest("GET", "/users/:id", "404")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get pod status
	if user.PodName != "" {
		ctx := context.Background()
		status, err := s.vpnManager.GetPodStatus(ctx, user.PodName)
		if err == nil {
			user.Status = status
		}
	}

	metrics.RecordAPIRequest("GET", "/users/:id", "200")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// DeleteUser deletes a user
func (s *Server) DeleteUser(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.RecordAPIRequestDuration("DELETE", "/users/:id", time.Since(start).Seconds())
	}()

	userID := c.Param("id")
	user, exists := s.users[userID]
	if !exists {
		metrics.RecordAPIRequest("DELETE", "/users/:id", "404")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete VPN pod
	ctx := context.Background()
	if err := s.vpnManager.DeleteUserVPN(ctx, user); err != nil {
		logrus.Errorf("Failed to delete VPN for user %s: %v", user.Username, err)
		metrics.RecordError("vpn_deletion", "api")
	}

	// Remove user from storage
	delete(s.users, userID)

	// Update metrics
	s.updateUserMetrics()

	metrics.RecordAPIRequest("DELETE", "/users/:id", "200")
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetUserConfig returns the VPN configuration for a user
func (s *Server) GetUserConfig(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.RecordAPIRequestDuration("GET", "/users/:id/config", time.Since(start).Seconds())
	}()

	userID := c.Param("id")
	user, exists := s.users[userID]
	if !exists {
		metrics.RecordAPIRequest("GET", "/users/:id/config", "404")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ConfigData == "" {
		metrics.RecordAPIRequest("GET", "/users/:id/config", "404")
		c.JSON(http.StatusNotFound, gin.H{"error": "VPN configuration not found"})
		return
	}

	// Set headers for file download
	c.Header("Content-Disposition", "attachment; filename=vpn-"+user.Username+".conf")
	c.Header("Content-Type", "text/plain")

	metrics.RecordAPIRequest("GET", "/users/:id/config", "200")
	c.String(http.StatusOK, user.ConfigData)
}

// GetMetrics returns system metrics
func (s *Server) GetMetrics(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.RecordAPIRequestDuration("GET", "/metrics", time.Since(start).Seconds())
	}()

	// Update pod metrics
	ctx := context.Background()
	if err := s.vpnManager.UpdatePodMetrics(ctx); err != nil {
		logrus.Errorf("Failed to update pod metrics: %v", err)
	}

	// Update user metrics
	s.updateUserMetrics()

	metrics.RecordAPIRequest("GET", "/metrics", "200")
	c.JSON(http.StatusOK, gin.H{
		"message": "Metrics updated successfully",
	})
}

// GetStats returns system statistics
func (s *Server) GetStats(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.RecordAPIRequestDuration("GET", "/stats", time.Since(start).Seconds())
	}()

	stats := s.calculateStats()

	metrics.RecordAPIRequest("GET", "/stats", "200")
	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// updateUserMetrics updates user-related metrics
func (s *Server) updateUserMetrics() {
	total := len(s.users)
	active, inactive, suspended := 0, 0, 0

	for _, user := range s.users {
		switch user.Status {
		case "active":
			active++
		case "inactive":
			inactive++
		case "suspended":
			suspended++
		}
	}

	metrics.UpdateUserMetrics(total, active, inactive, suspended)
}

// calculateStats calculates system statistics
func (s *Server) calculateStats() *models.UserStats {
	stats := &models.UserStats{}

	for _, user := range s.users {
		stats.TotalUsers++
		stats.TotalDataUsage += user.DataUsage
		stats.TotalConnections += user.ConnectionCount

		switch user.Status {
		case "active":
			stats.ActiveUsers++
		case "inactive":
			stats.InactiveUsers++
		case "suspended":
			stats.SuspendedUsers++
		}
	}

	return stats
}
