package grating

import "math"

// SymmetryPair groups the results of +m and -m for a grating at normal
// incidence, where the two orders must be exactly symmetric about the
// normal.
type SymmetryPair struct {
	// Plus is the result for order +m.
	Plus OrderResult

	// Minus is the result for order -m.
	Minus OrderResult
}

// SymmetricPair returns the +m and -m results at normal incidence. When the
// incident angle is exactly zero, the transmission equation reduces to
// sin(theta_m) = m*lambda/d, which is odd in m, so the pair must satisfy
// theta_{-m} = -theta_m.
func (g Grating) SymmetricPair(m int) (SymmetryPair, error) {
	if math.Abs(g.IncidentAngleRad) > 1e-12 {
		return SymmetryPair{}, ErrNotNormalIncidence()
	}
	plus, err := g.Eq(m)
	if err != nil {
		return SymmetryPair{}, err
	}
	minus, err := g.Eq(-m)
	if err != nil {
		return SymmetryPair{}, err
	}
	return SymmetryPair{Plus: plus, Minus: minus}, nil
}

// ErrNotNormalIncidence reports a symmetry query at non-normal incidence.
func ErrNotNormalIncidence() error {
	return &notNormalError{}
}

type notNormalError struct{}

func (e *notNormalError) Error() string {
	return "order symmetry is only defined at normal incidence"
}

// Antisymmetric reports whether the pair satisfies theta_{+m} = -theta_{-m}
// to within tolerance. This is the invariant that a sign error in the
// grating equation (adding sin(theta_i) instead of subtracting it) would
// break, because then theta_{+m} and theta_{-m} would no longer mirror each
// other at normal incidence.
func (p SymmetryPair) Antisymmetric() bool {
	if !p.Plus.Exists || !p.Minus.Exists {
		return p.Plus.Exists == p.Minus.Exists
	}
	return math.Abs(p.Plus.AngleRad+p.Minus.AngleRad) < AngleTolerance
}

// SameWavelengthSide reports whether two orders diffract to the same side
// of the normal.
func SameWavelengthSide(a, b OrderResult) bool {
	sa := math.Signbit(a.SineOfAngle)
	sb := math.Signbit(b.SineOfAngle)
	if math.Abs(a.SineOfAngle) < 1e-15 || math.Abs(b.SineOfAngle) < 1e-15 {
		// An order sitting on the normal counts as matching either side,
		// which keeps "all orders on one side" tests meaningful.
		return true
	}
	return sa == sb
}

// HigherOrderReport lists how many orders exist up to a given |m| bound.
// It is used by the cross-rule tests and by the "checks" style output.
type HigherOrderReport struct {
	// Bound is the scan limit that produced the counts.
	Bound int

	// Existing is the number of orders that exist within the bound.
	Existing int

	// Missing is the number that fail the |sin| <= 1 test.
	Missing int
}

// CountOrders runs the existence test across the scan window and tallies
// the outcome. The counts are cheap to compute and make dense-grating
// behaviour directly observable.
func (g Grating) CountOrders(limit int) (HigherOrderReport, error) {
	all, err := g.Orders(limit)
	if err != nil {
		return HigherOrderReport{}, err
	}
	rep := HigherOrderReport{Bound: limit}
	for _, r := range all {
		if r.Exists {
			rep.Existing++
		} else {
			rep.Missing++
		}
	}
	return rep, nil
}
