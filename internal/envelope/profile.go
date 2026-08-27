package envelope

import (
	"fmt"
	"math"

	"grating-eq/internal/grating"
)

type Sample struct {
	AngleRad   float64
	Phase      float64
	Normalized float64
	Exists     bool
}

type Profile struct {
	Order   int
	Slits   int
	PeakRad float64
	Samples []Sample
}

func SampleAround(g grating.Grating, m, n, count int) (Profile, error) {
	out := Profile{Order: m, Slits: n}
	if count < 3 {
		return out, fmt.Errorf("envelope: profile needs at least 3 samples")
	}
	if n < 1 {
		return out, grating.ErrBadSlits(n)
	}
	res, err := g.Eq(m)
	if err != nil {
		return out, err
	}
	if !res.Exists {
		return out, fmt.Errorf("envelope: order %d is not visible", m)
	}
	out.PeakRad = res.AngleRad
	hw, err := PrincipalHalfWidth(g, m, n)
	if err != nil {
		return out, err
	}
	span := 3 * hw.Radians
	out.Samples = make([]Sample, 0, count)
	for i := 0; i < count; i++ {
		frac := float64(i)/float64(count-1)*2 - 1
		theta := res.AngleRad + frac*span
		s := Sample{AngleRad: theta}
		if math.Abs(math.Sin(theta)) > 1 {
			out.Samples = append(out.Samples, s)
			continue
		}
		delta, err := PhaseDelta(g, theta)
		if err != nil {
			s.Exists = false
			out.Samples = append(out.Samples, s)
			continue
		}
		norm, err := NormalizedArray(n, delta)
		if err != nil {
			return out, err
		}
		s.Phase = delta
		s.Normalized = norm
		s.Exists = true
		out.Samples = append(out.Samples, s)
	}
	return out, nil
}

func (p Profile) PeakNormalized() float64 {
	best := 0.0
	for _, s := range p.Samples {
		if s.Exists && s.Normalized > best {
			best = s.Normalized
		}
	}
	return best
}

func (p Profile) MinAwayFromPeak(minOffsetRad float64) (float64, bool) {
	found := false
	minv := 0.0
	for _, s := range p.Samples {
		if !s.Exists {
			continue
		}
		if math.Abs(s.AngleRad-p.PeakRad) < minOffsetRad {
			continue
		}
		if !found || s.Normalized < minv {
			minv = s.Normalized
			found = true
		}
	}
	return minv, found
}

func NullIntensityAtFirstZero(g grating.Grating, m, n int) (float64, error) {
	delta0, err := FirstZeroPhase(n)
	if err != nil {
		return 0, err
	}
	centre := PhaseAtOrder(m)
	return NormalizedArray(n, centre+delta0)
}
