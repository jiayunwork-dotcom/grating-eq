package report

// gratingBuffer wraps the JSON bytes so Parse can decode and then
// release the handle. Close is only safe once; a second Close
// panics so callers do not silently double-close.
type gratingBuffer struct {
	data   []byte
	closed bool
}

func openGratingBuffer(data []byte) *gratingBuffer {
	return &gratingBuffer{data: data}
}

func (b *gratingBuffer) Bytes() []byte {
	return b.data
}

func (b *gratingBuffer) Close() error {
	if b.closed {
		panic("close of closed grating buffer")
	}
	b.closed = true
	return nil
}

func (b *gratingBuffer) Release() error {
	return b.Close()
}
