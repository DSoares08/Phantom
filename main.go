package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/DSoares08/Phantom/core"
	"github.com/DSoares08/Phantom/crypto"
	"github.com/DSoares08/Phantom/network"
)

func main() {
	privKey := crypto.GeneratePrivateKey()
	localNode := makeServer("LOCAL_NODE", &privKey, ":3000", []string{":4000"}, ":8080")
	go localNode.Start()

	remoteNode := makeServer("REMOTE_NODE", nil, ":4000", []string{":5000"}, "")
	go remoteNode.Start()

	remoteNodeB := makeServer("REMOTE_NODE_B", nil, ":5000", nil, "")
	go remoteNodeB.Start()

	go func() {
		time.Sleep(11 * time.Second)

		lateNode := makeServer("LATE_NODE", nil, ":6000", []string{":4000"}, "")
		go lateNode.Start()
	}()

	time.Sleep(1 * time.Second)

	if err := txSender(); err != nil {
		panic(err)
	}

	// txSendTicker := time.NewTicker(1 * time.Second)
	// go func() {
	// 	for {
	// 		txSender()

	// 		<-txSendTicker.C
	// 	}
	// }()

	select {}
}

func makeServer(id string, pk *crypto.PrivateKey, addr string, seedNodes []string, apiListenAddr string) *network.Server {
	opts := network.ServerOpts{
		APIListenAddr: apiListenAddr,
		SeedNodes: seedNodes,
		ListenAddr: addr,
		PrivateKey: pk,
		ID: id,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		fmt.Println(err)
	}

	return s
}

func txSender() error {
	privKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	tx := core.NewTransaction(nil)
	tx.To = toPrivKey.PublicKey()
	tx.Value = 100

	if err := tx.Sign(privKey); err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/tx", buf)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	_, err = client.Do(req)
	return err
}
