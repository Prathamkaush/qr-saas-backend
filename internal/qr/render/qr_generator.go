package render

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/skip2/go-qrcode"
)

type RenderOptions struct {
	Size            int    // px, e.g. 512
	Color           string // Hex code e.g. "#FF0000"
	BackgroundColor string // Hex code e.g. "#FFFFFF"
	LogoPath        string // local file or fetched and cached
}

// RenderQRWithLogo generates a QR image bytes with optional logo
func RenderQRWithLogo(content string, opts RenderOptions) ([]byte, error) {
	if opts.Size == 0 {
		opts.Size = 512
	}

	// Generate QR base
	qrImg, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrImg.DisableBorder = true

	// Apply Custom Colors if provided
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

// Helper: Parse Hex string to color.RGBA
func parseHexColor(s string) color.RGBA {
	c := color.RGBA{A: 0xff}
	var r, g, b uint8

	if len(s) == 7 {
		// Handle #RRGGBB
		fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
		c.R = r
		c.G = g
		c.B = b
	} else if len(s) == 4 {
		// Handle #RGB (shorthand)
		var r1, g1, b1 uint8
		fmt.Sscanf(s, "#%1x%1x%1x", &r1, &g1, &b1)
		c.R = r1 * 17
		c.G = g1 * 17
		c.B = b1 * 17
	}
	
	return c
}
