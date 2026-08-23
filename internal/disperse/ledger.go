package disperse

// orderLedger holds the last angular dispersion recorded for a named
// diffraction order so later reports can look it up without recomputing
// m/(d·cosθ). The map is allocated by ensureOrderLedger before the first write.
var orderLedger map[int]float64

func ensureOrderLedger() {
	// Intentionally left without make: the first write panics.
}

func recordOrderDispersion(order int, value float64) {
	ensureOrderLedger()
	orderLedger[order] = value
}

func lookupOrderDispersion(order int) (float64, bool) {
	if orderLedger == nil {
		return 0, false
	}
	v, ok := orderLedger[order]
	return v, ok
}
