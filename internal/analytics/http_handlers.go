// internal/analytics/http_handlers.go
package analytics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc: svc}

	// GET /api/analytics/:qrID/summary?from=2025-01-01&to=2025-01-31
	r.GET("/:qrID/summary", h.GetSummary)

	// GET /api/analytics/:qrID/timeseries?from=...&to=...&granularity=day
	r.GET("/:qrID/timeseries", h.GetTimeSeries)

	r.GET("/dashboard", h.GetDashboardStats)

	r.GET("/dashboard/timeseries", h.GetGlobalTimeSeries)
}

// parseDateRange gives default last 7 days if not provided
func parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr == "" {
		from = now.AddDate(0, 0, -7)
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if toStr == "" {
		to = now
	} else {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	return from, to, nil
}

// GetSummary godoc
// @Summary Get summary analytics for a QR code
// @Description Returns scan count, countries, devices, browsers
// @Tags Analytics
// @Produce json
// @Param qrID path string true "QR Code ID"
// @Param from query string false "From Date YYYY-MM-DD"
// @Param to query string false "To Date YYYY-MM-DD"
// @Security BearerAuth
// @Success 200 {object} SummaryResponse
// @Router /api/analytics/{qrID}/summary [get]
func (h *Handler) GetSummary(c *gin.Context) {
	qrID := c.Param("qrID")
	userID := c.GetString("user_id")

	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	// Call the Service (which calls the Repo we just fixed)
	summaryData, err := h.svc.GetSummary(c.Request.Context(), userID, qrID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map Repository struct to API Response struct
	response := SummaryResponse{
		TotalScans: int(summaryData.TotalScans),
		Countries:  summaryData.Countries,
		Devices:    summaryData.Devices,
		Browsers:   summaryData.Browsers,
	}

	c.JSON(http.StatusOK, response)
}

// GetTimeSeries godoc
// @Summary Get time-series analytics
// @Description Returns daily/hourly scan counts for a QR code
// @Tags Analytics
// @Produce json
// @Param qrID path string true "QR Code ID"
// @Param from query string false "From Date"
// @Param to query string false "To Date"
// @Param granularity query string false "day|hour" default(day)
// @Security BearerAuth
// @Success 200 {array} TimePoint
// @Router /api/analytics/{qrID}/timeseries [get]
func (h *Handler) GetTimeSeries(c *gin.Context) {
	qrID := c.Param("qrID")
	userID := c.GetString("user_id")

	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	granularity := c.DefaultQuery("granularity", "day")

	points, err := h.svc.GetTimeSeries(c.Request.Context(), userID, qrID, from, to, granularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Optional: limit points via ?limit=100
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit < len(points) {
			points = points[:limit]
		}
	}

	c.JSON(http.StatusOK, points)
}

// Handler
func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID := c.GetString("user_id")
	summary, err := h.svc.GetGlobalStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, summary)
}

func (h *Handler) GetGlobalTimeSeries(c *gin.Context) {
	userID := c.GetString("user_id")
	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date range"})
		return
	}

	points, err := h.svc.GetGlobalTimeSeries(c.Request.Context(), userID, from, to, "day")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, points)
}
