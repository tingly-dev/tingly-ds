package main

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type namedWindow struct {
	application.Window
	name string
}

func (w namedWindow) Name() string { return w.name }

func TestIsDeepSeekURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "chat", url: "https://chat.deepseek.com/", want: true},
		{name: "root", url: "https://deepseek.com", want: true},
		{name: "nested subdomain", url: "https://auth.eu.deepseek.com/login", want: true},
		{name: "explicit TLS port", url: "https://chat.deepseek.com:443/", want: true},
		{name: "trailing dot", url: "https://chat.deepseek.com./", want: true},
		{name: "plain HTTP", url: "http://chat.deepseek.com/", want: false},
		{name: "lookalike suffix", url: "https://deepseek.com.example.org/", want: false},
		{name: "lookalike prefix", url: "https://notdeepseek.com/", want: false},
		{name: "embedded credentials", url: "https://user@chat.deepseek.com/", want: false},
		{name: "unexpected port", url: "https://chat.deepseek.com:8443/", want: false},
		{name: "relative", url: "/chat", want: false},
		{name: "malformed", url: "://", want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := isDeepSeekURL(test.url); got != test.want {
				t.Fatalf("isDeepSeekURL(%q) = %v, want %v", test.url, got, test.want)
			}
		})
	}
}

func TestExternalURLFromMessage(t *testing.T) {
	t.Parallel()

	trustedOrigin := &application.OriginInfo{
		Origin:      "https://chat.deepseek.com/chat/s/123",
		IsMainFrame: true,
	}
	window := namedWindow{name: "deepseek"}

	tests := []struct {
		name    string
		window  application.Window
		message string
		origin  *application.OriginInfo
		want    string
		wantErr bool
	}{
		{
			name:    "external HTTPS URL",
			window:  window,
			message: `{"type":"open-external","url":"https://example.com/docs?q=1"}`,
			origin:  trustedOrigin,
			want:    "https://example.com/docs?q=1",
		},
		{
			name:    "external HTTP URL",
			window:  window,
			message: `{"type":"open-external","url":"http://example.com/"}`,
			origin:  trustedOrigin,
			want:    "http://example.com/",
		},
		{name: "wrong window", window: namedWindow{name: "settings"}, message: `{}`, origin: trustedOrigin, wantErr: true},
		{name: "nil origin", window: window, message: `{}`, origin: nil, wantErr: true},
		{
			name:    "subframe",
			window:  window,
			message: `{"type":"open-external","url":"https://example.com"}`,
			origin:  &application.OriginInfo{Origin: "https://chat.deepseek.com/", IsMainFrame: false},
			wantErr: true,
		},
		{
			name:    "lookalike sender",
			window:  window,
			message: `{"type":"open-external","url":"https://example.com"}`,
			origin:  &application.OriginInfo{Origin: "https://deepseek.com.evil.test/", IsMainFrame: true},
			wantErr: true,
		},
		{name: "malformed JSON", window: window, message: `{`, origin: trustedOrigin, wantErr: true},
		{name: "empty message", window: window, message: ``, origin: trustedOrigin, wantErr: true},
		{
			name:    "multiple JSON values",
			window:  window,
			message: `{"type":"open-external","url":"https://example.com"} {}`,
			origin:  trustedOrigin,
			wantErr: true,
		},
		{
			name:    "unknown field",
			window:  window,
			message: `{"type":"open-external","url":"https://example.com","extra":true}`,
			origin:  trustedOrigin,
			wantErr: true,
		},
		{
			name:    "unknown action",
			window:  window,
			message: `{"type":"quit","url":"https://example.com"}`,
			origin:  trustedOrigin,
			wantErr: true,
		},
		{
			name:    "DeepSeek stays embedded",
			window:  window,
			message: `{"type":"open-external","url":"https://api.deepseek.com/"}`,
			origin:  trustedOrigin,
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			window:  window,
			message: `{"type":"open-external","url":"ftp://example.com/file"}`,
			origin:  trustedOrigin,
			wantErr: true,
		},
		{
			name:    "credentials",
			window:  window,
			message: `{"type":"open-external","url":"https://user:secret@example.com/"}`,
			origin:  trustedOrigin,
			wantErr: true,
		},
		{
			name:    "oversized",
			window:  window,
			message: strings.Repeat("x", maxBridgeMessageBytes+1),
			origin:  trustedOrigin,
			wantErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := externalURLFromMessage(test.window, test.message, test.origin)
			if test.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got URL %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
