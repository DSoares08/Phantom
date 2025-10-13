package network

import (
	"bytes"
	"net"
	"os"
	"fmt"
	"time"

	"github.com/DSoares08/Phantom/types"
	"github.com/DSoares08/Phantom/crypto"
	"github.com/DSoares08/Phantom/core"
	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	TCPTransport *TCPTransport
	ID string
	Logger log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor RPCProcessor
	BlockTime time.Duration
	PrivateKey *crypto.PrivateKey
}

type Server struct{
	peerCh chan *TCPPeer
	peerMap map[net.Addr]*TCPPeer
	ServerOpts
	memPool *TxPool
	chain *core.Blockchain
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		// opts.Logger = log.With(opts.Logger, "addr", opts.Transport.Addr())
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	peerCh := make(chan *TCPPeer)
	s := &Server{
		peerCh: peerCh,
		peerMap: make(map[net.Addr]*TCPPeer),
		ServerOpts: opts,
		chain: chain,
			memPool: NewTxPool(1000),
			isValidator: opts.PrivateKey != nil,
			rpcCh:      make(chan RPC),
			quitCh:     make(chan struct{}, 1),
	}

	s.TCPTransport.peerCh = peerCh

	// Using server as default
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) Start() {
	s.TCPTransport.Start()

free:
	for {
		select {
		case peer := <-s.peerCh:
			fmt.Printf("new peer => %+v\n", peer)
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				fmt.Println(s.ID, err)
			}

			if err := s.ProcessMessage(msg); err != nil {
				if err != core.ErrBlockKnown {
					s.Logger.Log("error", err)
				}
			}
		case <-s.quitCh:
			break free
		}
	}

	fmt.Println(s.ID, "Server shutdown")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)

	fmt.Println(s.ID, "Starting validator loop", "blockTime", s.BlockTime)
	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	case *core.Block:
		return s.processBlock(t)
	case *GetStatusMessage:
		// return s.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		// return s.processStatusMessage(msg.From, t)
	case *GetBlocksMessage:
		return s.processGetBlocksMessage(msg.From, t)
	}

	return nil
}

func (s *Server) processGetBlocksMessage(from NetAddr, data *GetBlocksMessage) error {
	fmt.Printf("got get blocks message => %+v\n", data)

	return nil
}

// TODO: Remove the logic from the main function to here
// func (s *Server) sendGetStatusMessage(tr Transport) error {
// 	var (
// 		getStatusMsg = new(GetStatusMessage)
// 		buf = new(bytes.Buffer)
// 	)
// 	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
// 		return err
// 	}

// 	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())

// 	if err := s.Transport.SendMessage(tr.Addr(), msg.Bytes()); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *Server) broadcast(payload []byte) error {
	// for _, tr := range s.Transports {
	// 	if err := tr.Broadcast(payload); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// func (s *Server) processStatusMessage(from NetAddr, data *StatusMessage) error {
// 	if data.CurrentHeight <= s.chain.Height() {
// 		s.Logger.Log("msg", "cannot sync blockHeight too low", "ourHeight", s.chain.Height(), "theirHeight", data.CurrentHeight, "addr", from)
// 		return nil
// 	}

// 	getBlocksMessage := &GetBlocksMessage{
// 		From: s.chain.Height(),
// 		To: 0,
// 	}

// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(getBlocksMessage); err != nil {
// 		return err
// 	}

// 	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())

// 	// Send getBlocksMessage
// 	return s.Transport.SendMessage(from, msg.Bytes())

// }

// func (s *Server) processGetStatusMessage(from NetAddr, data *GetStatusMessage) error {
// 	fmt.Printf("=> received GetStatus msg from %s => %+v\n", from, data)

// 	statusMessage := &StatusMessage{
// 		CurrentHeight: s.chain.Height(),
// 		ID: s.ID,
// 	}

// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(statusMessage); err != nil {
// 		return err
// 	}

// 	msg := NewMessage(MessageTypeStatus, buf.Bytes())

// 	return s.Transport.SendMessage(from, msg.Bytes())
// }

func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.AddBlock(b); err != nil {
		return err
	}

	go s.broadcastBlock(b)

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Contains(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	// fmt.Println(s.ID, "adding new tx to the mempool", hash, s.memPool.PendingCount())

	go s.broadcastTx(tx)
	
	s.memPool.Add(tx)

	return nil
}

func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	txx := s.memPool.Pending()

	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	s.memPool.ClearPending()

	go s.broadcastBlock(block)

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version: 1,
		DataHash: types.Hash{},
		Timestamp: 000000,
		Height: 0,
	}

	b, _ := core.NewBlock(header, nil)
	return b
}