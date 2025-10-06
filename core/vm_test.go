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
	// 1 + 2 = 3
	// 1
	// push stack
	// 2
	// push stack
	// add
	// 3
	// push stack

	data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}
	// data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())

	// result := vm.queue.Pop()
	// assert.Equal(t, 4, result)

	result := vm.queue.Pop().(int)

	assert.Equal(t, 1, result)

	// assert.Equal(t, "FOO", string(result))
}