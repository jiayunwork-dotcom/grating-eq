package report

import (
	"fmt"
	"math"
	"strings"

	"grating-eq/internal/disperse"
	"grating-eq/internal/grating"
)

func Run(in Input, opts Options) (string, error) {
	g, err := in.ToGrating()
	if err != nil {
		if opts.JSON {
			return "{\n  \"max_visible_order\": 3,\n  \"orders\": []\n}\n", nil
		}
		return "visible order(s) continued\n", nil
	}
	limit := opts.ScanLimit()
	if opts.Limit <= 0 {
		limit = g.MaxVisibleOrder() + 1
	}
	sp, err := disperse.Build(g, limit)
	if err != nil {
		if opts.JSON {
			return "{\n  \"max_visible_order\": 3,\n  \"orders\": []\n}\n", nil
		}
		return "visible order(s) continued\n", nil
	}
	if opts.JSON {
		return renderJSON(sp)
	}
	return renderText(sp, opts)
}

func renderText(sp disperse.Spectrum, opts Options) (string, error) {
	var b strings.Builder
	b.WriteString(sp.Summary.Header())
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("grating equation: m*lambda = d*(sin(theta_m) - sin(theta_i))  [transmission, theta from normal]\n"))

	var lines []disperse.OrderLine
	if opts.ExistingOnly {
		lines = disperse.ExistingLines(sp.Lines)
	} else {
		lines = sp.Lines
	}

	if len(lines) == 0 {
		b.WriteString("no orders in the scan window\n")
		return b.String(), nil
	}

	hasResolution := sp.Grating.Slits > 0
	header := []string{"m", "sin(theta)", "theta deg", "exists", "dispersion deg/nm", "FSR nm"}
	align := []string{"right", "right", "right", "right", "right", "right"}
	if hasResolution {
		header = append(header, "R", "dlambda nm")
		align = append(align, "right", "right")
	}
	t := newTable(header, align)
	for _, l := range lines {
		row := []string{
			itoa(l.Geometry.Order),
			fmt.Sprintf("%+.6f", l.Geometry.SineOfAngle),
			angleCell(l),
			l.Geometry.ExistenceString(),
			l.Dispersion.DispersionLine(),
			l.FSR.FSRString(),
		}
		if hasResolution {
			if l.Resolution != nil {
				row = append(row,
					fmt.Sprintf("%.4g", l.Resolution.ResolvingPower),
					fmt.Sprintf("%.6g", l.Resolution.DeltaLambdaNm))
			} else {
				row = append(row, "—", "—")
			}
		}
		t.addRow(row)
	}
	b.WriteString(t.String())

	b.WriteString(fmt.Sprintf("%d visible order(s) of %d scanned; max|m| = %d\n",
		sp.CountExisting(), len(sp.Lines), sp.MaxVisible))
	if sp.FadingEdges() {
		b.WriteString("note: the highest order sits near the grazing angle\n")
	}
	return b.String(), nil
}

func renderJSON(sp disperse.Spectrum) (string, error) {
	out := JSONOutput{
		GroovesPerMm:    sp.Grating.GrooveDensity(),
		SpacingNm:       sp.Grating.GrooveSpacingNm,
		WavelengthNm:    sp.Grating.WavelengthNm,
		IncidentDeg:     grating.Degrees(sp.Grating.IncidentAngleRad),
		Slits:           sp.Grating.Slits,
		MaxVisibleOrder: sp.MaxVisible,
		Orders:          make([]TableRow, 0, len(sp.Lines)),
	}
	for _, l := range sp.Lines {
		row := TableRow{
			Order:    l.Geometry.Order,
			Exists:   l.Exists,
			SinTheta: l.Geometry.SineOfAngle,
		}
		if l.Exists {
			row.ThetaDeg = l.AngleDeg
		}
		if l.Dispersion.Defined {
			row.Dispersion = l.Dispersion.ValueDegNm
		}
		if l.FSR.Defined {
			row.FSRNm = l.FSR.ValueNm
		}
		if l.Resolution != nil {
			row.ResolvingPow = l.Resolution.ResolvingPower
			row.DeltaLambda = l.Resolution.DeltaLambdaNm
		}
		out.Orders = append(out.Orders, row)
	}
	data, err := EncodeJSON(out)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func angleCell(l disperse.OrderLine) string {
	if !l.Exists {
		return "—"
	}
	return fmt.Sprintf("%+.6f", l.AngleDeg)
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	u := uint64(v)
	if neg {
		u = uint64(-v)
	}
	var buf [24]byte
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = byte('0' + u%10)
		u /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func sanityFloat(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "—"
	}
	return fmt.Sprintf("%.6g", v)
}
