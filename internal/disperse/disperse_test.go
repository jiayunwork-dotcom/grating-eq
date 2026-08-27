package disperse

import (
	"math"
	"testing"

	"grating-eq/internal/grating"
)

func fixture600() grating.Grating {
	g, err := grating.New(grating.SpacingFromDensity(600), 550, 0, 0)
	if err != nil {
		panic(err)
	}
	return g
}

func TestAngularDispersionClosedForm(t *testing.T) {
	g := fixture600()
	res, err := g.Eq(1)
	if err != nil {
		t.Fatalf("Eq(1) error: %v", err)
	}
	d := AngularDispersion(g, res)
	if !d.Defined {
		t.Fatal("dispersion should be defined for order 1")
	}
	cos := math.Cos(res.AngleRad)
	want := 1.0 / (g.GrooveSpacingNm * cos)
	if math.Abs(d.Value-want) > 1e-12 {
		t.Errorf("dtheta/dlambda = %v, want %v", d.Value, want)
	}
}

func TestAngularDispersionSign(t *testing.T) {
	g := fixture600()
	for _, m := range []int{1, 2, 3} {
		res, err := g.Eq(m)
		if err != nil {
			t.Fatalf("Eq(%d) error: %v", m, err)
		}
		d := AngularDispersion(g, res)
		if !d.Defined {
			t.Fatalf("dispersion undefined for order %d", m)
		}
		if d.SignOf() != 1 {
			t.Errorf("order %d: dispersion sign = %d, want +1 under transmission convention",
				m, d.SignOf())
		}
	}
	for _, m := range []int{-1, -2, -3} {
		res, err := g.Eq(m)
		if err != nil {
			t.Fatalf("Eq(%d) error: %v", m, err)
		}
		d := AngularDispersion(g, res)
		if !d.Defined {
			t.Fatalf("dispersion undefined for order %d", m)
		}
		if d.SignOf() != -1 {
			t.Errorf("order %d: dispersion sign = %d, want -1", m, d.SignOf())
		}
	}
}

func TestFSRDefinition(t *testing.T) {
	g := fixture600()
	for _, m := range []int{1, -1, 2, -2, 3} {
		r := FreeSpectralRange(g.WavelengthNm, m)
		if !r.Defined {
			t.Fatalf("FSR should be defined for m=%d", m)
		}
		want := g.WavelengthNm / math.Abs(float64(m))
		if math.Abs(r.ValueNm-want) > 1e-12 {
			t.Errorf("FSR(m=%d) = %v, want %v", m, r.ValueNm, want)
		}
	}
	if r := FreeSpectralRange(g.WavelengthNm, 0); r.Defined {
		t.Error("zero order must not define an FSR")
	}
}

func TestResolvingPowerLinearInN(t *testing.T) {
	g := fixture600()
	r1, err := ResolvingPower(g.WavelengthNm, 2, 1000)
	if err != nil {
		t.Fatalf("ResolvingPower error: %v", err)
	}
	r2, err := ResolvingPower(g.WavelengthNm, 2, 2000)
	if err != nil {
		t.Fatalf("ResolvingPower error: %v", err)
	}
	if r2.ResolvingPower != 2*r1.ResolvingPower {
		t.Errorf("R(N=2000)=%v must equal 2*R(N=1000)=%v",
			r2.ResolvingPower, 2*r1.ResolvingPower)
	}
	if math.Abs(r1.DeltaLambdaNm-(g.WavelengthNm/2000)) > 1e-9 {
		t.Errorf("DeltaLambda = %v, want %v", r1.DeltaLambdaNm, g.WavelengthNm/2000)
	}
}

func TestResolvingPowerRejectsBadN(t *testing.T) {
	if _, err := ResolvingPower(550, 1, 0); err == nil {
		t.Error("N=0 must be rejected")
	}
	if _, err := ResolvingPower(550, 1, -3); err == nil {
		t.Error("N=-3 must be rejected")
	}
}

func TestEffectiveSlits(t *testing.T) {
	n, err := EffectiveSlits(2, 4000)
	if err != nil {
		t.Fatalf("EffectiveSlits error: %v", err)
	}
	if n != 2000 {
		t.Errorf("EffectiveSlits(2, 4000) = %d, want 2000", n)
	}
}

func TestCriticalDensity(t *testing.T) {
	d, err := CriticalDensity(550, 1)
	if err != nil {
		t.Fatalf("CriticalDensity error: %v", err)
	}
	if math.Abs(d-(grating.NanoPerMilli/550)) > 1e-6 {
		t.Errorf("critical density = %v, want %v", d, grating.NanoPerMilli/550)
	}
}

func TestVanishAnalysis(t *testing.T) {
	g := fixture600()
	rep, err := VanishAnalysis(g, 16)
	if err != nil {
		t.Fatalf("VanishAnalysis error: %v", err)
	}
	if rep.MaxOrder != 3 {
		t.Errorf("max surviving order = %d, want 3", rep.MaxOrder)
	}
	if rep.FirstMissing != 4 {
		t.Errorf("first missing order = %d, want 4", rep.FirstMissing)
	}
	gg, err := grating.New(grating.SpacingFromDensity(3000), 550, 0, 0)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	rep2, err := VanishAnalysis(gg, 8)
	if err != nil {
		t.Fatalf("VanishAnalysis error: %v", err)
	}
	if rep2.Existing != 0 || rep2.MaxOrder != 0 {
		t.Errorf("dense grating must hide all non-zero orders, got existing=%d max=%d",
			rep2.Existing, rep2.MaxOrder)
	}
}

func TestCrossRulesHold(t *testing.T) {
	g := fixture600()
	rules, err := CrossRules(g, 1, 1000)
	if err != nil {
		t.Fatalf("CrossRules error: %v", err)
	}
	if len(rules) != 5 {
		t.Fatalf("CrossRules returned %d rules, want 5", len(rules))
	}
	for _, r := range rules {
		if !r.Holds {
			t.Errorf("cross rule %q failed: %s", r.Name, r.Detail)
		}
	}
}

func TestScanOrderFade(t *testing.T) {
	g := fixture600()
	points, err := ScanOrder(g, 1, 21, 100, 3000)
	if err != nil {
		t.Fatalf("ScanOrder error: %v", err)
	}
	faded := false
	for _, p := range points {
		if !p.Exists {
			faded = true
		}
	}
	if !faded {
		t.Error("scan across lambda/d > 1 must show the order fading out")
	}
}

func TestInverseScan(t *testing.T) {
	g := fixture600()
	pts, err := InverseScan(g, 1, 21)
	if err != nil {
		t.Fatalf("InverseScan error: %v", err)
	}
	if len(pts) == 0 {
		t.Fatal("inverse scan returned no points")
	}
	found := false
	for _, p := range pts {
		if math.Abs(p.AngleRad-math.Asin(550/g.GrooveSpacingNm)) < 0.1 {
			if math.Abs(p.WavelengthNm-550) > 60 {
				t.Errorf("inverse wavelength at theta_1 = %v nm, want ~550", p.WavelengthNm)
			}
			found = true
		}
	}
	if !found {
		t.Error("inverse scan never sampled the example wavelength")
	}
}
