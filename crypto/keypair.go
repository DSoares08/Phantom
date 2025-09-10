package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/DSoares08/Phantom/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r, 
		S: s,
	}, nil
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{key: key}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{Key: &k.key.PublicKey}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

// GobEncode serializes the public key as compressed bytes, avoiding gob traversing
// the unexported fields of elliptic.Curve implementations.
func (k PublicKey) GobEncode() ([]byte, error) {
	if k.Key == nil {
		return nil, nil
	}
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y), nil
}

// GobDecode restores the public key from its compressed form using P-256.
func (k *PublicKey) GobDecode(data []byte) error {
	if len(data) == 0 {
		k.Key = nil
		return nil
	}
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), data)
	if x == nil || y == nil {
		return errors.New("invalid public key bytes")
	}
	k.Key = &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	return nil
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())

	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	S, R *big.Int
}

func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}
