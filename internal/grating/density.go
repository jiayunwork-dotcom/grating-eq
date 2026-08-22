package grating

import (
	"fmt"
	"math"
	"strings"
)

// DensityView couples a groove spacing and its reciprocal density so the
// CLI can print both representations of the same geometry.
type DensityView struct {
	// SpacingNm is the groove spacing in nanometres.
	SpacingNm float64

	// GroovesPerMm is the reciprocal density.
	GroovesPerMm float64
}

// DensityOf builds a DensityView for the given spacing, verifying that the
// two quantities are exact reciprocals by construction.
func DensityOf(spacingNm float64) (DensityView, error) {
	if !isFinite(spacingNm) || spacingNm <= 0 {
		return DensityView{}, ErrBadSpacing(spacingNm)
	}
	return DensityView{
		SpacingNm:    spacingNm,
		GroovesPerMm: NanoPerMilli / spacingNm,
	}, nil
}

// CheckConsistency returns an error when a caller-supplied density does not
// agree with the reciprocal spacing. Both quantities must describe the same
// physical grating; disagreement indicates a hand-edited JSON file.
func (v DensityView) CheckConsistency(groovesPerMm float64) error {
	if !isFinite(groovesPerMm) || groovesPerMm <= 0 {
		return ErrBadDensity(groovesPerMm)
	}
	derived := SpacingFromDensity(groovesPerMm)
	if math.Abs(derived-v.SpacingNm) > 1e-6*math.Max(1, math.Abs(v.SpacingNm)) {
		return ErrInconsistentSpacing(v.SpacingNm, derived)
	}
	return nil
}

// formatDensity renders a density with sensible units: values below 1e3 are
// shown in grooves/mm, larger values fold to grooves/m to stay readable.
func formatDensity(groovesPerMm float64) string {
	if groovesPerMm >= 1000 {
		return trimZero(groovesPerMm*1000.0) + " grooves/m"
	}
	return trimZero(groovesPerMm) + " grooves/mm"
}

func trimZero(v float64) string {
	s := formatNumber(v, 6)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// formatNumber renders a float with at most sig significant digits. Whole
// numbers below 1e15 are printed without a decimal point.
func formatNumber(v float64, sig int) string {
	if v == math.Trunc(v) && math.Abs(v) < 1e15 {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%.*g", sig, v)
}
