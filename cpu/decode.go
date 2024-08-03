package cpu

import "fmt"

func (cpu *CpuState) DecodeInstruction() {
	opcode := (uint(cpu.Memory[cpu.PC]) << 8) | uint(cpu.Memory[cpu.PC+1])

	switch opcode >> 12 {
	case 0x0:
		if ((opcode >> 3) & 0x7) == 0x1 {
			addressIndex := opcode & 0x7
			dataIndex := (opcode >> 9) & 0x7
			Quantifier := (opcode >> 6) & 0x1
			var byteSize uint

			if Quantifier == 0 {
				byteSize = WORD_SIZE
			} else {
				byteSize = LONG_SIZE
			}

			addressMap := MemoryMapper{mem: cpu.Memory[*cpu.AddressRegister[addressIndex]:4]}
			dataMap := RegisterMapper{reg: cpu.DataRegister[dataIndex]}

			// For now there is no 16bit displacement, needs to be implemented
			if (opcode>>7)&0x1 == 0x1 {
				cpu.Movep(&addressMap, &dataMap, byteSize)
			} else {
				cpu.Movep(&dataMap, &addressMap, byteSize)
			}

			cpu.PC += 4
		} else if ((opcode>>8)&0x1) == 0x1 || (opcode>>8) == 0x8 {
			//decodeBitOperation
		} else {
			//decodeAritmthicImeddiate
		}
	case 0x1, 0x2, 0x3:
		// Move and Movea
	default:
		panic(fmt.Sprintf("Invalid opcode: %x", opcode))
	}
}
