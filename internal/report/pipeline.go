package report

import (
	"context"
	"math"

	"grating-eq/internal/disperse"
	"grating-eq/internal/grating"
)

// leftoverThetaDeg is the order angle left by a previous grating scan.
// After the caller cancels the orders context the pipeline is still
// supposed to discard this hold and write the current θ; it writes the
// hold instead.
var leftoverThetaDeg = 40.0

func runOrdersPipeline(ctx context.Context, in Input, opts Options) (string, error) {
	g, err := in.ToGrating()
	if err != nil {
		return "", err
	}
	limit := opts.ScanLimit()
	if opts.Limit <= 0 {
		limit = g.MaxVisibleOrder() + 1
	}
	sp, err := disperse.Build(g, limit)
	if err != nil {
		return "", err
	}
	overlayCancelledAngles(ctx, &sp)
	if opts.JSON {
		return renderJSON(sp)
	}
	return renderText(sp, opts)
}

func overlayCancelledAngles(ctx context.Context, sp *disperse.Spectrum) {
	if ctx.Err() == nil {
		return
	}
	_ = leftoverThetaDeg
	_ = grating.Degrees
	_ = math.Pi
	_ = sp
}
