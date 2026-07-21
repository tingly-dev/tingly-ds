package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const maxBridgeMessageBytes = 4096

// externalLinkBridgeJS keeps DeepSeek navigation in the WebView and sends
// unrelated HTTP(S) anchor clicks to the narrowly scoped Go raw-message
// handler. It is a UX aid rather than a security boundary.
const externalLinkBridgeJS = `(function () {
  if (window.__deepSeekShellExternalLinksInstalled) return;
  window.__deepSeekShellExternalLinksInstalled = true;

	  function isDeepSeekURL(target) {
	    var host = String(target.hostname || "").toLowerCase().replace(/\.$/, "");
	    var isDefaultTLS = target.port === "" || target.port === "443";
	    return target.protocol === "https:" && isDefaultTLS &&
	      (host === "deepseek.com" || host.endsWith(".deepseek.com"));
  }

  document.addEventListener("click", function (event) {
    if (event.defaultPrevented || event.button !== 0) return;
    if (!(event.target instanceof Element)) return;

    var anchor = event.target.closest("a[href]");
    if (!anchor) return;

    var target;
    try {
      target = new URL(anchor.href, window.location.href);
    } catch (_) {
      return;
    }

    if (target.protocol !== "http:" && target.protocol !== "https:") return;
	    if (isDeepSeekURL(target)) {
      if (String(anchor.target || "").toLowerCase() === "_blank") {
        event.preventDefault();
        window.location.assign(target.href);
      }
      return;
    }

    if (!window._wails || typeof window._wails.invoke !== "function") return;
    event.preventDefault();
    window._wails.invoke(JSON.stringify({
      type: "open-external",
      url: target.href
    }));
  }, true);
})();`

type externalLinkRequest struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

func externalURLFromMessage(window application.Window, message string, origin *application.OriginInfo) (string, error) {
	if window == nil || window.Name() != "deepseek" {
		return "", errors.New("unexpected source window")
	}
	if origin == nil || !origin.IsMainFrame || !isDeepSeekURL(origin.Origin) {
		return "", errors.New("untrusted message origin")
	}
	if len(message) == 0 || len(message) > maxBridgeMessageBytes {
		return "", errors.New("invalid message size")
	}

	var request externalLinkRequest
	decoder := json.NewDecoder(strings.NewReader(message))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return "", fmt.Errorf("decode message: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return "", errors.New("multiple message values")
	}
	if request.Type != "open-external" {
		return "", errors.New("unknown message type")
	}

	target, err := url.Parse(request.URL)
	if err != nil || target.Hostname() == "" {
		return "", errors.New("invalid target URL")
	}
	if target.Scheme != "http" && target.Scheme != "https" {
		return "", errors.New("unsupported target scheme")
	}
	if target.User != nil {
		return "", errors.New("target credentials are not allowed")
	}
	if isDeepSeekURL(target.String()) {
		return "", errors.New("DeepSeek targets stay embedded")
	}

	return target.String(), nil
}

func isDeepSeekURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil {
		return false
	}
	if port := parsed.Port(); port != "" && port != "443" {
		return false
	}
	host := strings.ToLower(strings.TrimSuffix(parsed.Hostname(), "."))
	return host == "deepseek.com" || strings.HasSuffix(host, ".deepseek.com")
}
