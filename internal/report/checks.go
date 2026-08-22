package report

import (
	"fmt"
	"strings"

	"grating-eq/internal/disperse"
	"grating-eq/internal/grating"
)

// RunChecks executes the checks subcommand: it evaluates the five cross
// rules of the grating model on the supplied geometry and prints whether
// each one holds. The subcommand is the "self-consistency" view of the
// tool and returns a non-zero status when any rule fails, which happens
// only if the input geometry itself is inconsistent.
func RunChecks(in Input) (string, error) {
	g, err := in.ToGrating()
	if err != nil {
		return "", err
	}
	m := 1
	n := g.Slits
	if n < 1 {
		n = 1000
	}
	rules, err := disperse.CrossRules(g, m, n)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString(g.Summarize().Header())
	b.WriteString("\n")
	b.WriteString("cross rules of the plane grating model:\n")
	ok := true
	for _, r := range rules {
		mark := "PASS"
		if !r.Holds {
			mark = "FAIL"
			ok = false
		}
		b.WriteString(fmt.Sprintf("  [%s] %s: %s\n", mark, r.Name, r.Detail))
	}
	if !ok {
		return b.String(), fmt.Errorf("cross-rule check failed on the supplied geometry")
	}
	return b.String(), nil
}

// RunDensity prints the reciprocal relation between the groove spacing and
// the line density for the supplied geometry.
func RunDensity(in Input) (string, error) {
	g, err := in.ToGrating()
	if err != nil {
		return "", err
	}
	d, err := grating.DensityOf(g.GrooveSpacingNm)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("d = %s nm\n1/d = %s\n",
		formatNm(d.SpacingNm), g.LineDensityName()), nil
}

// formatNm renders a nanometre value for prose output.
func formatNm(v float64) string {
	return fmt.Sprintf("%.6g", v)
}
