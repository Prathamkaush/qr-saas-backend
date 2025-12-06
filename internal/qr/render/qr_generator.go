package render

import (
	"bytes"
	"fmt"
	"image"
	"image/color" // Needed for color types
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/skip2/go-qrcode"
)

type RenderOptions struct {
	Size            int    // px, e.g. 512
	Color           string // Hex code e.g. "#FF0000" (Renamed from ForegroundHex)
	BackgroundColor string // Hex code e.g. "#FFFFFF" (Renamed from BackgroundHex)
	LogoPath        string // local file or fetched and cached
}

// Generate QR image bytes with optional logo in center
func RenderQRWithLogo(content string, opts RenderOptions) ([]byte, error) {
	if opts.Size == 0 {
		opts.Size = 512
	}

	// Generate QR base
	qrImg, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrImg.DisableBorder = true // Usually looks better for custom designs

	// 🔥 FIX: Apply Custom Colors
	if opts.Color != "" {
		qrImg.ForegroundColor = parseHexColor(opts.Color)
	}
	if opts.BackgroundColor != "" {
		qrImg.BackgroundColor = parseHexColor(opts.BackgroundColor)
	}

	// Create the Image
	qrPNG := qrImg.Image(opts.Size)

	// Convert to NRGBA for imaging library manipulation
	base := imaging.Clone(qrPNG)

	// Overlay logo if provided
	if opts.LogoPath != "" {
		logoFile, err := os.Open(opts.LogoPath)
		if err == nil {
			defer logoFile.Close()
			logo, _, err := image.Decode(logoFile)
			if err == nil {
				// Resize logo to 20% of QR size
				logoSize := opts.Size / 5
				logo = imaging.Resize(logo, logoSize, logoSize, imaging.Lanczos)

				// Center position
				x := (base.Bounds().Dx() - logo.Bounds().Dx()) / 2
				y := (base.Bounds().Dy() - logo.Bounds().Dy()) / 2

				base = imaging.Overlay(base, logo, image.Pt(x, y), 1.0)
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, base); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Helper: Parse Hex string (#RRGGBB) to color.RGBA
func parseHexColor(s string) color.RGBA {
	c := color.RGBA{A: 0xff}
	var r, g, b uint8
	
    // Handle #RRGGBB
	if len(s) == 7 {
		fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
		c.R = r
		c.G = g
		c.B = b
	} 
    // Handle #RGB (shorthand)
    else if len(s) == 4 {
        var r1, g1, b1 uint8
        fmt.Sscanf(s, "#%1x%1x%1x", &r1, &g1, &b1)
        c.R = r1 * 17
        c.G = g1 * 17
        c.B = b1 * 17
    }
    // Default to black if invalid
	return c
}
```

### 2. Update `internal/qr/service.go` (Enable the Fields)

Now that the `render` package supports `Color` and `BackgroundColor`, you need to **uncomment** the lines in your service file that I told you to hide earlier.

Search for `GenerateQRImage` in `internal/qr/service.go` and update this block:

```go
    // ... (inside GenerateQRImage) ...

	// ---------------------------------------------------------
	// 🔥 FIX: Extract Colors from Saved Design JSON
	// ---------------------------------------------------------
	var design DesignConfig
	
	// Default Colors
	fgColor := "#000000"
	bgColor := "#ffffff"

	if qrData.DesignJSON != "" {
		if err := json.Unmarshal([]byte(qrData.DesignJSON), &design); err == nil {
			if design.Color != "" {
				fgColor = design.Color
			}
			if design.BgColor != "" {
				bgColor = design.BgColor
			}
		}
	}

	// Generate QR using the derived content
	qrBytes, err := render.RenderQRWithLogo(contentToEncode, render.RenderOptions{
		Size:            600,
        // 🔥 UNCOMMENT THESE LINES NOW:
		Color:           fgColor, 
		BackgroundColor: bgColor, 
	})
