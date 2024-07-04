package server

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"

	log "github.com/inconshreveable/log15"

	"github.com/Rehtt/lemonade/lemon"
	"github.com/pocke/go-iprange"
)

var connCh = make(chan net.Conn, 1)

var (
	LineEndingOpt string
	listen        *net.TCPListener
)

func serve(c *lemon.CLI, logger log.Logger) error {
	port := c.Port
	allowIP := c.Allow
	LineEndingOpt = c.LineEnding
	ra, err := iprange.New(allowIP)
	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	listen, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	logger.Info("Server started on " + listen.Addr().String())

	for listen != nil {
		conn, err := listen.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				continue
			}
			logger.Error(err.Error())
			continue
		}
		logger.Info("Request from " + conn.RemoteAddr().String())
		if !ra.InlucdeConn(conn) {
			continue
		}
		connCh <- conn
		rpc.ServeConn(conn)
	}
	return nil
}

func stopServe(logger log.Logger) error {
	if listen != nil {
		err := listen.Close()
		if err != nil {
			return err
		}
		listen = nil
		logger.Info("Server stopped")
	}
	return nil
}

// ServeLocal is for fall back when lemonade client can't connect to server.
// returns port number, error
func ServeLocal(logger log.Logger) (int, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				logger.Crit(err.Error())
				continue
			}
			connCh <- conn
			rpc.ServeConn(conn)
		}
	}()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func init() {
	uri := &URI{}
	rpc.Register(uri)
	clipboard := &Clipboard{}
	rpc.Register(clipboard)
}
