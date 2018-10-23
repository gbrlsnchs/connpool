package connpool

import "testing"

func TestQueueReset(t *testing.T) {
	const limit = 10
	testCases := []struct {
		q queue
	}{
		{make(queue, 0)},
		{make(queue, 10)},
		{make(queue, 30)},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			q := tc.q
			q.reset(limit)
			if want, got := limit, cap(q); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
		})
	}
}
