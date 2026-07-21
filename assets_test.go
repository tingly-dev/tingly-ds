package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestEmbeddedIcons(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		data       []byte
		wantWidth  int
		wantHeight int
	}{
		{name: "dock", data: deepSeekDockIcon, wantWidth: 1024, wantHeight: 1024},
		{name: "colour", data: deepSeekColorIcon, wantWidth: 640, wantHeight: 640},
		{name: "template", data: deepSeekTemplateIcon, wantWidth: 640, wantHeight: 640},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			config, err := png.DecodeConfig(bytes.NewReader(test.data))
			if err != nil {
				t.Fatalf("decode embedded icon: %v", err)
			}
			if config.Width != test.wantWidth || config.Height != test.wantHeight {
				t.Fatalf("icon dimensions = %dx%d, want %dx%d", config.Width, config.Height, test.wantWidth, test.wantHeight)
			}
		})
	}
}

func TestDockIconHasTileBackground(t *testing.T) {
	t.Parallel()

	icon, err := png.Decode(bytes.NewReader(deepSeekDockIcon))
	if err != nil {
		t.Fatalf("decode Dock icon: %v", err)
	}

	if alphaAt(icon, 0, 0) != 0 {
		t.Fatal("Dock icon corner is not transparent")
	}
	if alphaAt(icon, icon.Bounds().Dx()/2, icon.Bounds().Dy()/2) != 0xffff {
		t.Fatal("Dock icon tile centre is not opaque")
	}
}

func alphaAt(img image.Image, x, y int) uint32 {
	_, _, _, alpha := color.NRGBAModel.Convert(img.At(x, y)).RGBA()
	return alpha
}
