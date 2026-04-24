package routes

import (
	"strings"
	"testing"
)

func TestIsStandaloneCheckoutProxyURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{name: "pay subdomain shop", raw: "https://pay.ldxp.cn/shop/BSEJH4PV/4j9om0", want: true},
		{name: "checkout path", raw: "https://example.com/checkout/session", want: true},
		{name: "middle pay subdomain", raw: "https://shop.pay.example.com/products", want: true},
		{name: "regular docs page", raw: "https://example.com/docs", want: false},
		{name: "invalid url", raw: "not a url", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStandaloneCheckoutProxyURL(tt.raw); got != tt.want {
				t.Fatalf("isStandaloneCheckoutProxyURL(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestValidateCustomPageProxyURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{name: "rejects non checkout url", raw: "https://example.com/docs", wantErr: true},
		{name: "rejects http checkout url", raw: "http://example.com/checkout/session", wantErr: true},
		{name: "rejects private checkout url", raw: "https://127.0.0.1/checkout/session", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateCustomPageProxyURL(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateCustomPageProxyURL(%q) err = %v, wantErr %v", tt.raw, err, tt.wantErr)
			}
		})
	}
}

func TestInjectBaseHref(t *testing.T) {
	body := []byte("<html><head><title>Pay</title></head><body></body></html>")
	got := string(injectBaseHref(body, "https://pay.example.com/shop/a?x=1&y=2"))

	if !strings.Contains(got, `<base href="https://pay.example.com/shop/a?x=1&amp;y=2">`) {
		t.Fatalf("expected escaped base tag, got %s", got)
	}
	if strings.Index(got, "<base ") > strings.Index(got, "<title>") {
		t.Fatalf("expected base tag before title, got %s", got)
	}
}
