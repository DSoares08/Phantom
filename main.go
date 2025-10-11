package main

import (
	"fmt"
	"bytes"
	"encoding/gob"

	"github.com/DSoares08/Phantom/core"
	"github.com/DSoares08/Phantom/crypto"
	"github.com/DSoares08/Phantom/network"
)

var transports = []network.Transport{
	network.NewLocalTransport("LOCAL"),
	network.NewLocalTransport("REMOTE_A"),
	network.NewLocalTransport("REMOTE_B"),
	network.NewLocalTransport("REMOTE_C"),
}

func main() {
	// trLocal := network.NewLocalTransport("LOCAL")
	// trRemoteA := network.NewLocalTransport("REMOTE_A")
	// trRemoteB := network.NewLocalTransport("REMOTE_B")
	// trRemoteC := network.NewLocalTransport("REMOTE_C")

	initRemoteServers(transports)

	// go func() {
	// 	for {
	// 		if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	// if err := sendGetStatusMessage(trRemoteA, "REMOTE_B"); err != nil {
	// 	log.Fatal(err)
	// }

	// go func() {
	// 	time.Sleep(7 * time.Second)

	// 	trLate := network.NewLocalTransport("LATE_REMOTE")
	// 	trRemoteC.Connect(trLate)
	// 	lateServer := makeServer(string(trLate.Addr()), trLate, nil)

	// 	go lateServer.Start()
	// }()

	privKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", transports[0], &privKey)
	localServer.Start()
}

func initRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		Transport: tr,
		PrivateKey: pk,
		ID: id,
		Transports: transports,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		fmt.Println(err)
	}

	return s
}

func sendGetStatusMessage(tr network.Transport, to network.NetAddr) error {
	var (
		getStatusMsg = new(network.GetStatusMessage)
		buf = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	tx := core.NewTransaction(contract())
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}

func contract() []byte {
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}   
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0xae}
	data = append(data, pushFoo...)

	return data
}