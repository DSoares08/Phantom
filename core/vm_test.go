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
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}   
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0xae}

	data = append(data, pushFoo...)

	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	// fmt.Printf("%+v", vm.queue.data)
	value := vm.queue.Pop().([]byte)
	valueSerialized := deserializeInt64(value)

	assert.Equal(t, valueSerialized, int64(5))

	// valueBytes, err := contractState.Get([]byte("FOO"))
	// assert.Nil(t, err)
	// assert.Equal(t, value, int64(5))
}