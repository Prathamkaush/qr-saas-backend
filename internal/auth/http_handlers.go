package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc: svc}

	r.GET("/google/login", h.GoogleLogin)
	r.GET("/google/callback", h.GoogleCallback)

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}

// GoogleLogin godoc
// @Summary Google OAuth Login
// @Description Redirects user to Google OAuth login page
// @Tags Auth
// @Success 302 {string} string "redirect"
// @Router /api/auth/google/login [get]
func (h *Handler) GoogleLogin(c *gin.Context) {
	state := "dummy-state"
	url := h.svc.GoogleLoginURL(state)
	c.Redirect(http.StatusFound, url)
}

// GoogleCallback godoc
// @Summary Google OAuth Callback
// @Description Handles OAuth callback and returns JWT
// @Tags Auth
// @Produce json
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Router /api/auth/google/callback [get]
func (h *Handler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=no_code")
		return
	}

	token, user, err := h.svc.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		fmt.Println("GoogleCallback ERROR:", err) // 🔥 debug log
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=oauth_failed")
		return
	}

	redirectURL := fmt.Sprintf(
		"http://localhost:3000/auth/callback?token=%s&name=%s",
		token,
		user.Name,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// --- Handlers ---

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
