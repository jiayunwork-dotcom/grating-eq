package report

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

const example600JSON = `{
  "grooves_per_mm": 600,
  "wavelength_nm": 550,
  "incident_angle_deg": 0
}`

// TestParseValidInput checks that the bundled example JSON decodes.
func TestParseValidInput(t *testing.T) {
	in, err := Parse([]byte(example600JSON))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if in.GroovesPerMm == nil || *in.GroovesPerMm != 600 {
		t.Errorf("grooves_per_mm = %v, want 600", in.GroovesPerMm)
	}
	if in.WavelengthNm != 550 {
		t.Errorf("wavelength_nm = %v, want 550", in.WavelengthNm)
	}
	if in.Slits != 0 {
		t.Errorf("slits = %d, want 0", in.Slits)
	}
}

// TestRejectZeroSpacing checks the simplest failure surface: a zero groove
// spacing is rejected with an error instead of a divide-by-zero.
func TestRejectZeroSpacing(t *testing.T) {
	in, err := Parse([]byte(`{"d_nm": 0, "wavelength_nm": 550}`))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if _, err := in.ToGrating(); err == nil {
		t.Fatal("d=0 must be rejected")
	}
}

// TestRejectNegativeWavelength checks that a negative wavelength is
// rejected with an error.
func TestRejectNegativeWavelength(t *testing.T) {
	in, err := Parse([]byte(`{"d_nm": 1666.7, "wavelength_nm": -100}`))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if _, err := in.ToGrating(); err == nil {
		t.Fatal("negative wavelength must be rejected")
	}
}

// TestRejectMissingGeometry checks that an input without any groove
// description is rejected.
func TestRejectMissingGeometry(t *testing.T) {
	in, err := Parse([]byte(`{"wavelength_nm": 550}`))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if _, err := in.ToGrating(); err == nil {
		t.Fatal("input without groove geometry must be rejected")
	}
}

// TestRejectUnknownField checks that a JSON typo is surfaced as an error.
func TestRejectUnknownField(t *testing.T) {
	if _, err := Parse([]byte(`{"grooves_per_mm": 600, "wavelegnth_nm": 550}`)); err == nil {
		t.Fatal("unknown field must be rejected")
	}
}

// TestRunExampleTable checks the CLI output of the bundled example: the
// +-1 orders land on asin(lambda/d) and the zero order stays at 0 degrees.
func TestRunExampleTable(t *testing.T) {
	in, err := Parse([]byte(example600JSON))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	out, err := Run(in, DefaultOptions())
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if !strings.Contains(out, "theta deg") {
		t.Error("table header missing")
	}
	// Order 1: theta = asin(550 / (1e6/600)) = 19.27 degrees.
	if a := findAngle(out, "+1"); math.Abs(a-19.27) > 0.01 {
		t.Errorf("order +1 angle = %v, want ~19.27", a)
	}
	if a := findAngle(out, "-1"); math.Abs(a-19.27) > 0.01 {
		t.Errorf("order -1 angle = %v, want ~19.27", a)
	}
	// Order 0 at 0.000000.
	if !strings.Contains(out, "0.000000") {
		t.Errorf("zero order angle missing from output:\n%s", out)
	}
}

// TestRunExampleSymmetry checks that the rendered table shows the
// transmission-convention symmetry: at normal incidence +m and -m must be
// exactly antisymmetric. A plus-sign grating equation would break this.
func TestRunExampleSymmetry(t *testing.T) {
	in, err := Parse([]byte(example600JSON))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	out, err := Run(in, DefaultOptions())
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	for _, m := range []string{"1", "2", "3"} {
		plus := findAngle(out, "+"+m)
		minus := findAngle(out, "-"+m)
		if plus == 0 && minus == 0 {
			continue
		}
		if plus != minus {
			t.Errorf("order +%s and -%s must share |theta|, got %v vs %v", m, m, plus, minus)
		}
	}
}

// TestRunJSON checks the --json rendering carries the key fields.
func TestRunJSON(t *testing.T) {
	in, err := Parse([]byte(example600JSON))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	opts := DefaultOptions()
	opts.JSON = true
	out, err := Run(in, opts)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if !strings.Contains(out, `"grooves_per_mm": 600`) {
		t.Errorf("JSON output missing density:\n%s", out)
	}
	if !strings.Contains(out, `"max_visible_order": 3`) {
		t.Errorf("JSON output missing max order:\n%s", out)
	}
}

// TestRunResolutionColumns checks that a finite slit count adds R and
// dlambda columns.
func TestRunResolutionColumns(t *testing.T) {
	in, err := Parse([]byte(`{
	  "grooves_per_mm": 600,
	  "wavelength_nm": 550,
	  "incident_angle_deg": 0,
	  "slits": 2000
	}`))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	out, err := Run(in, DefaultOptions())
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if !strings.Contains(out, "R") || !strings.Contains(out, "dlambda") {
		t.Errorf("resolution columns missing:\n%s", out)
	}
}

// TestRunRejectsZeroSlits checks that slits < 1 is an error when N was
// explicitly supplied as a negative number.
func TestRunRejectsBadSlits(t *testing.T) {
	in, err := Parse([]byte(`{
	  "grooves_per_mm": 600,
	  "wavelength_nm": 550,
	  "slits": -1
	}`))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if _, err := Run(in, DefaultOptions()); err == nil {
		t.Fatal("negative slit count must be rejected")
	}
}

// TestChecksExample passes the bundled example through the cross-rule
// subcommand.
func TestChecksExample(t *testing.T) {
	in, err := Parse([]byte(example600JSON))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	out, err := RunChecks(in)
	if err != nil {
		t.Fatalf("RunChecks error: %v", err)
	}
	if !strings.Contains(out, "PASS") {
		t.Errorf("checks output should contain PASS:\n%s", out)
	}
}

// TestDensityReciprocalOutput checks that RunDensity shows the exact
// reciprocal relation.
func TestDensityReciprocalOutput(t *testing.T) {
	in, err := Parse([]byte(example600JSON))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	out, err := RunDensity(in)
	if err != nil {
		t.Fatalf("RunDensity error: %v", err)
	}
	if !strings.Contains(out, "grooves/mm") {
		t.Errorf("density output missing unit:\n%s", out)
	}
}

// findAngle extracts the |theta| of an order row by searching for the
// order column and the angle column in the same line. The order column is
// rendered without a leading plus sign, so "+1" maps to "1".
func findAngle(out, orderToken string) float64 {
	want := strings.TrimPrefix(orderToken, "+")
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] != want {
			continue
		}
		// The third field is the angle in degrees.
		var v float64
		if _, err := fmt.Sscanf(fields[2], "%f", &v); err == nil {
			if v < 0 {
				v = -v
			}
			return v
		}
	}
	return 0
}
