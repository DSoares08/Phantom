package network

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	conn net.Conn
}

type TCPTransport struct {
	listenAddr string
	listener net.Listener
}

func NewTCPTransport(addr string, peerCh chan *TCPPeer) *TCPTransport {
	return &TCPTransport{peerCh: peerCh, listenAddr: addr}
}

func (t *TCPTransport) readLoop(peer *TCPPeer) {
	buf := make([]byte, 2048)
	for {
		n, err := peer.conn.Read(buf)
		if err != nil {
			fmt.Printf("read error: %s", err)
			continue
		}

		msg := buf[:n]
		fmt.Printf(string(msg))
		// handleMessage => server
	}
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("accept error from %+v\n", conn)
			continue
		}

		peer := &TCPPeer{conn: conn}

		t.peerCh <- peer

		fmt.Printf("new incoming TCP connection => %+v\n", conn)

		go t.readLoop(peer)
	}
}

func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	t.listener = ln

	go t.acceptLoop()

	fmt.Println("TCP transport listening to port: ", t.listenAddr)

	return nil
}
