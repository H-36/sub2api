package routes

import (
	"bytes"
	"context"
	"errors"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"

	"github.com/gin-gonic/gin"
)

const (
	customPageProxyTimeout      = 15 * time.Second
	customPageProxyMaxBodyBytes = 10 << 20
	customPageProxyUserAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36"
)

var customPageProxyHTTPClient = &http.Client{
	Timeout: customPageProxyTimeout,
	Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           customPageProxyDialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("stopped after 5 redirects")
		}
		if strings.ToLower(req.URL.Scheme) != "https" {
			return errors.New("redirect to non-https url is not allowed")
		}
		return nil
	},
}

func handleCustomPageProxy(settingService *service.SettingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if settingService == nil {
			c.String(http.StatusServiceUnavailable, "settings service unavailable")
			return
		}

		targetURL, ok := resolveProxyableCustomPageURL(c.Request.Context(), settingService, c.Param("id"))
		if !ok {
			c.String(http.StatusNotFound, "custom page not found")
			return
		}

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, targetURL, nil)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid custom page url")
			return
		}
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("User-Agent", customPageProxyUserAgent)
		if acceptLanguage := strings.TrimSpace(c.GetHeader("Accept-Language")); acceptLanguage != "" {
			req.Header.Set("Accept-Language", acceptLanguage)
		}

		resp, err := customPageProxyHTTPClient.Do(req)
		if err != nil {
			c.String(http.StatusBadGateway, "failed to load custom page")
			return
		}
		defer resp.Body.Close()

		body, err := readProxyBody(resp.Body)
		if err != nil {
			c.String(http.StatusBadGateway, "custom page response is too large")
			return
		}

		contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
		if contentType == "" {
			contentType = "text/html; charset=utf-8"
		}
		if isHTMLContentType(contentType) {
			body = injectBaseHref(body, resp.Request.URL.String())
		}

		setCustomPageProxyFrameHeaders(c)
		c.Header("Content-Type", contentType)
		c.Status(resp.StatusCode)
		_, _ = c.Writer.Write(body)
	}
}

func resolveProxyableCustomPageURL(ctx context.Context, settingService *service.SettingService, menuItemID string) (string, bool) {
	menuItemID = strings.TrimSpace(menuItemID)
	if menuItemID == "" {
		return "", false
	}

	settings, err := settingService.GetPublicSettings(ctx)
	if err != nil {
		return "", false
	}

	for _, item := range dto.ParseUserVisibleMenuItems(settings.CustomMenuItems) {
		if item.ID != menuItemID {
			continue
		}
		targetURL, err := validateCustomPageProxyURL(item.URL)
		if err != nil {
			return "", false
		}
		return targetURL, true
	}
	return "", false
}

func validateCustomPageProxyURL(rawURL string) (string, error) {
	if !isStandaloneCheckoutProxyURL(rawURL) {
		return "", errors.New("custom page is not a standalone checkout url")
	}
	validated, err := urlvalidator.ValidateHTTPSURL(rawURL, urlvalidator.ValidationOptions{
		AllowPrivate: false,
	})
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(validated)
	if err != nil || parsed.Hostname() == "" {
		return "", errors.New("invalid custom page url")
	}
	if err := urlvalidator.ValidateResolvedIP(parsed.Hostname()); err != nil {
		return "", err
	}
	return validated, nil
}

func isStandaloneCheckoutProxyURL(rawURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Hostname() == "" {
		return false
	}
	hostname := strings.ToLower(parsed.Hostname())
	pathname := strings.ToLower(parsed.EscapedPath())
	return strings.HasPrefix(hostname, "pay.") ||
		strings.Contains(hostname, ".pay.") ||
		strings.HasPrefix(pathname, "/shop/") ||
		strings.Contains(pathname, "/checkout")
}

func customPageProxyDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if ip := net.ParseIP(host); ip == nil || isUnsafeCustomPageProxyIP(ip) {
		_ = conn.Close()
		return nil, errors.New("resolved ip is not allowed")
	}
	return conn, nil
}

func isUnsafeCustomPageProxyIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified()
}

func readProxyBody(body io.Reader) ([]byte, error) {
	limited := io.LimitReader(body, customPageProxyMaxBodyBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(data) > customPageProxyMaxBodyBytes {
		return nil, errors.New("response body too large")
	}
	return data, nil
}

func isHTMLContentType(contentType string) bool {
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	return mediaType == "" || mediaType == "text/html" || mediaType == "application/xhtml+xml"
}

func injectBaseHref(body []byte, pageURL string) []byte {
	baseTag := []byte(`<base href="` + html.EscapeString(pageURL) + `">`)
	lower := bytes.ToLower(body)
	headStart := bytes.Index(lower, []byte("<head"))
	if headStart < 0 {
		return append(append(baseTag, '\n'), body...)
	}
	headEnd := bytes.IndexByte(body[headStart:], '>')
	if headEnd < 0 {
		return append(append(baseTag, '\n'), body...)
	}
	insertAt := headStart + headEnd + 1
	out := make([]byte, 0, len(body)+len(baseTag))
	out = append(out, body[:insertAt]...)
	out = append(out, baseTag...)
	out = append(out, body[insertAt:]...)
	return out
}

func setCustomPageProxyFrameHeaders(c *gin.Context) {
	c.Header("X-Frame-Options", "SAMEORIGIN")
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("Content-Security-Policy", strings.Join([]string{
		"default-src * data: blob:",
		"script-src * 'unsafe-inline' 'unsafe-eval'",
		"style-src * 'unsafe-inline'",
		"img-src * data: blob:",
		"font-src * data:",
		"connect-src *",
		"frame-src * data: blob:",
		"media-src * blob:",
		"base-uri *",
		"form-action *",
		"frame-ancestors 'self'",
	}, "; "))
}
