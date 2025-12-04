package redirect

import (
	"context"
	"errors"
	"fmt"
	"time"

	"qr-saas/internal/analytics"
	"qr-saas/internal/qr"

	"github.com/google/uuid"        // Use UUIDs for unique events
	"github.com/mileusna/useragent" // <--- NEW: Import this
)

var ErrNotFound = errors.New("qr code not found or inactive")

type Service struct {
	qrRepo    qr.Repository
	analytics analytics.Service
}

func NewService(qrRepo qr.Repository, analyticsSvc analytics.Service) *Service {
	return &Service{
		qrRepo:    qrRepo,
		analytics: analyticsSvc,
	}
}

func (s *Service) ResolveAndLog(ctx context.Context, shortCode string, ip string, uaString string, referer string) (string, error) {

	// 1. Database Lookup (Postgres)
	qrData, err := s.qrRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		fmt.Printf("❌ Database Error: %v\n", err)
		return "", err
	}

	// 2. Validation
	if qrData == nil || !qrData.IsActive {
		fmt.Println("❌ QR not found or inactive")
		return "", ErrNotFound
	}

	fmt.Printf("ℹ️ QR Data Found - ID: %s | Type: '%s'\n", qrData.ID, qrData.QRType)

	// 3. Dynamic Check (Only dynamic QRs track analytics usually)
	targetURL := qrData.TargetURL
	if qrData.QRType == "dynamic" {
		if targetURL == "" {
			return "", ErrNotFound
		}

		// ---------------------------------------------------------
		// DATA ENRICHMENT (Parsing User Agent)
		// ---------------------------------------------------------
		ua := useragent.Parse(uaString)

		deviceType := "Desktop"
		if ua.Mobile {
			deviceType = "Mobile"
		} else if ua.Tablet {
			deviceType = "Tablet"
		} else if ua.Bot {
			deviceType = "Bot"
		}

		// 4. Construct Event
		ev := analytics.ScanEvent{
			EventID:   uuid.NewString(), // Generate a real UUID
			QRID:      qrData.ID,
			UserID:    qrData.UserID, // Important for billing/analytics
			ScannedAt: time.Now().UTC(),
			IP:        ip,
			UserAgent: uaString,
			Referer:   referer,

			// Parsed Data
			DeviceType: deviceType,
			OS:         ua.OS,   // e.g., "Windows 10", "iOS"
			Browser:    ua.Name, // e.g., "Chrome", "Firefox"

			// TODO: GeoIP Lookup (Requires MaxMind DB or external API)
			Country: "Unknown",
			City:    "Unknown",
		}

		// 5. Fire and Forget (Async)
		// We use context.Background() because the HTTP request ctx might cancel
		// before the analytics write finishes.
		go func() {
			fmt.Println("⏳ Attempting to insert into ClickHouse...")
			err := s.analytics.InsertScanEvent(context.Background(), ev)
			if err != nil {
				// 🔥 Print error to terminal so we can see it!
				fmt.Printf("❌ ANALYTICS INSERT ERROR: %v\n", err)
			} else {
				fmt.Printf("✅ Scan Logged for QR: %s\n", qrData.ID)
			}
		}()

		return targetURL, nil
	}

	// For static QR codes, we do not log analytics
	fmt.Println("⚠️ Skipping analytics: QR Type is not 'dynamic'")
	return targetURL, nil
}
