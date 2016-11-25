package elf

import (
	"encoding/binary"
	"fmt"
	"io"
	"errors"
	"golang.org/x/debug/elf"
	"reflect"
	"unsafe"
)

type ElfSigInfo struct {
	Sig  int32 /* signal number */
	Code int32 /* extra code */
	Err  int32 /* errno */
}

/*
  * Definitions to generate Intel SVR4-like core files.
  * These mostly have the same names as the SVR4 types with "elf_"
  * tacked on the front to prevent clashes with linux definitions,
  * and the typedef forms have been avoided.  This is mostly like
  * the SVR4 structure, but more Linuxy, with things that Linux does
  * not support and which gdb doesn't really use excluded.
  * Fields present but not used are marked with "XXX".
  *
  * http://lxr.free-electrons.com/source/include/uapi/linux/elfcore.h#L36
  * https://llvm.org/svn/llvm-project/lldb/trunk/source/Plugins/Process/elf-core/ProcessElfCore.cpp
  */
type PRStatus struct {
	Info    ElfSigInfo /* Info associated with signal */

	CurSig  int16      /* Current signal */
	SigPend uint       /* Set of pending signals */
	SigHold uint       /* Set of held signals */

	PID     KernelPid
	PPID    KernelPid
	PGRP    KernelPid
	SID     KernelPid

	UTime   TimeVal    /* User time */
	STime   TimeVal    /* System time */
	CUTime  TimeVal    /* Cumulative user time */
	CSTime  TimeVal    /* Cumulative system time */

	Regs    ElfGRegSet /* GP registers */
}

func GetUserRegs(set ElfGRegSet) UserRegs {
	r := UserRegs{}

	v := reflect.ValueOf(r)
	for i := 0; i < v.NumField(); i ++ {
		bits := v.Field(i).Type().Bits()
		if bits == 32 {
			p := unsafe.Pointer(&r)
			f := unsafe.Pointer(uintptr(p) + uintptr(i * 4))
			*(*uint64)(f) = uint64(set[i])
		} else if bits == 64 {
			p := unsafe.Pointer(&r)
			f := unsafe.Pointer(uintptr(p) + uintptr(i * 8))
			*(*uint64)(f) = uint64(set[i])
		} else {
        	panic("unknown class")
		}
	}

	return r
}

func ReadKernelPid(r io.Reader, o binary.ByteOrder) (KernelPid, error) {
	x, err := readInt(r, o, elf.ELFCLASS32)
	return KernelPid(x), err
}

func readInt(r io.Reader, o binary.ByteOrder, c elf.Class) (int, error) {
	if c == elf.ELFCLASS64 {
		var x int64
		err := binary.Read(r, o, &x)
		return int(x), err
	} else if c == elf.ELFCLASS32 {
		var x int32
		err := binary.Read(r, o, &x)
		return int(x), err
	} else {
		return 0, errors.New("unknown elf class")
	}
}

func readUInt(r io.Reader, o binary.ByteOrder, c elf.Class) (uint, error) {
	if c == elf.ELFCLASS64 {
		var x uint64
		err := binary.Read(r, o, &x)
		return uint(x), err
	} else if c == elf.ELFCLASS32 {
		var x uint32
		err := binary.Read(r, o, &x)
		return uint(x), err
	} else {
		return 0, errors.New("unknown elf class")
	}
}

func ReadPRStatus(n *Note, o binary.ByteOrder, c elf.Class) (*PRStatus, error) {
	if n.Type != NT_PRSTATUS {
		return nil, fmt.Errorf("invalid note type: %v", n)
	}

	var err error
	prs := &PRStatus{}

	r := n.Open()

	if err = binary.Read(r, o, &prs.Info.Sig); err != nil {
		return nil, fmt.Errorf("read sig failed: %v", err)
	}
	if err = binary.Read(r, o, &prs.Info.Code); err != nil {
		return nil, fmt.Errorf("read code failed: %v", err)
	}
	if err = binary.Read(r, o, &prs.Info.Err); err != nil {
		return nil, fmt.Errorf("read err failed: %v", err)
	}

	if err = binary.Read(r, o, &prs.CurSig); err != nil {
		return nil, fmt.Errorf("read cursig failed: %v", err)
	}
	r.Seek(2, io.SeekCurrent)

	if prs.SigPend, err = readUInt(r, o, c); err != nil {
		return nil, fmt.Errorf("read sigpend failed: %v", err)
	}
	if prs.SigHold, err = readUInt(r, o, c); err != nil {
		return nil, fmt.Errorf("read sighold failed: %v", err)
	}

	if prs.PID, err = ReadKernelPid(r, o); err != nil {
		return nil, fmt.Errorf("read pid failed: %v", err)
	}
	if prs.PPID, err = ReadKernelPid(r, o); err != nil {
		return nil, fmt.Errorf("read ppid failed: %v", err)
	}
	if prs.PGRP, err = ReadKernelPid(r, o); err != nil {
		return nil, fmt.Errorf("read pgrp failed: %v", err)
	}
	if prs.SID, err = ReadKernelPid(r, o); err != nil {
		return nil, fmt.Errorf("read sid failed: %v", err)
	}

	if err = binary.Read(r, o, &prs.UTime.Sec); err != nil {
		return nil, fmt.Errorf("read utime.sec failed: %v", err)
	}
	if err = binary.Read(r, o, &prs.UTime.USec); err != nil {
		return nil, fmt.Errorf("read utime.usec failed: %v", err)
	}

	if err = binary.Read(r, o, &prs.STime.Sec); err != nil {
		return nil, fmt.Errorf("read stime.sec failed: %v", err)
	}
	if err = binary.Read(r, o, &prs.STime.USec); err != nil {
		return nil, fmt.Errorf("read stime.usec failed: %v", err)
	}

	if err = binary.Read(r, o, &prs.CUTime.Sec); err != nil {
		return nil, fmt.Errorf("read cutime.sec failed: %v", err)
	}
	if err = binary.Read(r, o, &prs.CUTime.USec); err != nil {
		return nil, fmt.Errorf("read cutime.usec failed: %v", err)
	}

	if err = binary.Read(r, o, &prs.CSTime.Sec); err != nil {
		return nil, fmt.Errorf("read cstime.sec failed: %v", err)
	}
	if err = binary.Read(r, o, &prs.CSTime.USec); err != nil {
		return nil, fmt.Errorf("read cstime.usec failed: %v", err)
	}

	if c == elf.ELFCLASS64 {
		for i := uintptr(0); i < ELF_NGREG; i ++ {
			var x uint64
			if err := binary.Read(r, o, &x); err != nil {
				return nil, fmt.Errorf("read %d/%d reg failed: %v", 1 + i, ELF_NGREG, err)
			}
			prs.Regs[i] = ElfGReg(x)
		}
	} else if c == elf.ELFCLASS32 {
		for i := uintptr(0); i < ELF_NGREG; i ++ {
			var x uint32
			if err := binary.Read(r, o, &x); err != nil {
				return nil, fmt.Errorf("read %d/%d reg failed: %v", 1 + i, ELF_NGREG, err)
			}
			prs.Regs[i] = ElfGReg(x)
		}
	} else {
		return nil, errors.New("unknown elf class")
	}

	return prs, nil
}

