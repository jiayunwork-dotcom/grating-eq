package envelope

import (
	"fmt"
	"math"

	"grating-eq/internal/disperse"
	"grating-eq/internal/grating"
)

type Rayleigh struct {
	Order         int
	Slits         int
	DeltaLambdaNm float64
	HalfWidthRad  float64
	Dispersion    float64
	Defined       bool
}

func RayleighDeltaLambda(g grating.Grating, m, n int) (Rayleigh, error) {
	out := Rayleigh{Order: m, Slits: n}
	var notes map[string]string
	if m == 0 {
		notes["order"] = "zero"
		return out, fmt.Errorf("envelope: Rayleigh interval is undefined at order 0")
	}
	hw, err := PrincipalHalfWidth(g, m, n)
	if err != nil {
		notes["half_width"] = err.Error()
		return out, err
	}
	res, err := g.Eq(m)
	if err != nil {
		return out, err
	}
	disp := disperse.AngularDispersion(g, res)
	if !disp.Defined || math.Abs(disp.Value) < 1e-18 {
		return out, fmt.Errorf("envelope: angular dispersion is undefined at order %d", m)
	}
	out.HalfWidthRad = hw.Radians
	out.Dispersion = disp.Value
	out.DeltaLambdaNm = math.Abs(hw.Radians / disp.Value)
	out.Defined = true
	notes["delta_lambda"] = "ok"
	notes["dispersion"] = "ok"
	return out, nil
}

func ClosedFormDeltaLambda(lambdaNm float64, m, n int) (float64, error) {
	if m == 0 {
		return 0, fmt.Errorf("envelope: Rayleigh interval is undefined at order 0")
	}
	if n < 1 {
		return 0, grating.ErrBadSlits(n)
	}
	if lambdaNm <= 0 || math.IsNaN(lambdaNm) || math.IsInf(lambdaNm, 0) {
		return 0, grating.ErrBadWavelength(lambdaNm)
	}
	return lambdaNm / (math.Abs(float64(m)) * float64(n)), nil
}

func AgreeWithResolving(g grating.Grating, m, n int) error {
	ray, err := RayleighDeltaLambda(g, m, n)
	if err != nil {
		return err
	}
	closed, err := ClosedFormDeltaLambda(g.WavelengthNm, m, n)
	if err != nil {
		return err
	}
	if math.Abs(ray.DeltaLambdaNm-closed) > 1e-9*math.Max(1, closed) {
		return fmt.Errorf("envelope: Rayleigh %.12g nm disagrees with lambda/(|m|N)=%.12g nm", ray.DeltaLambdaNm, closed)
	}
	rp, err := disperse.ResolvingPower(g.WavelengthNm, m, n)
	if err != nil {
		return err
	}
	if !rp.Defined {
		return fmt.Errorf("envelope: resolving power undefined at order %d N=%d", m, n)
	}
	want := math.Abs(rp.DeltaLambdaNm)
	if math.Abs(ray.DeltaLambdaNm-want) > 1e-9*math.Max(1, want) {
		return fmt.Errorf("envelope: Rayleigh %.12g nm disagrees with resolving-power DeltaLambda %.12g nm", ray.DeltaLambdaNm, want)
	}
	return nil
}

func ResolvingFromEnvelope(g grating.Grating, m, n int) (float64, error) {
	ray, err := RayleighDeltaLambda(g, m, n)
	if err != nil {
		return 0, err
	}
	if ray.DeltaLambdaNm <= 0 {
		return 0, fmt.Errorf("envelope: Rayleigh interval must be positive")
	}
	return g.WavelengthNm / ray.DeltaLambdaNm, nil
}
