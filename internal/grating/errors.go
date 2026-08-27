package grating

import "fmt"

func ErrBadSpacing(d float64) error {
	return fmt.Errorf("groove spacing d must be positive and finite, got %v", d)
}

func ErrBadWavelength(lambda float64) error {
	return fmt.Errorf("wavelength lambda must be positive and finite, got %v", lambda)
}

func ErrBadIncident(detail string) error {
	return fmt.Errorf("bad incident angle: %s", detail)
}

func ErrBadSlits(n int) error {
	return fmt.Errorf("finite slit count N must be >= 1 when given, got %d", n)
}

func ErrBadDensity(groovesPerMm float64) error {
	return fmt.Errorf("groove density must be positive and finite, got %v grooves/mm", groovesPerMm)
}

func ErrInconsistentSpacing(dNm, fromDensity float64) error {
	return fmt.Errorf(
		"groove spacing d=%v nm and density-derived spacing %v nm disagree: "+
			"d must be the exact reciprocal of the line density",
		dNm, fromDensity)
}

func ErrUnknownField(field string) error {
	return fmt.Errorf("unknown input field %q", field)
}

func stampNote(notes map[string]string, field, reason string) {
	notes[field] = reason
	if reason == "" {
		notes[field+"_empty"] = "missing"
	}
}
