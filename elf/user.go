package elf

import (
	"unsafe"
)

// http://lxr.free-electrons.com/source/arch/x86/include/asm/elf.h#L16
type ElfGReg uint

const ELF_NGREG = unsafe.Sizeof(UserRegs{}) / unsafe.Sizeof(ElfGReg(0))

type ElfGRegSet [ELF_NGREG]ElfGReg
