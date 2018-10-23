package connpool_test

import (
	"net"
	"sync"
	"testing"

	. "github.com/gbrlsnchs/connpool"
)

func TestPool(t *testing.T) {
	const address = ":6060"
	l, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	testCases := []struct {
		maxIdleConns int
		maxOpenConns int
		expectedLen  int
		expectedCap  int
		count        int
	}{
		{0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1},
		{1, 1, 1, 1, 1},
		{10, 10, 1, 10, 1},
		{10, 1, 1, 1, 1},
		{100, 100, 1, 100, 1},
		{100, 1, 1, 1, 1},
		{2048, 2048, 1, 2048, 1},
		{2048, 1, 1, 1, 1},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(tc.count)
			for i := 0; i < tc.count; i++ {
				go func() {
					conn, err := l.Accept()
					if err != nil {
						t.Error(err)
					}
					conn.Close()
					wg.Done()
				}()
			}

			p := New("tcp", address)
			p.SetMaxIdleConns(tc.maxIdleConns)
			p.SetMaxOpenConns(tc.maxOpenConns)
			for i := 0; i < tc.count; i++ {
				conn, err := p.Get()
				if want, got := (error)(nil), err; want != got {
					t.Fatalf("want %v, got %v", want, got)
				}
				if want, got := (error)(nil), conn.Close(); want != got {
					t.Errorf("want %v, got %v", want, got)
				}
			}

			wg.Wait()
			if want, got := tc.expectedCap, p.Cap(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
			if want, got := tc.expectedLen, p.Len(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
		})
	}
}
