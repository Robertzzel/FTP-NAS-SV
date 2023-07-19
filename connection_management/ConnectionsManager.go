package connection_management

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
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
	publicIp      string
}

func NewConnectionManager(certificateFilePath, keyFilePath string) (ConnectionManager, error) {
	cert, err := tls.LoadX509KeyPair(certificateFilePath, keyFilePath)
	if err != nil {
		return ConnectionManager{}, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12, MaxVersion: tls.VersionTLS12}

	cm := ConnectionManager{tlsConfigs: config}

	ipString, err := getPublicIP()
	if err != nil {
		return ConnectionManager{}, err
	}
	cm.publicIp = ipString

	return cm, nil
}

func (manager *ConnectionManager) ListenForClientsPI() (net.Listener, error) {
	l, err := net.Listen(Network, fmt.Sprintf("%s:%d", Host, PIPort))
	if err != nil {
		return nil, err
	}
	return tls.NewListener(l, manager.tlsConfigs), nil
}

func (manager *ConnectionManager) ListenToAvailablePort() (net.Listener, error) {
	ip := fmt.Sprintf("%s:0", Host)
	l, err := net.Listen(Network, ip)
	if err != nil {
		return nil, err
	}
	return tls.NewListener(l, manager.tlsConfigs), nil
}

func (manager *ConnectionManager) GetPublicIP() string {
	return manager.publicIp
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}
