package disperse

import (
	"fmt"
	"math"
)

// ResolvingResult carries the resolving power and the minimum resolvable
// wavelength interval of one order on a grating with N illuminated slits.
type ResolvingResult struct {
	// Order is the diffraction order.
	Order int

	// Slits is the finite number of illuminated grooves.
	Slits int

	// ResolvingPower is R = m*N.
	ResolvingPower float64

	// DeltaLambdaNm is lambda/R, the smallest resolvable interval.
	DeltaLambdaNm float64

	// Defined reports whether R is positive and finite.
	Defined bool
}

// ResolvingPower computes R = m*N and the minimum resolvable interval
// DeltaLambda = lambda/R. The linear growth of R with N is an explicit
// cross rule of the tool.
func ResolvingPower(lambdaNm float64, m, n int) (ResolvingResult, error) {
	out := ResolvingResult{Order: m, Slits: n}
	if m == 0 {
		return out, nil
	}
	if n < 1 {
		return out, fmt.Errorf("finite slit count N must be >= 1 when given, got %d", n)
	}
	if lambdaNm <= 0 {
		return out, fmt.Errorf("wavelength must be positive, got %v", lambdaNm)
	}
	r := float64(m) * float64(n)
	out.ResolvingPower = r
	out.DeltaLambdaNm = lambdaNm / r
	out.Defined = true
	return out, nil
}

// MinimumResolvable returns DeltaLambda for a resolving power R.
func MinimumResolvable(lambdaNm, resolvingPower float64) (float64, error) {
	if resolvingPower <= 0 {
		return 0, fmt.Errorf("resolving power must be positive, got %v", resolvingPower)
	}
	if lambdaNm <= 0 {
		return 0, fmt.Errorf("wavelength must be positive, got %v", lambdaNm)
	}
	return lambdaNm / resolvingPower, nil
}

// ResolutionString renders R and DeltaLambda compactly.
func (r ResolvingResult) ResolutionString() string {
	if !r.Defined {
		return "—"
	}
	return fmt.Sprintf("R=%.4g  dlambda=%.6g nm",
		r.ResolvingPower, r.DeltaLambdaNm)
}

// DeltaLambdaAtR computes the resolvable interval for a hypothetical
// resolving power, used for "what if" queries.
func DeltaLambdaAtR(lambdaNm, r float64) (float64, error) {
	return MinimumResolvable(lambdaNm, r)
}

// EffectiveSlits returns the number of slits needed to reach a target
// resolving power R at order m.
func EffectiveSlits(m int, targetR float64) (int, error) {
	if m == 0 {
		return 0, fmt.Errorf("zero order has no resolving power")
	}
	if targetR <= 0 {
		return 0, fmt.Errorf("target resolving power must be positive, got %v", targetR)
	}
	n := int(math.Ceil(targetR / math.Abs(float64(m))))
	if n < 1 {
		n = 1
	}
	return n, nil
}
