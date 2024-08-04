package cpu

import (
	"encoding/binary"
	"fmt"
)

func (cpu *CpuState) getEffectiveAddress(opcode uint16, size uint) Mapper {
	var effectiveMapper Mapper

	switch (opcode >> 3) & 0x7 {
	case 0x0:
		effectiveMapper = &RegisterMapper{reg: &cpu.DataRegister[opcode&0x7]}
	case 0x1:
		effectiveMapper = &RegisterMapper{reg: &cpu.AddressRegister[opcode&0x7]}
	case 0x2:
		effectiveMapper = &MemoryMapper{mem: cpu.Memory[cpu.AddressRegister[opcode&0x7]:]}
	case 0x3:
		effectiveMapper = &MemoryMapper{mem: cpu.Memory[cpu.AddressRegister[opcode&0x7]:]}
		cpu.AddressRegister[opcode&0x7] += uint64(size)
	case 0x4:
		effectiveMapper = &MemoryMapper{mem: cpu.Memory[cpu.AddressRegister[opcode&0x7]:]}
		cpu.AddressRegister[opcode&0x7] -= uint64(size)
	case 0x5:
		displacement := binary.BigEndian.Uint16(cpu.Memory[cpu.PC+2 : cpu.PC+4])
		effectiveMapper = &MemoryMapper{mem: cpu.Memory[cpu.AddressRegister[int16(opcode&0x7)+int16(displacement)]:]}
	case 0x6:
		wordExtension := binary.BigEndian.Uint16(cpu.Memory[cpu.PC+2 : cpu.PC+4])
		regType := wordExtension & 0x8000 >> 15
		regNumber := (wordExtension & 0x7000) >> 12
		extensionSize := wordExtension & 0x800 >> 11
		displacement := wordExtension & 0xff
		var extensionMask uint64

		ea := cpu.AddressRegister[opcode&0x7]

		if extensionSize == 0x1 {
			extensionMask = 0xffff
		} else {
			extensionMask = 0xfffffff
		}

		if regType == 0x1 {
			ea += cpu.AddressRegister[regNumber] & extensionMask
		} else {
			ea += cpu.DataRegister[regNumber] & extensionMask
		}

		ea += uint64(displacement)
		effectiveMapper = &MemoryMapper{mem: cpu.Memory[ea:]}
	case 0x7:
		switch opcode & 0x7 {
		case 0x0:
			ea := int(binary.BigEndian.Uint16(cpu.Memory[cpu.PC+2 : cpu.PC+4]))

			if ea < 0 {
				effectiveMapper = &MemoryMapper{mem: cpu.Memory[len(cpu.Memory)-ea:]}
			} else {
				effectiveMapper = &MemoryMapper{mem: cpu.Memory[ea:]}
			}
		case 0x1:
			ea := int(binary.BigEndian.Uint32(cpu.Memory[cpu.PC+2 : cpu.PC+6]))

			if ea < 0 {
				effectiveMapper = &MemoryMapper{mem: cpu.Memory[len(cpu.Memory)-ea:]}
			} else {
				effectiveMapper = &MemoryMapper{mem: cpu.Memory[ea:]}
			}
		case 0x2:
			displacement := binary.BigEndian.Uint16(cpu.Memory[cpu.PC+2 : cpu.PC+4])
			effectiveMapper = &MemoryMapper{mem: cpu.Memory[cpu.PC+uint64(displacement):]}
		case 0x3:
			wordExtension := binary.BigEndian.Uint16(cpu.Memory[cpu.PC+2 : cpu.PC+4])
			regType := wordExtension & 0x8000 >> 15
			regNumber := (wordExtension & 0x7000) >> 12
			extensionSize := wordExtension & 0x800 >> 11
			displacement := wordExtension & 0xff
			var extensionMask uint64

			ea := cpu.PC

			if extensionSize == 0x1 {
				extensionMask = 0xffff
			} else {
				extensionMask = 0xfffffff
			}

			if regType == 0x1 {
				ea += cpu.AddressRegister[regNumber] & extensionMask
			} else {
				ea += cpu.DataRegister[regNumber] & extensionMask
			}

			ea += uint64(displacement)
			effectiveMapper = &MemoryMapper{mem: cpu.Memory[ea:]}
		case 0x4:
			var wordExtension uint64

			switch size {
			case BYTE_SIZE:
				wordExtension = binary.BigEndian.Uint64(cpu.Memory[cpu.PC+2 : cpu.PC+4])
			case WORD_SIZE:
				wordExtension = binary.BigEndian.Uint64(cpu.Memory[cpu.PC+2 : cpu.PC+4])
			case LONG_SIZE:
				wordExtension = binary.BigEndian.Uint64(cpu.Memory[cpu.PC+2 : cpu.PC+4])
			}

			effectiveMapper = &RegisterMapper{reg: &wordExtension}
		}
	}

	return effectiveMapper
}

func (cpu *CpuState) DecodeInstruction() {
	opcode := binary.BigEndian.Uint16(cpu.Memory[cpu.PC : cpu.PC+2])

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

			addressMap := MemoryMapper{mem: cpu.Memory[cpu.AddressRegister[addressIndex]:4]}
			dataMap := RegisterMapper{reg: &cpu.DataRegister[dataIndex]}

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
