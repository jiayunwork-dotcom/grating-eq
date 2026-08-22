package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

// OrderLine is the fully assembled row of the orders table: the geometry
// from the grating package combined with the dispersion quantities.
type OrderLine struct {
	// Geometry is the order result from the grating equation.
	Geometry grating.OrderResult

	// AngleDeg is theta_m in degrees (empty for missing orders).
	AngleDeg float64

	// Exists mirrors Geometry.Exists.
	Exists bool

	// Dispersion is the angular dispersion in deg/nm.
	Dispersion AngularDispersionResult

	// FSR is the free spectral range in nanometres.
	FSR FSRResult

	// Resolution is present when a finite slit count was supplied.
	Resolution *ResolvingResult
}

// BuildLines evaluates every order in the scan window and decorates it
// with the dispersing quantities.
func BuildLines(g grating.Grating, limit int) ([]OrderLine, error) {
	orders, err := g.Orders(limit)
	if err != nil {
		return nil, err
	}
	lines := make([]OrderLine, 0, len(orders))
	for _, res := range orders {
		line := OrderLine{
			Geometry:   res,
			Exists:     res.Exists,
			Dispersion: AngularDispersion(g, res),
			FSR:        FreeSpectralRange(g.WavelengthNm, res.Order),
		}
		if res.Exists {
			line.AngleDeg = grating.Degrees(res.AngleRad)
		}
		if g.Slits > 0 && res.Order != 0 {
			r, err := ResolvingPower(g.WavelengthNm, res.Order, g.Slits)
			if err != nil {
				return nil, err
			}
			cp := r
			line.Resolution = &cp
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// ExistingLines filters the table to the orders that actually diffract.
func ExistingLines(lines []OrderLine) []OrderLine {
	out := make([]OrderLine, 0, len(lines))
	for _, l := range lines {
		if l.Exists {
			out = append(out, l)
		}
	}
	return out
}

// Spectrum holds the assembled output of one orders run.
type Spectrum struct {
	// Grating is the input geometry.
	Grating grating.Grating

	// Lines are the per-order rows.
	Lines []OrderLine

	// Summary describes the grating in prose form.
	Summary grating.Summary

	// MaxVisible is the highest order that exists.
	MaxVisible int
}

// Build assembles a Spectrum for the orders subcommand.
func Build(g grating.Grating, limit int) (Spectrum, error) {
	if limit <= 0 {
		limit = grating.DefaultMaxOrderScan
	}
	lines, err := BuildLines(g, limit)
	if err != nil {
		return Spectrum{}, err
	}
	return Spectrum{
		Grating:    g,
		Lines:      lines,
		Summary:    g.Summarize(),
		MaxVisible: g.MaxVisibleOrder(),
	}, nil
}

// CountExisting returns the number of visible orders in the spectrum.
func (s Spectrum) CountExisting() int {
	n := 0
	for _, l := range s.Lines {
		if l.Exists {
			n++
		}
	}
	return n
}

// FadingEdges returns true when the highest existing order sits within
// 1e-3 rad of the grazing angle, the signature of a near-dense grating.
func (s Spectrum) FadingEdges() bool {
	for _, l := range s.Lines {
		if l.Exists && math.Abs(l.Geometry.SineOfAngle) > 1-1e-3 {
			return true
		}
	}
	return false
}
