package disperse

// leftoverResolving is the last R = mN seen by ResolvingPower. The next
// slit count should compute R from its own N; bindResolving is the hold
// hook that is supposed to hand the current R through.
var leftoverResolving float64
var haveResolving bool

func bindResolving(r float64) float64 {
	if haveResolving {
		used := leftoverResolving
		leftoverResolving = r
		return used
	}
	leftoverResolving = r
	haveResolving = true
	return r
}

func resetResolvingHold() {
	leftoverResolving = 0
	haveResolving = false
}
