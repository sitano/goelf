package elf

import (
	"fmt"
	"encoding/binary"

	"golang.org/x/debug/elf"
	"strings"
	"io"
	"errors"
)

const ELF_PRARGSZ = 80 /* Number of chars for args */

type KernelUid uint32
type KernelGid uint32

type KernelPid uint32

// Linux kernel data sizes: https://static.lwn.net/images/pdf/LDD3/ch11.pdf
// http://lxr.free-electrons.com/source/include/uapi/linux/elfcore.h#L78
type PRPSInfo struct {
	State byte /* numeric process state */
	SName string /* char for pr_state (byte) */
	Zomb byte  /* zombie */
	Nice byte  /* nice val */

	Flag uint  /* flags */

	UID KernelUid
	GID KernelGid
	PID, PPID, PGRP, SID KernelPid

	FName string /* filename of executable ([16]byte) */
	PSArgs string /* initial part of arg list ([ELF_PRARGSZ]byte) */
}

func readKernelUid(r io.Reader, o binary.ByteOrder, c elf.Class) (KernelUid, error) {
	if c == elf.ELFCLASS64 {
		var x uint32
		err := binary.Read(r, o, &x)
		return KernelUid(x), err
	} else if c == elf.ELFCLASS32 {
		var x uint16
		err := binary.Read(r, o, &x)
		return KernelUid(x), err
	} else {
		return 0, errors.New("unknown elf class")
	}
}

func readKernelGid(r io.Reader, o binary.ByteOrder, c elf.Class) (KernelGid, error) {
	if c == elf.ELFCLASS64 {
		var x uint32
		err := binary.Read(r, o, &x)
		return KernelGid(x), err
	} else if c == elf.ELFCLASS32 {
		var x uint16
		err := binary.Read(r, o, &x)
		return KernelGid(x), err
	} else {
		return 0, errors.New("unknown elf class")
	}
}

func ReadPRPSInfo(n *Note, o binary.ByteOrder, c elf.Class) (*PRPSInfo, error) {
	if n.Type != NT_PRPSINFO {
		return nil, fmt.Errorf("invalid note type: %v", n)
	}

	var err error
	prps := &PRPSInfo{}

	r := n.Open()

	if err = binary.Read(r, o, &prps.State); err != nil {
		return nil, fmt.Errorf("read state failed: %v", err)
	}

	var sname byte
	if err = binary.Read(r, o, &sname); err != nil {
		return nil, fmt.Errorf("read sname failed: %v", err)
	} else {
		prps.SName = string(sname)
	}

	if err = binary.Read(r, o, &prps.Zomb); err != nil {
		return nil, fmt.Errorf("read zomb failed: %v", err)
	}
	if err = binary.Read(r, o, &prps.Nice); err != nil {
		return nil, fmt.Errorf("read nice failed: %v", err)
	}

	if prps.Flag, err = readUInt(r, o, c); err != nil {
		return nil, fmt.Errorf("read flag failed: %v", err)
	}
	if c == elf.ELFCLASS64 {
		r.Seek(4, io.SeekCurrent)
	}

	if prps.UID, err = readKernelUid(r, o, c); err != nil {
		return nil, fmt.Errorf("read uid failed: %v", err)
	}
	if prps.GID, err = readKernelGid(r, o, c); err != nil {
		return nil, fmt.Errorf("read gid failed: %v", err)
	}

	if prps.PID, err = ReadKernelPid(r, o); err != nil {
		return nil, fmt.Errorf("read pid failed: %v", err)
	}
	if prps.PPID, err = ReadKernelPid(r, o); err != nil {
		return nil, fmt.Errorf("read ppid failed: %v", err)
	}
	if prps.PGRP, err = ReadKernelPid(r, o); err != nil {
		return nil, fmt.Errorf("read pgrp failed: %v", err)
	}
	if prps.SID, err = ReadKernelPid(r, o); err != nil {
		return nil, fmt.Errorf("read sid failed: %v", err)
	}

	var fname [16]byte
	if err = binary.Read(r, o, &fname); err != nil {
		return nil, fmt.Errorf("read fname failed: %v", err)
	} else {
		prps.FName = strings.Trim(string(fname[:]), "\x00")
	}

	var psargs [ELF_PRARGSZ]byte
	if err = binary.Read(r, o, &psargs); err != nil {
		return nil, fmt.Errorf("read psargs failed: %v", err)
	} else {
		prps.PSArgs = strings.Trim(string(psargs[:]), "\x00")
	}

	return prps, nil
}

