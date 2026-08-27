package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"

	"grating-eq/internal/grating"
)

type Input struct {
	GroovesPerMm *float64 `json:"grooves_per_mm"`

	SpacingNm *float64 `json:"d_nm"`

	WavelengthNm float64 `json:"wavelength_nm"`

	IncidentDeg float64 `json:"incident_angle_deg"`

	Slits int `json:"slits"`
}

func (in Input) HasGeometry() bool {
	return in.GroovesPerMm != nil || in.SpacingNm != nil
}

func ParseFile(path string) (Input, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Input{}, fmt.Errorf("read grating file %q: %w", path, err)
	}
	return Parse(data)
}

func Parse(data []byte) (Input, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var in Input
	if err := dec.Decode(&in); err != nil {
		return Input{}, fmt.Errorf("parse grating JSON: %w", err)
	}
	if dec.More() {
		return Input{}, fmt.Errorf("parse grating JSON: unexpected trailing content")
	}
	return in, nil
}

func ParseReader(r io.Reader) (Input, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Input{}, fmt.Errorf("read grating JSON: %w", err)
	}
	return Parse(data)
}

func (in Input) ToGrating() (grating.Grating, error) {
	if err := in.validate(); err != nil {
		return grating.Grating{}, err
	}
	var d float64
	switch {
	case in.SpacingNm != nil && in.GroovesPerMm != nil:
		if *in.SpacingNm <= 0 || !isFinite(*in.SpacingNm) {
			return grating.Grating{}, grating.ErrBadSpacing(*in.SpacingNm)
		}
		if *in.GroovesPerMm <= 0 || !isFinite(*in.GroovesPerMm) {
			return grating.Grating{}, grating.ErrBadDensity(*in.GroovesPerMm)
		}
		view, err := grating.DensityOf(*in.SpacingNm)
		if err != nil {
			return grating.Grating{}, err
		}
		if err := view.CheckConsistency(*in.GroovesPerMm); err != nil {
			return grating.Grating{}, err
		}
		d = *in.SpacingNm
	case in.SpacingNm != nil:
		if *in.SpacingNm <= 0 || !isFinite(*in.SpacingNm) {
			return grating.Grating{}, grating.ErrBadSpacing(*in.SpacingNm)
		}
		d = *in.SpacingNm
	case in.GroovesPerMm != nil:
		if *in.GroovesPerMm <= 0 || !isFinite(*in.GroovesPerMm) {
			return grating.Grating{}, grating.ErrBadDensity(*in.GroovesPerMm)
		}
		d = grating.SpacingFromDensity(*in.GroovesPerMm)
	default:
		return grating.Grating{}, fmt.Errorf("input must set groove density (grooves_per_mm) or spacing (d_nm)")
	}
	inc := grating.Radians(in.IncidentDeg)
	g, err := grating.New(d, in.WavelengthNm, inc, in.Slits)
	if err != nil {
		return grating.Grating{}, err
	}
	return g, nil
}

func (in Input) validate() error {
	if math.IsNaN(in.WavelengthNm) || math.IsInf(in.WavelengthNm, 0) || in.WavelengthNm <= 0 {
		return grating.ErrBadWavelength(in.WavelengthNm)
	}
	if in.Slits < 0 {
		return grating.ErrBadSlits(in.Slits)
	}
	return nil
}

func isFinite(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}
