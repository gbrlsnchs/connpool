package connpool

type queue chan struct{}

func (q *queue) reset(limit int) {
	qq := *q
	// Reset channel only if size has changed.
	if limit != cap(qq) {
		if qq != nil {
			close(qq)
		}
		if limit > 0 {
			*q = make(queue, limit)
			return
		}
		// A nil queue channel means open connections are limitless.
		*q = nil
	}
}
