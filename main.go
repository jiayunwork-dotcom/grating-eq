// Command grating-eq is a plane diffraction grating calculator.
//
// It reads a grating description (groove spacing or line density,
// wavelength, incident angle and optional finite slit count) from a JSON
// file and reports, for each diffraction order m, the angle theta_m, the
// existence of the order, the angular dispersion, the free spectral range
// and, when a finite slit count is given, the resolving power and the
// minimum resolvable wavelength interval.
//
// This file only wires the subcommands; the physics lives in the internal
// packages grating, disperse and report.
package main

import (
	"fmt"
	"os"
	"strings"

	"grating-eq/internal/report"
)

const usage = `grating-eq: plane diffraction grating calculator.

The grating equation m*lambda = d*(sin(theta_m) - sin(theta_i)) is solved
for every order in a scan window. Angles are measured from the grating
normal and the transmission convention is used throughout. A JSON input
file describes the groove geometry with either "grooves_per_mm" or "d_nm"
(spacing in nanometres), the wavelength in "wavelength_nm", the incident
angle in degrees in "incident_angle_deg" and an optional finite slit count
in "slits".

usage:
  grating-eq orders <grating.json> [--limit n] [--json] [--existing]
  grating-eq checks <grating.json>
  grating-eq density <grating.json>
  grating-eq help

Examples:
  grating-eq orders example/600lpmm.json
  grating-eq orders example/600lpmm.json --limit 5
  grating-eq checks example/600lpmm.json

Illegal inputs (d <= 0, lambda <= 0, |sin(incident)| > 1, N < 1 when
given, unknown JSON fields, missing files) are reported on stderr and the
process exits non-zero.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "orders":
		err = runOrders(os.Args[2:])
	case "checks":
		err = runChecks(os.Args[2:])
	case "density":
		err = runDensity(os.Args[2:])
	case "help", "-h", "--help":
		fmt.Print(usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "grating-eq: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
	if err != nil {
		if report.IsHelp(err) {
			fmt.Print(usage)
			return
		}
		fmt.Fprintf(os.Stderr, "grating-eq: %v\n", err)
		os.Exit(1)
	}
}

// reorderArgs moves every positional argument to the end, so flags may
// appear before or after the JSON file name. flagsWithValue lists the
// flags that consume the following token as their value.
func reorderArgs(args []string, flagsWithValue ...string) []string {
	takesValue := make(map[string]bool, len(flagsWithValue))
	for _, f := range flagsWithValue {
		takesValue[f] = true
	}
	var flags, positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-") && a != "-" {
			flags = append(flags, a)
			name := strings.TrimLeft(a, "-")
			if eq := strings.IndexByte(name, '='); eq >= 0 {
				name = name[:eq]
			}
			if takesValue[name] && !strings.Contains(a, "=") && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i++
			}
			continue
		}
		positional = append(positional, a)
	}
	return append(flags, positional...)
}

func runOrders(args []string) error {
	opts, pos, err := report.ParseOptions(reorderArgs(args, "limit"))
	if err != nil {
		return err
	}
	if len(pos) != 1 {
		return fmt.Errorf("orders needs exactly one grating JSON file")
	}
	in, err := report.ParseFile(pos[0])
	if err != nil {
		return err
	}
	out, err := report.Run(in, opts)
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}

func runChecks(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("checks needs exactly one grating JSON file")
	}
	in, err := report.ParseFile(args[0])
	if err != nil {
		return err
	}
	out, err := report.RunChecks(in)
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}

func runDensity(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("density needs exactly one grating JSON file")
	}
	in, err := report.ParseFile(args[0])
	if err != nil {
		return err
	}
	out, err := report.RunDensity(in)
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}
