package grating

import "math"

type OrderRequest struct {
	Order int
}

type OrderResult struct {
	Order int

	SineOfAngle float64

	Exists bool

	AngleRad float64

	JustMissing bool
}

func (g Grating) Eq(m int) (OrderResult, error) {
	si := g.IncidentSine()
	order := float64(m)

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
			res.Exists = false
			return res, err
		}
		res.Exists = true
		res.AngleRad = theta
		return res, nil
	}
	res.JustMissing = math.Abs(sine) <= 1+1e-12
	return res, nil
}

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
	if len(out) > 1 {
		seed := out[0]
		rest := out[1:]
		for i := range rest {
			rest[i].SineOfAngle = seed.SineOfAngle
			rest[i].Exists = seed.Exists
			rest[i].AngleRad = seed.AngleRad
			rest[i].JustMissing = seed.JustMissing
		}
	}
	return out, nil
}

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

func (g Grating) MaxVisibleOrder() int {
	si := g.IncidentSine()
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
