package socks

import (
	"io"
	"io/ioutil"
	"net"

	adapters "github.com/slimemice/clash/adapters/inbound"
	"github.com/slimemice/clash/component/socks5"
	C "github.com/slimemice/clash/constant"
	"github.com/slimemice/clash/log"
	authStore "github.com/slimemice/clash/proxy/auth"
	"github.com/slimemice/clash/tunnel"
)

type SockListener struct {
	net.Listener
	address string
	closed  bool
}

func NewSocksProxy(addr string) (*SockListener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	sl := &SockListener{l, addr, false}
	go func() {
		log.Infoln("SOCKS proxy listening at: %s", addr)
		for {
			c, err := l.Accept()
			if err != nil {
				if sl.closed {
					break
				}
				continue
			}
			go handleSocks(c)
		}
	}()

	return sl, nil
}

func (l *SockListener) Close() {
	l.closed = true
	l.Listener.Close()
}

func (l *SockListener) Address() string {
	return l.address
}

func handleSocks(conn net.Conn) {
	target, command, err := socks5.ServerHandshake(conn, authStore.Authenticator())
	if err != nil {
		conn.Close()
		return
	}
	conn.(*net.TCPConn).SetKeepAlive(true)
	if command == socks5.CmdUDPAssociate {
		defer conn.Close()
		io.Copy(ioutil.Discard, conn)
		return
	}
	tunnel.Add(adapters.NewSocket(target, conn, C.SOCKS, C.TCP))
}
