// Copyright 2014 Wandoujia Inc. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package redispool

import (
	"bufio"
	"net"
)

//not thread-safe
type Conn struct {
	DB int

	addr string
	net.Conn
	closed bool
	r      *bufio.Reader
}

func (c *Conn) Close() {
	c.Conn.Close()
	c.closed = true
}

func (c *Conn) IsClosed() bool {
	return c.closed
}

type PooledConn struct {
	*Conn
	pool *ConnectionPool
}

func (pc *PooledConn) Recycle() {
	if pc.IsClosed() {
		pc.pool.Put(nil)
	} else {
		pc.pool.Put(pc)
	}
}

//requre read to use bufio
func (pc *PooledConn) Read(p []byte) (int, error) {
	panic("not allowed")
}

func (pc *PooledConn) Write(p []byte) (int, error) {
	return pc.Conn.Write(p)
}

func (pc *PooledConn) BufioReader() *bufio.Reader {
	return pc.r
}

func NewConnection(key string) (*Conn, error) {
	addr := key

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Conn{
		DB:   0,
		addr: addr,
		Conn: conn,
		r:    bufio.NewReaderSize(conn, 204800),
	}, nil
}

func ConnectionCreator(key string) CreateConnectionFunc {
	return func(pool *ConnectionPool) (PoolConnection, error) {
		c, err := NewConnection(key)
		if err != nil {
			return nil, err
		}
		return &PooledConn{c, pool}, nil
	}
}
