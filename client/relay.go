package client

import (
	"log"
	"path"
	"strings"
	"time"

	"github.com/nchuxyz/firefly/pkg/chain"
	"github.com/nchuxyz/firefly/pkg/gosocks"
)

type RelayHandler struct {
	Basic                     *gosocks.BasicSocksHandler
	NextHop                   string
	EmbeddedTunnellingDomains map[string]bool
	CustomTunnellingDomains   []string
	TunnellingAll             bool
}

func (r *RelayHandler) lookup(dst string, conn *gosocks.SocksConn) chain.SocksChain {
	chain := &chain.SocksSocksChain{
		SocksDialer: &gosocks.SocksDialer{
			Timeout: conn.Timeout,
			Auth:    &gosocks.AnonymousClientAuthenticator{},
		},
		SocksAddr: r.NextHop,
	}

	if r.TunnellingAll {
		return chain
	}

	labels := strings.Split(dst, ".")
	for i := 0; i < len(labels); i++ {
		_, ok := r.EmbeddedTunnellingDomains[strings.Join(labels[i:], ".")]
		if ok {
			return chain
		}
	}

	for _, v := range r.CustomTunnellingDomains {
		matched, _ := path.Match(v, dst)
		if matched {
			return chain
		}
	}

	return nil
}

func (r *RelayHandler) handleUDPAssociate(req *gosocks.SocksRequest, conn *gosocks.SocksConn) {
	clientBind, clientAssociate, udpReq, clientAddr, err := r.Basic.UDPAssociateFirstPacket(req, conn)
	if err != nil {
		conn.Close()
		return
	}
	chain := r.lookup(udpReq.DstHost, conn)
	if chain != nil {
		chain.UDPAssociate(req, conn, clientBind, clientAssociate, udpReq, clientAddr)
	} else {
		r.Basic.UDPAssociateForwarding(conn, clientBind, clientAssociate, udpReq, clientAddr)
	}
}

func (r *RelayHandler) ServeSocks(conn *gosocks.SocksConn) {
	conn.SetReadDeadline(time.Now().Add(conn.Timeout))
	req, err := gosocks.ReadSocksRequest(conn)
	if err != nil {
		log.Printf("error in ReadSocksRequest: %s", err)
		return
	}

	switch req.Cmd {
	case gosocks.SocksCmdConnect:
		chain := r.lookup(req.DstHost, conn)
		if chain != nil {
			chain.TCP(req, conn)
		} else {
			r.Basic.HandleCmdConnect(req, conn)
		}
		return
	case gosocks.SocksCmdUDPAssociate:
		r.handleUDPAssociate(req, conn)
		return
	case gosocks.SocksCmdBind:
		conn.Close()
		return
	default:
		return
	}
}

func (r *RelayHandler) Quit() {}
