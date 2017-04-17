package elf

type UserRegs struct {
	R15 uint64
	R14 uint64
	R13 uint64
	R12 uint64
	BP uint64
	BX uint64
	R11 uint64
	R10 uint64
	R9 uint64
	R8 uint64
	AX uint64
	CX uint64
	DX uint64
	SI uint64
	DI uint64
	OrigAX uint64
	IP uint64
	CS uint64
	Flags uint64
	SP uint64
	SS uint64
	FSBase uint64
	GSBase uint64
	DS uint64
	ES uint64
	FS uint64
	GS uint64
}
