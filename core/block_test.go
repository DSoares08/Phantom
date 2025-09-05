package core

import (
	"testing"
	"time"

	"github.com/DSoares08/Phantom/crypto"
	"github.com/DSoares08/Phantom/types"
	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version: 1,
		PrevBlock: types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height: height,
	}
    tx := Transaction{
		Data: []byte("foo"),
	}

	return NewBlock(header, []Transaction{tx})
}

func randomBlockWithSignature(t *testing.T, height uint32) *Block {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(height)
	assert.Nil(t, b.Sign(privKey))

	return b
}

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())
}