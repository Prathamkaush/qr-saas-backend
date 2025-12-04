package render

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)

type CompositeOptions struct {
	BackgroundPath string // e.g. "assets/person_pizza.png"
	QRBytes        []byte // from RenderQRWithLogo
	PosX           int    // X position where QR should start
	PosY           int    // Y position
	Width          int    // width to resize QR into
	Height         int
}

// Returns final composite PNG bytes
func ComposeQROnBackground(opts CompositeOptions) ([]byte, error) {
	bgFile, err := os.Open(opts.BackgroundPath)
	if err != nil {
		return nil, err
	}
	defer bgFile.Close()

	bg, _, err := image.Decode(bgFile)
	if err != nil {
		return nil, err
	}

	qrImg, err := png.Decode(bytes.NewReader(opts.QRBytes))
	if err != nil {
		return nil, err
	}

	// resize QR to fit “plate/box” area
	if opts.Width > 0 && opts.Height > 0 {
		qrImg = imaging.Resize(qrImg, opts.Width, opts.Height, imaging.Lanczos)
	}

	// Overlay QR on background
	composite := imaging.Overlay(bg, qrImg, image.Pt(opts.PosX, opts.PosY), 1.0)

	var out bytes.Buffer
	if err := png.Encode(&out, composite); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
