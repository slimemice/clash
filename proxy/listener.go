package proxy

import (
	"fmt"
	"net"
	"strconv"

	"github.com/slimemice/clash/component/resolver"
	"github.com/slimemice/clash/config"
	"github.com/slimemice/clash/dns"

	"github.com/slimemice/clash/proxy/http"
	"github.com/slimemice/clash/proxy/redir"
	"github.com/slimemice/clash/proxy/socks"
	"github.com/slimemice/clash/proxy/tun"
)

var (
	allowLan    = false
	bindAddress = "*"

	socksListener    *socks.SockListener
	socksUDPListener *socks.SockUDPListener
	httpListener     *http.HttpListener
	redirListener    *redir.RedirListener
	tunAdapter       tun.TunAdapter
)

type listener interface {
	Close()
	Address() string
}

type Ports struct {
	Port      int `json:"port"`
	SocksPort int `json:"socks-port"`
	RedirPort int `json:"redir-port"`
}

func AllowLan() bool {
	return allowLan
}

func BindAddress() string {
	return bindAddress
}

func SetAllowLan(al bool) {
	allowLan = al
}

func Tun() config.Tun {
	if tunAdapter == nil {
		return config.Tun{}
	}
	return config.Tun{
		Enable:    true,
		DeviceURL: tunAdapter.DeviceURL(),
		DNSListen: tunAdapter.DNSListen(),
	}
}

func SetBindAddress(host string) {
	bindAddress = host
}

func ReCreateHTTP(port int) error {
	addr := genAddr(bindAddress, port, allowLan)

	if httpListener != nil {
		if httpListener.Address() == addr {
			return nil
		}
		httpListener.Close()
		httpListener = nil
	}

	if portIsZero(addr) {
		return nil
	}

	var err error
	httpListener, err = http.NewHttpProxy(addr)
	if err != nil {
		return err
	}

	return nil
}

func ReCreateSocks(port int) error {
	addr := genAddr(bindAddress, port, allowLan)

	shouldTCPIgnore := false
	shouldUDPIgnore := false

	if socksListener != nil {
		if socksListener.Address() != addr {
			socksListener.Close()
			socksListener = nil
		} else {
			shouldTCPIgnore = true
		}
	}

	if socksUDPListener != nil {
		if socksUDPListener.Address() != addr {
			socksUDPListener.Close()
			socksUDPListener = nil
		} else {
			shouldUDPIgnore = true
		}
	}

	if shouldTCPIgnore && shouldUDPIgnore {
		return nil
	}

	if portIsZero(addr) {
		return nil
	}

	tcpListener, err := socks.NewSocksProxy(addr)
	if err != nil {
		return err
	}

	udpListener, err := socks.NewSocksUDPProxy(addr)
	if err != nil {
		tcpListener.Close()
		return err
	}

	socksListener = tcpListener
	socksUDPListener = udpListener

	return nil
}

func ReCreateRedir(port int) error {
	addr := genAddr(bindAddress, port, allowLan)

	if redirListener != nil {
		if redirListener.Address() == addr {
			return nil
		}
		redirListener.Close()
		redirListener = nil
	}

	if portIsZero(addr) {
		return nil
	}

	var err error
	redirListener, err = redir.NewRedirProxy(addr)
	if err != nil {
		return err
	}

	return nil
}

func ReCreateTun(conf config.Tun) error {
	enable := conf.Enable
	url := conf.DeviceURL
	if tunAdapter != nil {
		if enable && (url == "" || url == tunAdapter.DeviceURL()) {
			// Though we don't need to recreate tun device, we should update tun DNSServer
			return tunAdapter.ReCreateDNSServer(resolver.DefaultResolver.(*dns.Resolver), conf.DNSListen)
		}
		tunAdapter.Close()
		tunAdapter = nil
	}
	if !enable {
		return nil
	}
	var err error
	tunAdapter, err = tun.NewTunProxy(url)
	if err != nil {
		return err
	}
	if resolver.DefaultResolver != nil {
		return tunAdapter.ReCreateDNSServer(resolver.DefaultResolver.(*dns.Resolver), conf.DNSListen)
	}
	return nil
}

// GetPorts return the ports of proxy servers
func GetPorts() *Ports {
	ports := &Ports{}

	if httpListener != nil {
		_, portStr, _ := net.SplitHostPort(httpListener.Address())
		port, _ := strconv.Atoi(portStr)
		ports.Port = port
	}

	if socksListener != nil {
		_, portStr, _ := net.SplitHostPort(socksListener.Address())
		port, _ := strconv.Atoi(portStr)
		ports.SocksPort = port
	}

	if redirListener != nil {
		_, portStr, _ := net.SplitHostPort(redirListener.Address())
		port, _ := strconv.Atoi(portStr)
		ports.RedirPort = port
	}

	return ports
}

func portIsZero(addr string) bool {
	_, port, err := net.SplitHostPort(addr)
	if port == "0" || port == "" || err != nil {
		return true
	}
	return false
}

func genAddr(host string, port int, allowLan bool) string {
	if allowLan {
		if host == "*" {
			return fmt.Sprintf(":%d", port)
		} else {
			return fmt.Sprintf("%s:%d", host, port)
		}
	}

	return fmt.Sprintf("127.0.0.1:%d", port)
}
