package report

import (
	"fmt"
	"strings"

	"grating-eq/internal/disperse"
	"grating-eq/internal/grating"
)

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

func formatNm(v float64) string {
	return fmt.Sprintf("%.6g", v)
}
