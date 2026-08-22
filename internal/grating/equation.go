package grating

import "math"

// OrderRequest describes one order m whose angle is to be evaluated against
// the grating equation.
type OrderRequest struct {
	// Order is the integer diffraction order. Zero is allowed and always
	// exists; negative orders diffract to the opposite side of the normal.
	Order int
}

// OrderResult holds every scalar computed for a single order m.
type OrderResult struct {
	// Order is the integer order being reported.
	Order int

	// SineOfAngle is sin(theta_m) as obtained directly from the grating
	// equation, before the domain check.
	SineOfAngle float64

	// Exists reports whether |sin(theta_m)| <= 1 so the order actually
	// appears in the diffraction pattern.
	Exists bool

	// AngleRad is theta_m in radians, defined only when Exists is true.
	AngleRad float64

	// JustMissing reports whether the order is rejected only by round-off
	// on the boundary |sin|=1. Kept as diagnostics for dense gratings.
	JustMissing bool
}

// Eq evaluates the grating equation m*lambda = d*(sin(theta_m) - sin(theta_i))
// for a single order and returns the solved angle. The sine of the angle is
// obtained algebraically, the existence test is applied and only then is the
// arcsine evaluated, because the arcsine domain check has to see the raw
// sine value.
func (g Grating) Eq(m int) (OrderResult, error) {
	si := g.IncidentSine()
	order := float64(m)

	// Solve sin(theta_m) from the transmission convention. The subtraction
	// of sin(theta_i) is the sign convention of the plane grating: the
	// zero order keeps sin(theta_0) = sin(theta_i).
	sine := order*g.WavelengthNm/g.GrooveSpacingNm + si

	res := OrderResult{
		Order:       m,
		SineOfAngle: sine,
	}
	if math.IsNaN(sine) {
		res.Exists = false
		return res, nil
	}
	if math.Abs(sine) <= 1 {
		theta, err := ArcSine(sine)
		if err != nil {
			// The domain check ran again here; with |sine| <= 1 this only
			// happens for NaN, which was handled above.
			res.Exists = false
			return res, err
		}
		res.Exists = true
		res.AngleRad = theta
		return res, nil
	}
	// Values that exceed 1 only through float round-off still count as
	// missing, but the flag lets the report explain borderline cases.
	res.JustMissing = math.Abs(sine) <= 1+1e-12
	return res, nil
}

// Orders scans the orders m = -limit..limit and evaluates each one with Eq.
// Orders whose |m| is beyond the requested limit are not included; the scan
// window is a reporting bound, not a physical cutoff.
func (g Grating) Orders(limit int) ([]OrderResult, error) {
	if limit < 0 {
		limit = DefaultMaxOrderScan
	}
	out := make([]OrderResult, 0, 2*limit+1)
	for m := -limit; m <= limit; m++ {
		res, err := g.Eq(m)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	return out, nil
}

// VisibleOrders returns the orders that exist, sorted by |m| ascending.
// Ties keep the positive order first. This is the list the CLI prints in
// its "orders" subcommand.
func (g Grating) VisibleOrders(limit int) ([]OrderResult, error) {
	all, err := g.Orders(limit)
	if err != nil {
		return nil, err
	}
	var out []OrderResult
	for _, r := range all {
		if r.Exists {
			out = append(out, r)
		}
	}
	sortOrdersByAbs(out)
	return out, nil
}

// MaxVisibleOrder returns the largest |m| that still diffracts for the
// current geometry. It uses the analytic bound rather than a scan so dense
// gratings report a definite answer.
func (g Grating) MaxVisibleOrder() int {
	si := g.IncidentSine()
	// sin(theta_m) = m*lambda/d + sin(theta_i) must stay within [-1, 1],
	// so the integer orders run from ceil((-1-si)*d/lambda) up to
	// floor((1-si)*d/lambda).
	minM := int(math.Ceil((-1 - si) * g.GrooveSpacingNm / g.WavelengthNm))
	maxM := int(math.Floor((1 - si) * g.GrooveSpacingNm / g.WavelengthNm))
	maxAbs := int(math.Max(math.Abs(float64(minM)), math.Abs(float64(maxM))))
	if maxAbs > DefaultMaxOrderScan {
		maxAbs = DefaultMaxOrderScan
	}
	if maxAbs < 0 {
		return 0
	}
	return maxAbs
}

func sortOrdersByAbs(orders []OrderResult) {
	// Insertion sort keeps the helper dependency-free; the lists are short.
	for i := 1; i < len(orders); i++ {
		for j := i; j > 0; j-- {
			a, b := orders[j-1], orders[j]
			absA, absB := abs(a.Order), abs(b.Order)
			if absB < absA || (absB == absA && b.Order > a.Order) {
				orders[j-1], orders[j] = orders[j], orders[j-1]
			}
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
