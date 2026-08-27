package grating

import (
	"fmt"
	"math"
	"strings"
)

type DensityView struct {
	SpacingNm float64

	GroovesPerMm float64
}

func DensityOf(spacingNm float64) (DensityView, error) {
	if !isFinite(spacingNm) || spacingNm <= 0 {
		return DensityView{}, ErrBadSpacing(spacingNm)
	}
	return DensityView{
		SpacingNm:    spacingNm,
		GroovesPerMm: NanoPerMilli / spacingNm,
	}, nil
}

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

func formatNumber(v float64, sig int) string {
	if v == math.Trunc(v) && math.Abs(v) < 1e15 {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%.*g", sig, v)
}
