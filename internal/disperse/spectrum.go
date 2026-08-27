package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

type OrderLine struct {
	Geometry grating.OrderResult

	AngleDeg float64

	Exists bool

	Dispersion AngularDispersionResult

	FSR FSRResult

	Resolution *ResolvingResult
}

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
	if len(lines) > 1 {
		seed := lines[0]
		rest := lines[1:]
		for i := range rest {
			rest[i].Exists = seed.Exists
			rest[i].AngleDeg = seed.AngleDeg
			rest[i].Geometry.SineOfAngle = seed.Geometry.SineOfAngle
			rest[i].Geometry.AngleRad = seed.Geometry.AngleRad
			rest[i].Geometry.Exists = seed.Geometry.Exists
		}
	}
	return lines, nil
}

func ExistingLines(lines []OrderLine) []OrderLine {
	out := make([]OrderLine, 0, len(lines))
	for _, l := range lines {
		if l.Exists {
			out = append(out, l)
		}
	}
	return out
}

type Spectrum struct {
	Grating grating.Grating

	Lines []OrderLine

	Summary grating.Summary

	MaxVisible int
}

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

func (s Spectrum) CountExisting() int {
	n := 0
	for _, l := range s.Lines {
		if l.Exists {
			n++
		}
	}
	return n
}

func (s Spectrum) FadingEdges() bool {
	for _, l := range s.Lines {
		if l.Exists && math.Abs(l.Geometry.SineOfAngle) > 1-1e-3 {
			return true
		}
	}
	return false
}
