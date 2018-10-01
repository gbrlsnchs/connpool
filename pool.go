package cpool

import (
	"context"
	"net"
)

const defaultMaxIdleConns = 2

// Pool is a connection pool.
type Pool struct {
	network      string
	address      string
	maxOpenConns int
	maxIdleConns int
	c            chan net.Conn // conn pool
	q            chan struct{} // open conns queue
}

// New creates a new connection pool.
func New(network, address string) *Pool {
	p := &Pool{network: network, address: address}
	p.SetMaxIdleConns(defaultMaxIdleConns)
	return p
}

// Dial tries to stablish a new connection.
func (p *Pool) Dial() (net.Conn, error) {
	return p.DialContext(context.Background())
}

// DialContext tries to stablish a new connection before a context is canceled.
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
	return &conn{Conn: c, p: p.c, q: p.q}, nil
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

// SetMaxIdleConns limits the amount of idle connections in the pool.
func (p *Pool) SetMaxIdleConns(maxConns int) {
	p.maxIdleConns = maxConns
	if p.maxOpenConns > 0 && p.maxIdleConns > p.maxOpenConns {
		p.maxIdleConns = p.maxOpenConns
	}
	p.resetPool()
}

// SetMaxOpenConns limits the amount of open connections.
func (p *Pool) SetMaxOpenConns(maxConns int) {
	p.maxOpenConns = maxConns
	if p.maxOpenConns > 0 && p.maxIdleConns > p.maxOpenConns {
		p.maxIdleConns = p.maxOpenConns
	}
	p.resetPool()
	p.resetQueue()
}

func (p *Pool) resetPool() {
	// Reset channel only if size has changed.
	if p.maxIdleConns != cap(p.c) {
		if p.c != nil {
			close(p.c)
		}
		if p.maxIdleConns > 0 {
			p.c = make(chan net.Conn, p.maxIdleConns)
			return
		}
		// Don't reuse any connections.
		p.c = nil
	}
}

func (p *Pool) resetQueue() {
	// Reset channel only if size has changed.
	if p.maxOpenConns != cap(p.q) {
		if p.q != nil {
			close(p.q)
		}
		if p.maxOpenConns > 0 {
			p.q = make(chan struct{}, p.maxOpenConns)
			return
		}
		// A nil queue channel means open connections are limitless.
		p.q = nil
	}
}

func (p *Pool) wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if p.maxOpenConns > 0 {
			p.q <- struct{}{}
		}
		return nil
	}
}
