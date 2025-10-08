package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := NewQueue(128)

	q.Push(1)
	q.Push(2)

	value := q.Pop()
	assert.Equal(t, 1, value)

	value = q.Pop()
	assert.Equal(t, 2, value)
}

func TestVM(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f} //, 0x4f, 0x0c, 0x0d} //, 0x0e}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	valueBytes, err := contractState.Get([]byte("FOO"))
	value := deserializeInt64(valueBytes)
	assert.Nil(t, err)
	assert.Equal(t, value, int64(5))
}