package qr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc: svc}

	r.POST("/dynamic/url", h.CreateDynamicURL)
	r.GET("/", h.ListMyQRCodes)
	r.GET("/:id/image", h.GetQRImage)
	r.DELETE("/:id", h.DeleteQR)
}

type CreateDynamicURLRequest struct {
	Name      string      `json:"name"`
    // Removed "url" validation so it accepts WiFi/vCard strings
	TargetURL string      `json:"target_url" binding:"required"` 
    // 🔥 ADDED: Field to capture the type (wifi, vcard, etc.)
	QRType    string      `json:"qr_type" binding:"required"`    
	Design    interface{} `json:"design"`
}

// CreateDynamicURL godoc
// @Summary Create Dynamic QR Code
// @Tags QR
// @Accept json
// @Produce json
// @Router /api/qr/dynamic/url [post]
// @Security BearerAuth
func (h *Handler) CreateDynamicURL(c *gin.Context) {
	var req CreateDynamicURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	userID := c.GetString("user_id") // set by JWT middleware

    // 🔥 FIX: Pass req.QRType as the 5th argument
	qr, err := h.svc.CreateDynamicURL(
		c.Request.Context(),
		userID,
		req.Name,
		req.TargetURL,
		req.QRType, // <--- This was missing!
		req.Design,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create QR: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, qr)
}

// GetQRImage godoc
// @Router /api/qr/{id}/image [get]
func (h *Handler) GetQRImage(c *gin.Context) {
	qrID := c.Param("id")
	userID := c.GetString("user_id")
	scene := c.DefaultQuery("scene", "plain")

	img, err := h.svc.GenerateQRImage(
		c.Request.Context(),
		qrID,
		userID,
		scene,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate image: " + err.Error(),
		})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Writer.Write(img)
}

func (h *Handler) ListMyQRCodes(c *gin.Context) {
	userID := c.GetString("user_id")

	qrs, err := h.svc.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load qr codes"})
		return
	}

	c.JSON(http.StatusOK, qrs)
}

func (h *Handler) DeleteQR(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.Status(http.StatusNoContent)
}
