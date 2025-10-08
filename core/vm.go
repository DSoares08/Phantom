package core

import (
	"encoding/binary"
)

type Instruction byte

const (
	InstrPushInt Instruction = 0x0a // 10
	InstrAdd Instruction = 0x0b // 11
	InstrPushByte Instruction = 0x0c 
	InstrPack Instruction = 0x0d
	InstrSub Instruction = 0x0e // 1
	InstrStore Instruction = 0x0f 
)

type Queue struct {
	data []any
	head int	
}

func NewQueue(size int) *Queue {
	return &Queue{
		data: make([]any, size),
		head: 0,
	}
}

func (q *Queue) Push(v any) {
	q.data[q.head] = v
	q.head++
}


func (q *Queue) Pop() any {
	value := q.data[0]	
	q.data = append(q.data[:0], q.data[1:]...)
	q.head--

	return value
}

type VM struct {
	data []byte
	ip int // instruction pointer
	queue *Queue
	contractState *State
}

func NewVM(data []byte, contractState *State) *VM {
	return &VM{
		data: data,
		ip: 0,
		queue: NewQueue(128),
		contractState: contractState,
	}
}

func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])

		if err := vm.Exec(instr); err != nil {
			return err
		}

		vm.ip++

		if vm.ip > len(vm.data)-1 {
			break
		}
	}
	return nil
}

func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrStore:
		var (
			key = vm.queue.Pop().([]byte)
			value = vm.queue.Pop()
			serializedValue []byte
		)

		switch v := value.(type) {
		case int:
			serializedValue = serializeInt64(int64(v))
		default:
			panic("TODO: unkown type")
		}

		vm.contractState.Put(key, serializedValue)

	case InstrPushInt:
		vm.queue.Push(int(vm.data[vm.ip - 1]))

	case InstrPushByte:
		vm.queue.Push(byte(vm.data[vm.ip - 1]))

	case InstrPack:
		n := vm.queue.Pop().(int)
		b := make([]byte, n)

		for i := 0; i < n; i++ {
			b[i] = vm.queue.Pop().(byte)
		}

		vm.queue.Push(b)

	case InstrSub:
		a := vm.queue.Pop().(int)
		b := vm.queue.Pop().(int)
		c := a - b
		vm.queue.Push(c)

	case InstrAdd:
		a := vm.queue.Pop().(int)
		b := vm.queue.Pop().(int)
		c := a + b
		vm.queue.Push(c)
	}
	return nil
}

func serializeInt64(value int64) []byte {
	buf := make([]byte, 8)

	binary.LittleEndian.PutUint64(buf, uint64(value))

	return buf
}

func deserializeInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}