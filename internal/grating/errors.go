package grating

import "fmt"

// Error values produced by validation. They are returned from Validate and
// propagated by ParseFile, so the CLI can print a plain message on stderr
// and exit non-zero.

// ErrBadSpacing describes a groove spacing that is missing, non-finite or
// not positive. A zero or negative spacing is a degenerate grating that
// produces no diffraction at all.
func ErrBadSpacing(d float64) error {
	return fmt.Errorf("groove spacing d must be positive and finite, got %v", d)
}

// ErrBadWavelength describes a wavelength that is not positive. Negative
// wavelengths are meaningless and a zero wavelength makes every order angle
// collapse onto the incident direction.
func ErrBadWavelength(lambda float64) error {
	return fmt.Errorf("wavelength lambda must be positive and finite, got %v", lambda)
}

// ErrBadIncident describes an incident angle whose sine falls outside
// [-1, 1] or whose value is not finite. Such an angle cannot point at the
// grating from any real direction.
func ErrBadIncident(detail string) error {
	return fmt.Errorf("bad incident angle: %s", detail)
}

// ErrBadSlits describes a finite slit count below 1. A grating with fewer
// than one illuminated groove resolves nothing.
func ErrBadSlits(n int) error {
	return fmt.Errorf("finite slit count N must be >= 1 when given, got %d", n)
}

// ErrBadDensity describes a groove density that is not positive and finite.
func ErrBadDensity(groovesPerMm float64) error {
	return fmt.Errorf("groove density must be positive and finite, got %v grooves/mm", groovesPerMm)
}

// ErrInconsistentSpacing is returned when a JSON file supplies both a groove
// density and a groove spacing that do not agree with each other. The two
// descriptions must always be exact reciprocals.
func ErrInconsistentSpacing(dNm, fromDensity float64) error {
	return fmt.Errorf(
		"groove spacing d=%v nm and density-derived spacing %v nm disagree: "+
			"d must be the exact reciprocal of the line density",
		dNm, fromDensity)
}

// ErrUnknownField is returned by the JSON decoder when the input carries a
// field that this tool does not understand, so typos surface immediately.
func ErrUnknownField(field string) error {
	return fmt.Errorf("unknown input field %q", field)
}
