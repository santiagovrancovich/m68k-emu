package cpu

import (
	"encoding/binary"
	"errors"
	"math"
)

type MemoryMapper struct {
	mem []byte
}

type RegisterMapper struct {
	reg *uint64
}

type Mapper interface {
	Read(size uint) (uint64, error)
	Write(value uint64, size uint) (uint, error)
}

func (m *MemoryMapper) Read(s uint) (uint64, error) {
	switch s {
	case BYTE_SIZE:
		return uint64(m.mem[0]), nil
	case WORD_SIZE:
		return uint64(binary.BigEndian.Uint16(m.mem[0:2])), nil
	case LONG_SIZE:
		return uint64(binary.BigEndian.Uint32(m.mem[:])), nil
	default:
		return 0, errors.New("Invalid read size")
	}
}

func (m *MemoryMapper) Write(value uint64, size uint) (uint, error) {
	switch size {
	case BYTE_SIZE:
		m.mem[0] = byte(value)
	case WORD_SIZE:
		binary.BigEndian.PutUint16(m.mem, uint16(value&math.MaxUint16))
	case LONG_SIZE:
		binary.BigEndian.PutUint32(m.mem, uint32(value&math.MaxUint32))
	default:
		return 0, errors.New("Invalid read size")
	}

	return size, nil
}

func (m *RegisterMapper) Read(size uint) (uint64, error) {
	return *m.reg & ((1 << (size * 8)) - 1), nil
}

func (m *RegisterMapper) Write(value uint64, size uint) (uint, error) {
	*m.reg = value & ((1 << (size * 8)) - 1)
	return size, nil
}
