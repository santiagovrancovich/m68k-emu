package cpu

const (
	BYTE_SIZE = 1
	WORD_SIZE = 2
	LONG_SIZE = 4
)

type CpuState struct {
	DataRegister    [8]uint64
	AddressRegister [8]uint64
	SSP             uint64
	PC              uint64
	SR              uint16
	Memory          []byte
}
