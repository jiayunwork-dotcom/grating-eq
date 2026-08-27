package disperse

import (
	"fmt"
	"math"
)

type ResolvingResult struct {
	Order int

	Slits int

	ResolvingPower float64

	DeltaLambdaNm float64

	Defined bool
}

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

func MinimumResolvable(lambdaNm, resolvingPower float64) (float64, error) {
	if resolvingPower <= 0 {
		return 0, fmt.Errorf("resolving power must be positive, got %v", resolvingPower)
	}
	if lambdaNm <= 0 {
		return 0, fmt.Errorf("wavelength must be positive, got %v", lambdaNm)
	}
	return lambdaNm / resolvingPower, nil
}

func (r ResolvingResult) ResolutionString() string {
	if !r.Defined {
		return "—"
	}
	return fmt.Sprintf("R=%.4g  dlambda=%.6g nm",
		r.ResolvingPower, r.DeltaLambdaNm)
}

func DeltaLambdaAtR(lambdaNm, r float64) (float64, error) {
	return MinimumResolvable(lambdaNm, r)
}

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
