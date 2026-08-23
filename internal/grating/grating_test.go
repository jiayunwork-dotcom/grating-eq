package grating

import (
	"math"
	"testing"
)

// fixture600 returns the bundled example geometry: 600 grooves/mm,
// lambda = 550 nm, normal incidence.
func fixture600() Grating {
	d := SpacingFromDensity(600)
	g, err := New(d, 550, 0, 0)
	if err != nil {
		panic(err)
	}
	return g
}

// TestFirstOrderNormalIncidence pins the closed form of the +-1 orders at
// normal incidence: sin(theta_1) = lambda/d, so theta_1 = asin(lambda/d).
func TestFirstOrderNormalIncidence(t *testing.T) {
	g := fixture600()
	want := math.Asin(g.WavelengthNm / g.GrooveSpacingNm)
	for _, m := range []int{1, -1} {
		res, err := g.Eq(m)
		if err != nil {
			t.Fatalf("Eq(%d) error: %v", m, err)
		}
		if !res.Exists {
			t.Fatalf("order %d should exist for the example geometry", m)
		}
		if got := math.Abs(res.AngleRad); math.Abs(got-want) > 1e-9 {
			t.Errorf("order %d: |theta| = %v, want %v (asin(lambda/d))", m, got, want)
		}
	}
}

// TestZeroOrderFollowsIncident checks that m=0 always continues along the
// incident direction, including at non-zero incidence.
func TestZeroOrderFollowsIncident(t *testing.T) {
	inc := 0.4 // radians
	g, err := New(1666.7, 550, inc, 0)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	res, err := g.Eq(0)
	if err != nil {
		t.Fatalf("Eq(0) error: %v", err)
	}
	if !res.Exists {
		t.Fatal("zero order must always exist")
	}
	if math.Abs(res.AngleRad-inc) > 1e-9 {
		t.Errorf("zero order angle = %v, want incident angle %v", res.AngleRad, inc)
	}
}

// TestTransmissionConventionSignSymmetric guards the transmission sign
// convention: at normal incidence the pair +m / -m must be exactly
// antisymmetric. Writing the equation as m*lambda = d*(sin+sin) breaks
// this symmetry, so the test fails for that implementation.
func TestTransmissionConventionSignSymmetric(t *testing.T) {
	g := fixture600()
	pair, err := g.SymmetricPair(1)
	if err != nil {
		t.Fatalf("SymmetricPair error: %v", err)
	}
	if !pair.Antisymmetric() {
		t.Errorf(
			"normal-incidence +-1 orders must be antisymmetric: "+
				"theta_+1=%v theta_-1=%v sum=%v",
			pair.Plus.AngleRad, pair.Minus.AngleRad, pair.Plus.AngleRad+pair.Minus.AngleRad)
	}
	if pair.Plus.SineOfAngle+pair.Minus.SineOfAngle > 1e-12 {
		t.Errorf("sin values must be odd in m, got %v and %v",
			pair.Plus.SineOfAngle, pair.Minus.SineOfAngle)
	}
}

// TestDenseGratingHighOrderVanishes checks that when lambda/d becomes
// large enough, high orders exceed |sin| <= 1 and disappear.
func TestDenseGratingHighOrderVanishes(t *testing.T) {
	// 3000 grooves/mm at 550 nm: lambda/d = 1.65, only m=0 survives.
	d := SpacingFromDensity(3000)
	g, err := New(d, 550, 0, 0)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	for m := -3; m <= 3; m++ {
		res, err := g.Eq(m)
		if err != nil {
			t.Fatalf("Eq(%d) error: %v", m, err)
		}
		if m == 0 && !res.Exists {
			t.Error("zero order must exist even on a dense grating")
		}
		if m != 0 && res.Exists {
			t.Errorf("order %d must vanish when lambda/d=%.3f > 1", m, 550/d)
		}
	}
}

// TestValidateRejectsZeroSpacing checks the simplest failure surface: a
// groove spacing of zero is rejected with an error.
func TestValidateRejectsZeroSpacing(t *testing.T) {
	_, err := New(0, 550, 0, 0)
	if err == nil {
		t.Fatal("d=0 must be rejected")
	}
}

// TestValidateRejectsNegativeWavelength checks that a negative wavelength
// is rejected with an error.
func TestValidateRejectsNegativeWavelength(t *testing.T) {
	_, err := New(1666.7, -100, 0, 0)
	if err == nil {
		t.Fatal("negative wavelength must be rejected")
	}
}

// TestValidateRejectsInvalidIncidentSine checks that an incident angle
// whose sine is not a real number is rejected.
func TestValidateRejectsInvalidIncidentSine(t *testing.T) {
	_, err := New(1666.7, 550, math.Inf(1), 0)
	if err == nil {
		t.Fatal("non-finite incident angle must be rejected")
	}
	_, err = New(1666.7, 550, math.NaN(), 0)
	if err == nil {
		t.Fatal("NaN incident angle must be rejected")
	}
}

// TestLineDensityReciprocal checks that d and the line density are exact
// reciprocals.
func TestLineDensityReciprocal(t *testing.T) {
	d := SpacingFromDensity(600)
	if got := NanoPerMilli / d; math.Abs(got-600) > 1e-9 {
		t.Errorf("reciprocal density = %v, want 600", got)
	}
	g := fixture600()
	if got := g.GrooveDensity(); math.Abs(got-600) > 1e-9 {
		t.Errorf("GrooveDensity = %v, want 600", got)
	}
}

// TestDensityConsistencyRejectsMismatch checks that supplying a spacing and
// a density that disagree is reported as an error.
func TestDensityConsistencyRejectsMismatch(t *testing.T) {
	d := SpacingFromDensity(600)
	v, err := DensityOf(d)
	if err != nil {
		t.Fatalf("DensityOf error: %v", err)
	}
	if err := v.CheckConsistency(500); err == nil {
		t.Error("conflicting density must be rejected")
	}
	if err := v.CheckConsistency(600); err != nil {
		t.Errorf("consistent density rejected: %v", err)
	}
}

// TestWavelengthIncreaseRaisesAngle checks the cross rule that raising the
// wavelength at a fixed order increases |theta|.
func TestWavelengthIncreaseRaisesAngle(t *testing.T) {
	g := fixture600()
	a, err := g.Eq(1)
	if err != nil {
		t.Fatalf("Eq error: %v", err)
	}
	g.WavelengthNm = 600
	b, err := g.Eq(1)
	if err != nil {
		t.Fatalf("Eq error: %v", err)
	}
	if !(math.Abs(b.AngleRad) > math.Abs(a.AngleRad)) {
		t.Errorf("|theta| must grow from %.6g to %.6g when lambda goes 550->600",
			math.Abs(a.AngleRad), math.Abs(b.AngleRad))
	}
}

// TestDenserGratingRaisesAngle checks the cross rule that a smaller d
// (denser grating) pushes the same order further from the normal.
func TestDenserGratingRaisesAngle(t *testing.T) {
	g := fixture600()
	a, err := g.Eq(1)
	if err != nil {
		t.Fatalf("Eq error: %v", err)
	}
	g.GrooveSpacingNm = 1200
	b, err := g.Eq(1)
	if err != nil {
		t.Fatalf("Eq error: %v", err)
	}
	if !(math.Abs(b.AngleRad) > math.Abs(a.AngleRad)) {
		t.Errorf("|theta| must grow when d shrinks: %.6g -> %.6g",
			math.Abs(a.AngleRad), math.Abs(b.AngleRad))
	}
}

// TestMaxVisibleOrder checks the analytic bound on the highest order.
func TestMaxVisibleOrder(t *testing.T) {
	g := fixture600()
	// lambda/d = 0.33, so |m| <= 3.
	if got := g.MaxVisibleOrder(); got != 3 {
		t.Errorf("MaxVisibleOrder = %d, want 3", got)
	}
	// Very coarse grating: lambda/d tiny, many orders exist, capped by scan.
	g.GrooveSpacingNm = 50000
	if got := g.MaxVisibleOrder(); got != DefaultMaxOrderScan {
		t.Errorf("MaxVisibleOrder on coarse grating = %d, want scan cap %d",
			got, DefaultMaxOrderScan)
	}
}

// TestArcSineDomain guards the domain check that must precede asin.
func TestArcSineDomain(t *testing.T) {
	if _, err := ArcSine(1.5); err == nil {
		t.Error("asin(1.5) must be rejected")
	}
	if _, err := ArcSine(-1.0000001); err == nil {
		t.Error("asin(-1.0000001) must be rejected")
	}
	v, err := ArcSine(1)
	if err != nil || math.Abs(v-math.Pi/2) > 1e-12 {
		t.Errorf("asin(1) = %v, err %v; want pi/2", v, err)
	}
}

// TestZeroOrderExistsEverywhere checks that m=0 is never hidden.
func TestZeroOrderExistsEverywhere(t *testing.T) {
	densities := []float64{100, 600, 2400}
	wavelengths := []float64{100, 550, 2500}
	incs := []float64{-0.5, 0, 0.6}
	for _, dd := range densities {
		for _, wl := range wavelengths {
			for _, inc := range incs {
				g, err := New(SpacingFromDensity(dd), wl, inc, 0)
				if err != nil {
					t.Fatalf("New(%v,%v,%v) error: %v", dd, wl, inc, err)
				}
				res, err := g.Eq(0)
				if err != nil {
					t.Fatalf("Eq(0) error: %v", err)
				}
				if !res.Exists {
					t.Errorf("zero order missing for d=%v lambda=%v inc=%v", dd, wl, inc)
				}
			}
		}
	}
}

// TestLookupRoundTrip checks that LookupAngle and WavelengthForOrder are
// inverses of each other for a visible order.
func TestLookupRoundTrip(t *testing.T) {
	g := fixture600()
	res, err := g.Eq(2)
	if err != nil {
		t.Fatalf("Eq(2) error: %v", err)
	}
	if !res.Exists {
		t.Fatal("order 2 should exist for the example geometry")
	}
	back, err := WavelengthForOrder(g.GrooveSpacingNm, 0, res.AngleRad, 2)
	if err != nil {
		t.Fatalf("WavelengthForOrder error: %v", err)
	}
	if math.Abs(back-g.WavelengthNm) > 1e-6 {
		t.Errorf("round trip wavelength = %v, want %v", back, g.WavelengthNm)
	}
}
