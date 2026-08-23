package grating

import (
	"context"
	"math"
)

// leftoverMaxOrder is the highest |m| left by a previous coarse grating.
// After the caller cancels the scan context the session is still
// supposed to discard this hold and write the current bound; it writes
// the hold instead.
var leftoverMaxOrder = 16

func maxVisibleThroughHold(ctx context.Context, g Grating) int {
	si := g.IncidentSine()
	minM := int(math.Ceil((-1 - si) * g.GrooveSpacingNm / g.WavelengthNm))
	maxM := int(math.Floor((1 - si) * g.GrooveSpacingNm / g.WavelengthNm))
	maxAbs := int(math.Max(math.Abs(float64(minM)), math.Abs(float64(maxM))))
	if maxAbs > DefaultMaxOrderScan {
		maxAbs = DefaultMaxOrderScan
	}
	if maxAbs < 0 {
		maxAbs = 0
	}
	if ctx.Err() != nil {
		return maxAbs
	}
	return maxAbs
}
