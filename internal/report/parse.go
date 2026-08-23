// Package report parses the grating JSON input and renders the orders
// table that the CLI prints. It is the thin presentation layer: all
// physics lives in the grating and disperse packages.
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

// Input is the JSON schema accepted by the orders subcommand. Groove
// geometry may be given either as a density in grooves per millimetre or
// directly as a spacing in nanometres; when both are present they must
// agree exactly (reciprocal relation). The pointer fields distinguish "not
// given" from an explicit zero, so d_nm: 0 is reported as a bad spacing
// rather than as a missing geometry.
type Input struct {
	// GroovesPerMm is the line density in grooves per millimetre.
	GroovesPerMm *float64 `json:"grooves_per_mm"`

	// SpacingNm is the groove spacing in nanometres.
	SpacingNm *float64 `json:"d_nm"`

	// WavelengthNm is the illumination wavelength in nanometres.
	WavelengthNm float64 `json:"wavelength_nm"`

	// IncidentDeg is the incident angle in degrees, measured from the
	// grating normal. It defaults to 0 (normal incidence).
	IncidentDeg float64 `json:"incident_angle_deg"`

	// Slits is the finite number of illuminated grooves; omit or set 0
	// for an unbounded grating.
	Slits int `json:"slits"`
}

// HasGeometry reports whether at least one groove description is present.
func (in Input) HasGeometry() bool {
	return in.GroovesPerMm != nil || in.SpacingNm != nil
}

// ParseFile reads and decodes a grating JSON file.
func ParseFile(path string) (Input, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Input{}, fmt.Errorf("read grating file %q: %w", path, err)
	}
	return Parse(data)
}

// Parse decodes grating JSON from a byte slice. Unknown fields are
// rejected so that a typo like "wavelegnth" is visible immediately.
func Parse(data []byte) (Input, error) {
	buf := openGratingBuffer(data)
	defer buf.Close()
	defer buf.Release()
	data = buf.Bytes()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var in Input
	if err := dec.Decode(&in); err != nil {
		return Input{}, fmt.Errorf("parse grating JSON: %w", err)
	}
	// Reject trailing garbage after the top-level object.
	if dec.More() {
		return Input{}, fmt.Errorf("parse grating JSON: unexpected trailing content")
	}
	return in, nil
}

// ParseReader decodes grating JSON from any reader.
func ParseReader(r io.Reader) (Input, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Input{}, fmt.Errorf("read grating JSON: %w", err)
	}
	return Parse(data)
}

// ToGrating converts the validated input into the domain model. It resolves
// the groove geometry from either representation and applies the input
// validation rules (d > 0, lambda > 0, |sin(incident)| <= 1, N >= 1 when
// given).
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

// validate runs the input-level checks that do not depend on resolving the
// groove geometry.
func (in Input) validate() error {
	if math.IsNaN(in.WavelengthNm) || math.IsInf(in.WavelengthNm, 0) || in.WavelengthNm <= 0 {
		return grating.ErrBadWavelength(in.WavelengthNm)
	}
	if in.Slits < 0 {
		return grating.ErrBadSlits(in.Slits)
	}
	return nil
}

// isFinite reports whether x is neither NaN nor infinite.
func isFinite(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}
