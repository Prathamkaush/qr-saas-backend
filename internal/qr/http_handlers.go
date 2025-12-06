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
	TargetURL string      `json:"target_url" binding:"required"` 
    
    QRType    string      `json:"qr_type" binding:"required"`
    Design    interface{} `json:"design"`
}

// CreateDynamicURL godoc
// @Summary Create Dynamic QR Code
// @Description Creates a QR code with redirect tracking
// @Tags QR
// @Accept json
// @Produce json
// @Param data body CreateDynamicURLRequest true "QR Data"
// @Success 201 {object} QRCode
// @Failure 400 {object} map[string]string
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

	qr, err := h.svc.CreateDynamicURL(
		c.Request.Context(),
		userID,
		req.Name,
		req.TargetURL,
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
// @Summary Get QR Image
// @Description Returns QR Code Image for given ID
// @Tags QR
// @Produce png
// @Param id path string true "QR Code ID"
// @Param scene query string false "scene (plain, logo, frame)"
// @Success 200 {file} png
// @Failure 404 {object} map[string]string
// @Router /api/qr/{id}/image [get]
// @Security BearerAuth
func (h *Handler) GetQRImage(c *gin.Context) {
	qrID := c.Param("id")
	userID := c.GetString("user_id")
	scene := c.DefaultQuery("scene", "plain") // plain / logo / person_pizza

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
	userID := c.GetString("user_id") // from JWT middleware

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

	// Call service to delete
	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.Status(http.StatusNoContent) // 204 Success
}
