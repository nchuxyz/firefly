package chain

import (
	"net"

	"github.com/nchuxyz/firefly/pkg/gosocks"
)

type SocksChain interface {
	TCP(*gosocks.SocksRequest, *gosocks.SocksConn)
	UDPAssociate(*gosocks.SocksRequest, *gosocks.SocksConn, *net.UDPConn, *net.UDPAddr, *gosocks.UDPRequest, *net.UDPAddr)
}
