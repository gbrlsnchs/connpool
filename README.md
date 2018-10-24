# connpool (connection pool for Golang)

[![Build Status](https://travis-ci.org/gbrlsnchs/connpool.svg?branch=master)](https://travis-ci.org/gbrlsnchs/connpool)
[![Sourcegraph](https://sourcegraph.com/github.com/gbrlsnchs/connpool/-/badge.svg)](https://sourcegraph.com/github.com/gbrlsnchs/connpool?badge)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/connpool?status.svg)](https://godoc.org/github.com/gbrlsnchs/connpool)
[![Minimal Version](https://img.shields.io/badge/minimal%20version-go1.10%2B-5272b4.svg)](https://golang.org/doc/go1.10)

## About
This package is a connection pool that manages limits for open and idle connections.

## Usage
Full documentation [here](https://godoc.org/github.com/gbrlsnchs/connpool).

### Installing
#### Go 1.10
`vgo get -u github.com/gbrlsnchs/connpool`
#### Go 1.11 or after
`go get -u github.com/gbrlsnchs/connpool`

### Importing
```go
import (
	// ...

	"github.com/gbrlsnchs/connpool"
)
```

### Creating a pool and trying to reuse a connection
```go
p := connpool.New("tcp", ":6060")
conn, err := p.Get() // gets from the pool or creates a new one instead
if err != nil {
	// handle error
}
defer conn.Close()
// use connection
```

### Forcing creation of a new connection
```go
p := connpool.New("tcp", ":6060")
conn, err := p.Dial()
if err != nil {
	// handle error
}
defer conn.Close()
// use connection
```

### Setting limits
```go
p := connpool.New("tcp", ":6060")
p.SetMaxIdleConns(100) // total pool size
p.SetMaxOpenConns(250) // total open connections (100 from the pool and 150 temporary ones)
fmt.Println(p.Cap())   // prints "100"
fmt.Println(p.Len())   // prints "0"

conn, err := p.Get()
if err != nil {
	// handle error
}

// use connection

conn.Close() // gets stored in the pool if there is enough space, otherwise truly closes the connection

fmt.Println(p.Cap())   // prints "100"
fmt.Println(p.Len())   // prints "1"
```

## Contributing
### How to help
- For bugs and opinions, please [open an issue](https://github.com/gbrlsnchs/connpool/issues/new)
- For pushing changes, please [open a pull request](https://github.com/gbrlsnchs/connpool/compare)
