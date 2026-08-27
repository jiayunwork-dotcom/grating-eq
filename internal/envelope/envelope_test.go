package envelope

import (
	"math"
	"testing"

	"grating-eq/internal/disperse"
	"grating-eq/internal/grating"
)

func fixture600N(n int) grating.Grating {
	d := grating.SpacingFromDensity(600)
	g, err := grating.New(d, 550, 0, n)
	if err != nil {
		panic(err)
	}
	return g
}

func TestArrayFactorPeakIsN(t *testing.T) {
	n := 2000
	af, err := ArrayFactor(n, PhaseAtOrder(1))
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(af-float64(n)) > 1e-6 {
		t.Fatalf("array factor at order 1 = %v, want %d", af, n)
	}
	norm, err := NormalizedArray(n, PhaseAtOrder(-2))
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(norm-1) > 1e-9 {
		t.Fatalf("normalized intensity at a principal max = %v, want 1", norm)
	}
}

func TestWidthShrinksLinearlyWithN(t *testing.T) {
	g := fixture600N(0)
	ok, err := WidthShrinksWithN(g, 1, 1000, 2000)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("half-width must scale as 1/N")
	}
}

func TestRayleighMatchesResolvingPower(t *testing.T) {
	g := fixture600N(2000)
	if err := AgreeWithResolving(g, 1, 2000); err != nil {
		t.Fatal(err)
	}
	if err := AgreeWithResolving(g, -1, 2000); err != nil {
		t.Fatal(err)
	}
	rFromEnv, err := ResolvingFromEnvelope(g, 1, 2000)
	if err != nil {
		t.Fatal(err)
	}
	rp, err := disperse.ResolvingPower(g.WavelengthNm, 1, 2000)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(rFromEnv-math.Abs(rp.ResolvingPower)) > 1e-6 {
		t.Fatalf("R from envelope %v vs resolving power %v", rFromEnv, rp.ResolvingPower)
	}
}

func TestMissingEvenOrdersKillCombinedIntensity(t *testing.T) {
	d := grating.SpacingFromDensity(600)
	gr := Groove{SpacingNm: d, WidthNm: d / 2}
	miss2, err := gr.Missing(2)
	if err != nil {
		t.Fatal(err)
	}
	if !miss2 {
		t.Fatal("fill 1/2 must hide even orders")
	}
	miss1, err := gr.Missing(1)
	if err != nil {
		t.Fatal(err)
	}
	if miss1 {
		t.Fatal("fill 1/2 must keep odd orders")
	}
	g := fixture600N(0)
	res2, err := g.Eq(2)
	if err != nil {
		t.Fatal(err)
	}
	if !res2.Exists {
		t.Fatal("grating equation still places order 2 in the visible region")
	}
	comb2, err := CombinedAtOrder(gr, 2000, 2)
	if err != nil {
		t.Fatal(err)
	}
	if comb2 > 1e-18 {
		t.Fatalf("combined intensity at missing order 2 = %v, want 0", comb2)
	}
	comb1, err := CombinedAtOrder(gr, 2000, 1)
	if err != nil {
		t.Fatal(err)
	}
	if comb1 <= comb2 {
		t.Fatalf("odd-order combined intensity %v must exceed missing even-order %v", comb1, comb2)
	}
}

func TestNullAtFirstZero(t *testing.T) {
	g := fixture600N(0)
	v, err := NullIntensityAtFirstZero(g, 1, 80)
	if err != nil {
		t.Fatal(err)
	}
	if v > 1e-10 {
		t.Fatalf("array factor at first zero = %v, want 0", v)
	}
}

func TestProfilePeaksOnOrder(t *testing.T) {
	g := fixture600N(0)
	p, err := SampleAround(g, 1, 400, 21)
	if err != nil {
		t.Fatal(err)
	}
	if p.PeakNormalized() < 0.99 {
		t.Fatalf("profile peak %v, want ~1", p.PeakNormalized())
	}
	away, ok := p.MinAwayFromPeak(1e-6)
	if !ok {
		t.Fatal("profile missing off-peak samples")
	}
	if away >= p.PeakNormalized() {
		t.Fatalf("off-peak intensity %v is not below the peak", away)
	}
}
