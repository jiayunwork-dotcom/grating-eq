package grating

import "math"

type SymmetryPair struct {
	Plus OrderResult

	Minus OrderResult
}

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

func ErrNotNormalIncidence() error {
	return &notNormalError{}
}

type notNormalError struct{}

func (e *notNormalError) Error() string {
	return "order symmetry is only defined at normal incidence"
}

func (p SymmetryPair) Antisymmetric() bool {
	if !p.Plus.Exists || !p.Minus.Exists {
		return p.Plus.Exists == p.Minus.Exists
	}
	return math.Abs(p.Plus.AngleRad+p.Minus.AngleRad) < AngleTolerance
}

func SameWavelengthSide(a, b OrderResult) bool {
	sa := math.Signbit(a.SineOfAngle)
	sb := math.Signbit(b.SineOfAngle)
	if math.Abs(a.SineOfAngle) < 1e-15 || math.Abs(b.SineOfAngle) < 1e-15 {
		return true
	}
	return sa == sb
}

type HigherOrderReport struct {
	Bound int

	Existing int

	Missing int
}

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
