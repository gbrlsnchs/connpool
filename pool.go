package connpool

import (
	"context"
	"net"
)

const defaultMaxIdleConns = 2

// Pool is a connection pool.
type Pool struct {
	network string
	address string
	c       chan net.Conn // conn pool
	q       chan struct{} // open conns queue
}

// New creates a new connection pool.
func New(network, address string) *Pool {
	p := &Pool{network: network, address: address}
	p.SetMaxIdleConns(defaultMaxIdleConns)
	return p
}

// Cap returns the capacity of the pool,
// which is the number of maximum idle connections.
func (p *Pool) Cap() int {
	return cap(p.c)
}

// Dial tries to stablish a new connection
// while respecting the limit for open connections.
func (p *Pool) Dial() (net.Conn, error) {
	return p.DialContext(context.Background())
}

// DialContext tries to stablish a new connection before a context is canceled
// while respecting the limit for open connections.
func (p *Pool) DialContext(ctx context.Context) (net.Conn, error) {
	var err error
	if err = p.wait(ctx); err != nil {
		return nil, err
	}
	var d net.Dialer
	c, err := d.DialContext(ctx, p.network, p.address)
	if err != nil {
		return nil, err
	}
	return &conn{c, p.c, p.q}, nil
}

// Get retrieves a new connection if any is available,
// otherwise it spawns a new connection.
func (p *Pool) Get() (net.Conn, error) {
	return p.GetContext(context.Background())
}

// GetContext retrieves a new connection or spawns a new one
// if the context is still not done.
func (p *Pool) GetContext(ctx context.Context) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		select {
		case conn := <-p.c:
			return conn, nil
		default:
			return p.DialContext(ctx)
		}
	}
}

// Len returns the number of idle connections.
func (p *Pool) Len() int {
	return len(p.c)
}

// SetMaxIdleConns limits the amount of idle connections in the pool.
func (p *Pool) SetMaxIdleConns(n int) {
	if maxOpenConns := cap(p.q); maxOpenConns > 0 && n > maxOpenConns {
		n = maxOpenConns
	}
	if n != cap(p.c) {
		if p.c != nil {
			close(p.c)
		}
		if n == 0 {
			// Don't reuse any connections.
			p.c = nil
			return
		}
		p.c = make(chan net.Conn, n)
	}
}

// SetMaxOpenConns limits the amount of open connections.
func (p *Pool) SetMaxOpenConns(n int) {
	if n > 0 {
		p.SetMaxIdleConns(n)
	}
	if n != cap(p.q) {
		if p.q != nil {
			close(p.q)
		}
		if n == 0 {
			// A nil queue channel means open connections are limitless.
			p.q = nil
			return
		}
		p.q = make(chan struct{}, n)
	}
}

func (p *Pool) wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if cap(p.q) > 0 {
			p.q <- struct{}{}
		}
		return nil
	}
}
