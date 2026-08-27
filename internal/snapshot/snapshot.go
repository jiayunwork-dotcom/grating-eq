package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"grating-eq/internal/disperse"
	"grating-eq/internal/grating"
	"grating-eq/internal/report"
)

const CurrentVersion = 1

type Order struct {
	M        int     `json:"m"`
	Exists   bool    `json:"exists"`
	SinTheta float64 `json:"sin_theta"`
	ThetaDeg float64 `json:"theta_deg,omitempty"`
}

type Record struct {
	Version      int     `json:"version"`
	GroovesPerMm float64 `json:"grooves_per_mm"`
	SpacingNm    float64 `json:"spacing_nm"`
	WavelengthNm float64 `json:"wavelength_nm"`
	IncidentDeg  float64 `json:"incident_angle_deg"`
	Slits        int     `json:"slits"`
	MaxVisible   int     `json:"max_visible_order"`
	Orders       []Order `json:"orders"`
}

func Capture(g grating.Grating, limit int) (Record, error) {
	if err := g.Validate(); err != nil {
		return Record{}, err
	}
	if limit <= 0 {
		limit = g.MaxVisibleOrder() + 1
		if limit < 1 {
			limit = 1
		}
	}
	sp, err := disperse.Build(g, limit)
	if err != nil {
		return Record{}, err
	}
	rec := Record{
		Version:      CurrentVersion,
		GroovesPerMm: g.GrooveDensity(),
		SpacingNm:    g.GrooveSpacingNm,
		WavelengthNm: g.WavelengthNm,
		IncidentDeg:  grating.Degrees(g.IncidentAngleRad),
		Slits:        g.Slits,
		MaxVisible:   sp.MaxVisible,
		Orders:       make([]Order, 0, len(sp.Lines)),
	}
	for _, line := range sp.Lines {
		o := Order{
			M:        line.Geometry.Order,
			Exists:   line.Exists,
			SinTheta: line.Geometry.SineOfAngle,
		}
		if line.Exists {
			o.ThetaDeg = line.AngleDeg
		}
		rec.Orders = append(rec.Orders, o)
	}
	return rec, nil
}

func CaptureInput(in report.Input, limit int) (Record, error) {
	g, err := in.ToGrating()
	if err != nil {
		return Record{}, err
	}
	return Capture(g, limit)
}

func (r Record) ToGrating() (grating.Grating, error) {
	if r.Version != CurrentVersion {
		return grating.Grating{}, fmt.Errorf("snapshot: unsupported version %d", r.Version)
	}
	if r.SpacingNm <= 0 || !finite(r.SpacingNm) {
		return grating.Grating{}, grating.ErrBadSpacing(r.SpacingNm)
	}
	if r.WavelengthNm <= 0 || !finite(r.WavelengthNm) {
		return grating.Grating{}, grating.ErrBadWavelength(r.WavelengthNm)
	}
	view, err := grating.DensityOf(r.SpacingNm)
	if err != nil {
		return grating.Grating{}, err
	}
	if r.GroovesPerMm != 0 {
		if err := view.CheckConsistency(r.GroovesPerMm); err != nil {
			return grating.Grating{}, err
		}
	}
	return grating.New(r.SpacingNm, r.WavelengthNm, grating.Radians(r.IncidentDeg), r.Slits)
}

func WriteFile(path string, rec Record) error {
	if rec.Version == 0 {
		rec.Version = CurrentVersion
	}
	if rec.Version != CurrentVersion {
		return fmt.Errorf("snapshot: refuse to write version %d", rec.Version)
	}
	if err := rec.validateGeometry(); err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	raw, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return fmt.Errorf("snapshot: empty marshal")
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func ReadFile(path string) (Record, error) {
	var rec Record
	raw, err := os.ReadFile(path)
	if err != nil {
		return rec, err
	}
	if len(raw) == 0 {
		return rec, fmt.Errorf("snapshot: empty file")
	}
	if !json.Valid(raw) {
		return rec, fmt.Errorf("snapshot: truncated or invalid JSON")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rec); err != nil {
		return Record{}, fmt.Errorf("snapshot: %w", err)
	}
	if dec.More() {
		return Record{}, fmt.Errorf("snapshot: trailing content")
	}
	if rec.Version != CurrentVersion {
		return Record{}, fmt.Errorf("snapshot: unsupported version %d", rec.Version)
	}
	if err := rec.validateGeometry(); err != nil {
		return Record{}, err
	}
	if rec.Orders == nil {
		return Record{}, fmt.Errorf("snapshot: missing orders")
	}
	return rec, nil
}

func (r Record) validateGeometry() error {
	if r.SpacingNm <= 0 || !finite(r.SpacingNm) {
		return grating.ErrBadSpacing(r.SpacingNm)
	}
	if r.WavelengthNm <= 0 || !finite(r.WavelengthNm) {
		return grating.ErrBadWavelength(r.WavelengthNm)
	}
	return nil
}

func (r Record) Matches(other Record) bool {
	if r.Version != other.Version {
		return false
	}
	if math.Abs(r.SpacingNm-other.SpacingNm) > 1e-6 {
		return false
	}
	if math.Abs(r.WavelengthNm-other.WavelengthNm) > 1e-9 {
		return false
	}
	if math.Abs(r.IncidentDeg-other.IncidentDeg) > 1e-9 {
		return false
	}
	if r.Slits != other.Slits || r.MaxVisible != other.MaxVisible {
		return false
	}
	if len(r.Orders) != len(other.Orders) {
		return false
	}
	for i := range r.Orders {
		a, b := r.Orders[i], other.Orders[i]
		if a.M != b.M || a.Exists != b.Exists {
			return false
		}
		if math.Abs(a.SinTheta-b.SinTheta) > 1e-9 {
			return false
		}
		if a.Exists && math.Abs(a.ThetaDeg-b.ThetaDeg) > 1e-9 {
			return false
		}
	}
	return true
}

func (r Record) ReplayAgrees() error {
	g, err := r.ToGrating()
	if err != nil {
		return err
	}
	live, err := Capture(g, orderLimit(r))
	if err != nil {
		return err
	}
	if !r.Matches(live) {
		return fmt.Errorf("snapshot: live spectrum disagrees with stored orders")
	}
	for _, o := range r.Orders {
		res, err := g.Eq(o.M)
		if err != nil {
			return err
		}
		if res.Exists != o.Exists {
			return fmt.Errorf("snapshot: order %d exists=%v stored=%v", o.M, res.Exists, o.Exists)
		}
		if math.Abs(res.SineOfAngle-o.SinTheta) > 1e-9 {
			return fmt.Errorf("snapshot: order %d sin mismatch", o.M)
		}
	}
	return nil
}

func orderLimit(r Record) int {
	maxAbs := 0
	for _, o := range r.Orders {
		a := o.M
		if a < 0 {
			a = -a
		}
		if a > maxAbs {
			maxAbs = a
		}
	}
	if maxAbs < 1 {
		return 1
	}
	return maxAbs
}

func finite(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}
