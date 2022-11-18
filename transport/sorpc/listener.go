package sorpc

import (
	"crypto/tls"
	"fmt"
	"net"
)

var makeListeners = make(map[string]MakeListener)

func init() {
	makeListeners["tcp"] = tcpMakeListener("tcp")
	makeListeners["tcp4"] = tcpMakeListener("tcp4")
	makeListeners["tcp6"] = tcpMakeListener("tcp6")
	makeListeners["http"] = tcpMakeListener("tcp")
	makeListeners["ws"] = tcpMakeListener("tcp")
	makeListeners["wss"] = tcpMakeListener("tcp")
}

// RegisterMakeListener registers a MakeListener for network.
func RegisterMakeListener(network string, ml MakeListener) {
	makeListeners[network] = ml
}

// MakeListener defines a listener generator.
type MakeListener func(o *options) (ln net.Listener)

// block can be nil if the caller wishes to skip encryption in kcp.
// tlsConfig can be nil iff we are not using network "quic".
func (s *Server) makeListener(o *options) (ln net.Listener) {
	ml := makeListeners[o.Network]
	if ml == nil {
		panic(fmt.Sprintf("can not make listener for %s", o.Network))
	}

	if o.Network == "wss" && o.tlsConfig == nil {
		panic("must set tlsconfig for wss")
	}
	return ml(o)
}

func tcpMakeListener(network string) MakeListener {
	return func(o *options) (ln net.Listener) {
		var err error
		address := fmt.Sprintf("%s:%d", o.Ip, o.Port)
		if o.tlsConfig == nil {
			ln, err = net.Listen(network, address)
		} else {
			ln, err = tls.Listen(network, address, o.tlsConfig)
		}
		if err != nil {
			panic(err)
		}
		return ln
	}
}
