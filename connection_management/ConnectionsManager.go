package connection_management

import (
	"crypto/tls"
	"fmt"
	"net"
)

const (
	Network = "tcp"
	Host    = "localhost"
	PIPort  = 4040
)

type ConnectionManager struct {
	PIConnection  *net.Listener  // Protocol interpreter connection
	DTConnections []net.Listener // Data transfer connections
	tlsConfigs    *tls.Config
}

func NewConnectionManager(certificateFilePath, keyFilePath string) (ConnectionManager, error) {
	cert, err := tls.LoadX509KeyPair(certificateFilePath, keyFilePath)
	if err != nil {
		return ConnectionManager{}, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	return ConnectionManager{tlsConfigs: config}, nil
}

func (manager *ConnectionManager) ListenForClients() (net.Listener, error) {
	return tls.Listen(Network, fmt.Sprintf("%s:%d", Host, PIPort), manager.tlsConfigs)
}

func (manager *ConnectionManager) ListenToAvailablePort() (net.Listener, error) {
	listener, err := net.Listen(Network, fmt.Sprintf("%s:0", Host))
	if err != nil {
		return nil, err
	}

	manager.DTConnections = append(manager.DTConnections, listener)

	return listener, nil
}
