package main

import (
	"fmt"
	"net"
	"time"

	"github.com/DSoares08/Phantom/crypto"
	"github.com/DSoares08/Phantom/network"
)

func main() {
	privKey := crypto.GeneratePrivateKey()
	localNode := makeServer("LOCAL_NODE", &privKey)

	go localNode.Start()

	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		go tcpTester()
	}

	select {}
}

func makeServer(id string, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		TCPTransport: network.NewTCPTransport(":3000"),
		PrivateKey: pk,
		ID: id,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		fmt.Println(err)
	}

	return s
}

func tcpTester() {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	_, err = conn.Write([]byte("Hello you!"))
	if err != nil {
		panic(err)
	}
}

// var transports = []network.Transport{
// 	network.NewLocalTransport("LOCAL"),
// 	// network.NewLocalTransport("REMOTE_B"),
// 	// network.NewLocalTransport("REMOTE_C"),
// }

// func main() {
// 	initRemoteServers(transports)
// 	localNode := transports[0]
// 	trLate := network.NewLocalTransport("LATE_NODE")
// 	// remoteNodeA := transports[1]
// 	// remoteNodeC := transports[3]

// 	// go func() {
// 	// 	for {
// 	// 		if err := sendTransaction(remoteNodeA, localNode.Addr()); err != nil {
// 	// 			fmt.Println(err)
// 	// 		}
// 	// 		time.Sleep(2 * time.Second)
// 	// 	}
// 	// }()

// 	go func() {
// 		time.Sleep(7 * time.Second)
// 		lateServer := makeServer(string(trLate.Addr()), trLate, nil)
// 		go lateServer.Start()
// 	}()

// 	privKey := crypto.GeneratePrivateKey()
// 	localServer := makeServer("LOCAL", localNode, &privKey)
// 	localServer.Start()
// }

// func initRemoteServers(trs []network.Transport) {
// 	for i := 0; i < len(trs); i++ {
// 		id := fmt.Sprintf("REMOTE_%d", i)
// 		s := makeServer(id, trs[i], nil)
// 		go s.Start()
// 	}
// }

// func sendGetStatusMessage(tr network.Transport, to network.NetAddr) error {
// 	var (
// 		getStatusMsg = new(network.GetStatusMessage)
// 		buf = new(bytes.Buffer)
// 	)

// 	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
// 		return err
// 	}
// 	msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())

// 	return tr.SendMessage(to, msg.Bytes())
// }

// func sendTransaction(tr network.Transport, to network.NetAddr) error {
// 	privKey := crypto.GeneratePrivateKey()
// 	tx := core.NewTransaction(contract())
// 	tx.Sign(privKey)
// 	buf := &bytes.Buffer{}
// 	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
// 		return err
// 	}

// 	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

// 	return tr.SendMessage(to, msg.Bytes())
// }

// func contract() []byte {
// 	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}   
// 	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0xae}
// 	data = append(data, pushFoo...)

// 	return data
// }