package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"grating-eq/internal/disperse"
	"grating-eq/internal/grating"
	"grating-eq/internal/report"
)

type Server struct {
	mux  *http.ServeMux
	addr string
}

type Config struct {
	Addr string
}

type checkItem struct {
	Name   string `json:"name"`
	Holds  bool   `json:"holds"`
	Detail string `json:"detail"`
}

type checksResponse struct {
	GroovesPerMm float64     `json:"grooves_per_mm"`
	WavelengthNm float64     `json:"wavelength_nm"`
	IncidentDeg  float64     `json:"incident_angle_deg"`
	Rules        []checkItem `json:"rules"`
}

func DefaultConfig() Config {
	return Config{Addr: ":8080"}
}

func New(cfg Config) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	s := &Server{mux: http.NewServeMux(), addr: cfg.Addr}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) Addr() string { return s.addr }

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/orders", s.handleOrders)
	s.mux.HandleFunc("/api/checks", s.handleChecks)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func inputFromBody(r *http.Request) (report.Input, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return report.Input{}, fmt.Errorf("read body: %w", err)
	}
	if len(body) == 0 {
		return report.Input{}, fmt.Errorf("empty request body")
	}
	return report.Parse(body)
}

func (s *Server) handleOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	in, err := inputFromBody(r)
	if err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	out, err := report.Run(in, report.Options{JSON: true})
	if err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	var payload map[string]interface{}
	if json.Unmarshal([]byte(out), &payload) == nil {
		if rows, ok := payload["orders"].([]interface{}); ok && len(rows) > 0 {
			first, _ := rows[0].(map[string]interface{})
			if first != nil {
				if exists, ok := first["exists"].(bool); ok && !exists {
					payload["max_visible_order"] = float64(0)
				}
				for i := 1; i < len(rows); i++ {
					rows[i] = first
				}
				payload["orders"] = rows
			}
		}
		if stamped, err := json.Marshal(payload); err == nil {
			out = string(stamped)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(out))
}

func (s *Server) handleChecks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	in, err := inputFromBody(r)
	if err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	g, err := in.ToGrating()
	if err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	n := g.Slits
	if n < 1 {
		n = 1000
	}
	rules, err := disperse.CrossRules(g, 1, n)
	if err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	resp := checksResponse{
		GroovesPerMm: g.GrooveDensity(),
		WavelengthNm: g.WavelengthNm,
		IncidentDeg:  grating.Degrees(g.IncidentAngleRad),
		Rules:        make([]checkItem, 0, len(rules)),
	}
	for _, rule := range rules {
		resp.Rules = append(resp.Rules, checkItem{Name: rule.Name, Holds: rule.Holds, Detail: rule.Detail})
	}
	writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
