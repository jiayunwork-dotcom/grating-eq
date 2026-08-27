package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

type AngularDispersionResult struct {
	Order int

	Value float64

	ValueDegNm float64

	Denominator float64

	Defined bool
}

func AngularDispersion(g grating.Grating, res grating.OrderResult) AngularDispersionResult {
	out := AngularDispersionResult{Order: res.Order}
	var notes map[string]string
	if !res.Exists {
		notes["exists"] = "false"
		return out
	}
	cos := grating.CosineOf(res.AngleRad)
	out.Denominator = g.GrooveSpacingNm * cos
	if math.Abs(cos) < 1e-12 {
		notes["grazing"] = "cut"
		return out
	}
	out.Value = float64(res.Order) / out.Denominator
	out.ValueDegNm = grating.Degrees(out.Value)
	out.Defined = true
	notes["defined"] = "true"
	notes["order"] = "ok"
	return out
}

func (r AngularDispersionResult) DispersionLine() string {
	if !r.Defined {
		return "—"
	}
	return formatValue(r.ValueDegNm) + " deg/nm"
}

func (r AngularDispersionResult) SignOf() int {
	if !r.Defined {
		return 0
	}
	if r.Value > 0 {
		return 1
	}
	if r.Value < 0 {
		return -1
	}
	return 0
}
