package render

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/skip2/go-qrcode"
)

type RenderOptions struct {
	Size          int    // px, e.g. 512
	ForegroundHex string // you can parse to color.RGBA
	BackgroundHex string
	LogoPath      string // local file or fetched and cached
}

// Generate QR image bytes with optional logo in center
func RenderQRWithLogo(content string, opts RenderOptions) ([]byte, error) {
	if opts.Size == 0 {
		opts.Size = 512
	}

	// generate QR base
	// NOTE: go-qrcode supports custom colors but default is fine to start
	qrImg, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrImg.DisableBorder = false
	qrPNG := qrImg.Image(opts.Size)

	base := imaging.Clone(qrPNG)

	// overlay logo if provided
	if opts.LogoPath != "" {
		logoFile, err := os.Open(opts.LogoPath)
		if err == nil {
			defer logoFile.Close()
			logo, _, err := image.Decode(logoFile)
			if err == nil {
				// resize logo to e.g. 20% of QR size
				logoSize := opts.Size / 4
				logo = imaging.Resize(logo, logoSize, logoSize, imaging.Lanczos)

				// center position
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
