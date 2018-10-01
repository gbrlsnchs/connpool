package cpool

import "net"

// Conn is a Redis connection.
type conn struct {
	net.Conn
	p chan<- net.Conn
	q <-chan struct{}
}

// Close either returns the connection to the pool or,
// if the pool is full, it simply closes the connection.
//
// If it can't return the connection to the pool,
// it tries to dequeue the connection in order to respect
// the max open connections limit before truly closing it.
func (c *conn) Close() error {
	// Try to send the connection back to the pool,
	// otherwise simply close it.
	select {
	case c.p <- c:
		return nil
	default:
	}

	select {
	// If a limit of open conns is set (maxOpenConns > 0),
	// dequeue the connection.
	case <-c.q:
	default:
	}
	return c.Conn.Close()
}
