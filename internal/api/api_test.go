package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func exampleBody(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile("../../example/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestHealth(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("health code %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "ok") {
		t.Fatalf("health body %q", rec.Body.String())
	}
}

func TestOrdersEndpoint(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(exampleBody(t, "600lpmm.json")))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("orders code %d body %s", rec.Code, rec.Body.String())
	}
	var out map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	max, ok := out["max_visible_order"].(float64)
	if !ok || max < 1 {
		t.Fatalf("max_visible_order = %v", out["max_visible_order"])
	}
}

func TestOrdersRejectsBadSpacing(t *testing.T) {
	srv := New(DefaultConfig())
	raw := []byte(`{"d_nm":0,"wavelength_nm":550,"incident_angle_deg":0}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(raw))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d body %s", rec.Code, rec.Body.String())
	}
}

func TestOrdersRejectsEmpty(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(nil))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}

func TestChecksEndpoint(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/checks", bytes.NewReader(exampleBody(t, "600lpmm.json")))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("checks code %d body %s", rec.Code, rec.Body.String())
	}
	var resp checksResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Rules) == 0 {
		t.Fatal("expected cross rules")
	}
	for _, r := range resp.Rules {
		if !r.Holds {
			t.Fatalf("rule %q failed: %s", r.Name, r.Detail)
		}
	}
}

func TestOrdersMethodNotAllowed(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
