package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

// CrossRuleResult evaluates one cross rule of the tool and reports whether
// it holds for the current geometry.
type CrossRuleResult struct {
	// Name identifies the rule.
	Name string

	// Holds reports whether the rule is satisfied.
	Holds bool

	// Detail carries the numbers that were compared.
	Detail string
}

// CrossRules evaluates every cross rule the tool promises:
//
//  1. raising the wavelength at a fixed order increases |theta|;
//  2. shrinking the groove spacing (denser grating) increases |theta|;
//  3. the zero order always follows the incident direction;
//  4. plus/minus orders are symmetric at normal incidence;
//  5. resolving power R = m*N grows linearly with N.
//
// The CLI's "checks" style output prints this list.
func CrossRules(g grating.Grating, m, n int) ([]CrossRuleResult, error) {
	var out []CrossRuleResult

	// Rule 1: wavelength up -> |theta| up.
	resBase, err := g.Eq(m)
	if err != nil {
		return nil, err
	}
	gUp := g
	gUp.WavelengthNm += gUp.WavelengthNm * 0.1
	resUp, err := gUp.Eq(m)
	if err != nil {
		return nil, err
	}
	r1 := CrossRuleResult{Name: "wavelength increase raises |theta|"}
	if resBase.Exists && resUp.Exists {
		r1.Holds = math.Abs(resUp.AngleRad) > math.Abs(resBase.AngleRad)
		r1.Detail = f2(math.Abs(resBase.AngleRad)) + " -> " + f2(math.Abs(resUp.AngleRad)) + " rad"
	} else {
		r1.Detail = "order leaves the visible region during the probe"
	}
	out = append(out, r1)

	// Rule 2: denser grating -> |theta| up.
	gDense := g
	gDense.GrooveSpacingNm *= 0.9
	resDense, err := gDense.Eq(m)
	if err != nil {
		return nil, err
	}
	r2 := CrossRuleResult{Name: "denser grating raises |theta|"}
	if resBase.Exists && resDense.Exists {
		r2.Holds = math.Abs(resDense.AngleRad) > math.Abs(resBase.AngleRad)
		r2.Detail = f2(math.Abs(resBase.AngleRad)) + " -> " + f2(math.Abs(resDense.AngleRad)) + " rad"
	} else {
		r2.Detail = "order leaves the visible region during the probe"
	}
	out = append(out, r2)

	// Rule 3: zero order along the incident direction.
	res0, err := g.Eq(0)
	if err != nil {
		return nil, err
	}
	r3 := CrossRuleResult{Name: "zero order follows the incident direction"}
	r3.Holds = res0.Exists && math.Abs(res0.AngleRad-g.IncidentAngleRad) < 1e-9
	r3.Detail = "theta_0 = " + f2(res0.AngleRad) + " rad vs theta_i = " + f2(g.IncidentAngleRad) + " rad"
	out = append(out, r3)

	// Rule 4: plus/minus symmetry at normal incidence.
	r4 := CrossRuleResult{Name: "plus/minus order symmetry at normal incidence"}
	if math.Abs(g.IncidentAngleRad) < 1e-12 && math.Abs(float64(m)) > 0 {
		pair, err := g.SymmetricPair(m)
		if err == nil {
			r4.Holds = pair.Antisymmetric()
			r4.Detail = "theta_+" + itoa(m) + "=" + f2(pair.Plus.AngleRad) +
				" theta_-" + itoa(m) + "=" + f2(pair.Minus.AngleRad) + " rad"
		} else {
			r4.Detail = err.Error()
		}
	} else {
		r4.Holds = true
		r4.Detail = "not applicable at non-normal incidence"
	}
	out = append(out, r4)

	// Rule 5: R linear in N.
	r5 := CrossRuleResult{Name: "resolving power grows linearly with N"}
	if math.Abs(float64(m)) > 0 && n >= 1 {
		r1, err := ResolvingPower(g.WavelengthNm, m, n)
		if err == nil {
			r2, err := ResolvingPower(g.WavelengthNm, m, 2*n)
			if err == nil {
				r5.Holds = math.Abs(r2.ResolvingPower-2*r1.ResolvingPower) < 1e-9
				r5.Detail = "R(" + itoa(n) + ")=" + f2(r1.ResolvingPower) +
					" R(" + itoa(2*n) + ")=" + f2(r2.ResolvingPower)
			} else {
				r5.Detail = err.Error()
			}
		} else {
			r5.Detail = err.Error()
		}
	} else {
		r5.Holds = true
		r5.Detail = "not applicable (m=0 or no finite N)"
	}
	out = append(out, r5)

	return out, nil
}

func f2(v float64) string {
	return formatValue(v)
}
