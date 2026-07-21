package main

import (
	"bytes"
	"image/png"
	"testing"
)

func TestEmbeddedTrayIcons(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data []byte
	}{
		{name: "colour", data: deepSeekColorIcon},
		{name: "template", data: deepSeekTemplateIcon},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			config, err := png.DecodeConfig(bytes.NewReader(test.data))
			if err != nil {
				t.Fatalf("decode embedded icon: %v", err)
			}
			if config.Width != 640 || config.Height != 640 {
				t.Fatalf("icon dimensions = %dx%d, want 640x640", config.Width, config.Height)
			}
		})
	}
}
