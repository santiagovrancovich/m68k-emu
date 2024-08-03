package cpu

// Reference instruction implementation
func (cpu *CpuState) Add(destination, source Mapper, size uint) {
	var mask uint64 = (0x1 << (size * 8)) - 1

	destinationValue, err := destination.Read(size)
	sourceValue, err := source.Read(size)

	if err != nil {
		panic(err.Error())
	}

	operation := destinationValue + destinationValue
	signedOperation := int64(destinationValue) + int64(sourceValue)

	negativeFlag := (operation >> ((size * 8) - 1))
	carryFlag := (operation >> (size * 8))
	var zeroFlag uint64 = 0
	var overflowFlag uint64 = 0

	if (operation & mask) == 0 {
		zeroFlag = 1
	}

	if operation != uint64(signedOperation) {
		overflowFlag = 1
	}

	var newCCR uint16 = uint16(negativeFlag) << 3 & uint16(overflowFlag) << 2 & uint16(zeroFlag) << 1 & uint16(carryFlag)
	cpu.SR = (cpu.SR & 0b10000) | newCCR

	destination.Write(operation, size)
}

func (cpu *CpuState) Movep(destination, source Mapper, size uint) {
	//Missing implentation
}
