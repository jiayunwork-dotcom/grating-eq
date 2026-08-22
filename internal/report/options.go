package report

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Options control the orders subcommand output.
type Options struct {
	// Limit is the scan half-width: orders m in [-Limit, Limit] are
	// evaluated. Zero means the default of 64.
	Limit int

	// JSON reports the table as a JSON document instead of aligned text.
	JSON bool

	// ExistingOnly prints only the orders that exist.
	ExistingOnly bool
}

// DefaultOptions returns the option set used when no flags are given.
func DefaultOptions() Options {
	return Options{Limit: 0, JSON: false, ExistingOnly: false}
}

// ScanLimit resolves the effective scan half-width.
func (o Options) ScanLimit() int {
	if o.Limit <= 0 {
		return 64
	}
	return o.Limit
}

// ParseOptions parses command-line flags that the main package forwards.
// Flags may appear before or after the positional file argument.
func ParseOptions(args []string) (Options, []string, error) {
	opts := DefaultOptions()
	var rest []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--limit" || a == "-l":
			if i+1 >= len(args) {
				return opts, nil, fmt.Errorf("flag %s needs a value", a)
			}
			i++
			if _, err := fmt.Sscanf(args[i], "%d", &opts.Limit); err != nil {
				return opts, nil, fmt.Errorf("invalid --limit value %q", args[i])
			}
		case strings.HasPrefix(a, "--limit="):
			if _, err := fmt.Sscanf(a[len("--limit="):], "%d", &opts.Limit); err != nil {
				return opts, nil, fmt.Errorf("invalid --limit value %q", a)
			}
		case a == "--json":
			opts.JSON = true
		case a == "--existing":
			opts.ExistingOnly = true
		case a == "--help" || a == "-h":
			return opts, nil, errHelp
		case strings.HasPrefix(a, "-"):
			return opts, nil, fmt.Errorf("unknown flag %q", a)
		default:
			rest = append(rest, a)
		}
	}
	return opts, rest, nil
}

// errHelp is a sentinel so the main package can print the usage text.
var errHelp = fmt.Errorf("help requested")

// IsHelp reports whether the error came from a --help flag.
func IsHelp(err error) bool {
	return err == errHelp
}

// TableRow is the JSON form of one order line, mirroring what the text
// table prints.
type TableRow struct {
	Order        int     `json:"order"`
	Exists       bool    `json:"exists"`
	ThetaDeg     float64 `json:"theta_deg,omitempty"`
	SinTheta     float64 `json:"sin_theta"`
	Dispersion   float64 `json:"dispersion_deg_per_nm,omitempty"`
	FSRNm        float64 `json:"fsr_nm,omitempty"`
	ResolvingPow float64 `json:"resolving_power,omitempty"`
	DeltaLambda  float64 `json:"delta_lambda_nm,omitempty"`
}

// JSONOutput is the full JSON document emitted with --json.
type JSONOutput struct {
	GroovesPerMm    float64    `json:"grooves_per_mm"`
	SpacingNm       float64    `json:"spacing_nm"`
	WavelengthNm    float64    `json:"wavelength_nm"`
	IncidentDeg     float64    `json:"incident_angle_deg"`
	Slits           int        `json:"slits"`
	MaxVisibleOrder int        `json:"max_visible_order"`
	Orders          []TableRow `json:"orders"`
}

// EncodeJSON renders the output document.
func EncodeJSON(out JSONOutput) ([]byte, error) {
	return json.MarshalIndent(out, "", "  ")
}
