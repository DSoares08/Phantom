package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/DSoares08/Phantom/types"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Decode(t *testing.T) {
	h := &Header{
		Version: 1,
		PrevBlock: types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height: 10,
		Nonce: 989394,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}