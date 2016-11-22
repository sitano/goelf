package elf

import (
	"io"
	"fmt"
	"encoding/binary"
	"golang.org/x/debug/elf"
	"strings"
	"strconv"
	"bytes"
	"errors"
)

type NoteType uint32

type Note struct {
	Name string
	Type NoteType
	Data []byte

	io.ReaderAt
}

// Open returns a new ReadSeeker reading the ELF section.
func (n *Note) Open() io.ReadSeeker { return io.NewSectionReader(bytes.NewReader(n.Data), 0, int64(len(n.Data))) }

const (
	NT_GO_BUILD	NoteType = 0x4

	/*
	* Notes used in ET_CORE. Architectures export some of the arch register sets
	* using the corresponding note types via the PTRACE_GETREGSET and
	* PTRACE_SETREGSET requests.
	* http://lxr.free-electrons.com/source/include/uapi/linux/elf.h#L400
	*/
	NT_PRSTATUS     NoteType = 0x1
	NT_PRFPREG      NoteType = 0x2
	NT_PRPSINFO     NoteType = 0x3
	NT_PRXREG       NoteType = 0x4
	NT_PLATFORM     NoteType = 0x5
	NT_AUXV         NoteType = 0x6

	NT_SIGINFO      NoteType = 0x53494749
	NT_FILE         NoteType = 0x46494c45
	NT_PRXFPREG     NoteType = 0x46e62b7f /* copied from gdb5.1/include/elf/common.h */
	NT_PPC_VMX      NoteType = 0x100 /* PowerPC Altivec/VMX registers */
	NT_PPC_SPE      NoteType = 0x101 /* PowerPC SPE/EVR registers */
	NT_PPC_VSX      NoteType = 0x102 /* PowerPC VSX registers */
	NT_PPC_TAR      NoteType = 0x103 /* Target Address Register */
	NT_PPC_PPR      NoteType = 0x104 /* Program Priority Register */
	NT_PPC_DSCR     NoteType = 0x105 /* Data Stream Control Register */
	NT_PPC_EBB      NoteType = 0x106 /* Event Based Branch Registers */
	NT_PPC_PMU      NoteType = 0x107 /* Performance Monitor Registers */
	NT_PPC_TM_CGPR  NoteType = 0x108 /* TM checkpointed GPR Registers */
	NT_PPC_TM_CFPR  NoteType = 0x109 /* TM checkpointed FPR Registers */
	NT_PPC_TM_CVMX  NoteType = 0x10a /* TM checkpointed VMX Registers */
	NT_PPC_TM_CVSX  NoteType = 0x10b /* TM checkpointed VSX Registers */
	NT_PPC_TM_SPR   NoteType = 0x10c /* TM Special Purpose Registers */
	NT_PPC_TM_CTAR  NoteType = 0x10d /* TM checkpointed Target Address Register */
	NT_PPC_TM_CPPR  NoteType = 0x10e /* TM checkpointed Program Priority Register */
	NT_PPC_TM_CDSCR NoteType = 0x10f /* TM checkpointed Data Stream Control Register */
	NT_386_TLS      NoteType = 0x200 /* i386 TLS slots (struct user_desc) */
	NT_386_IOPERM   NoteType = 0x201
	NT_X86_XSTATE   NoteType = 0x202
	NT_S390_HIGH_GPRS       NoteType = 0x300 /* s390 upper register halves */
	NT_S390_TIMER   NoteType = 0x301 /* s390 timer register */
	NT_S390_TODCMP  NoteType = 0x302 /* s390 TOD clock comparator register */
	NT_S390_TODPREG NoteType = 0x303 /* s390 TOD programmable register */
	NT_S390_CTRS    NoteType = 0x304 /* s390 control registers */
	NT_S390_PREFIX  NoteType = 0x305 /* s390 prefix register */
	NT_S390_LAST_BREAK      NoteType = 0x306 /* s390 breaking event address */
	NT_S390_SYSTEM_CALL     NoteType = 0x307 /* s390 system call restart data */
	NT_S390_TDB     NoteType = 0x308 /* s390 transaction diagnostic block */
	NT_S390_VXRS_LOW        NoteType = 0x309 /* s390 vector registers 0-15 upper half */
	NT_S390_VXRS_HIGH       NoteType = 0x30a /* s390 vector registers 16-31 */
	NT_ARM_VFP      NoteType = 0x400 /* ARM VFP/NEON registers */
	NT_ARM_TLS      NoteType = 0x401 /* ARM TLS register */
	NT_ARM_HW_BREAK NoteType = 0x402 /* ARM hardware breakpoint registers */
	NT_ARM_HW_WATCH NoteType = 0x403 /* ARM hardware watchpoint registers */
	NT_ARM_SYSTEM_CALL      NoteType = 0x404 /* ARM system call number */
	NT_METAG_CBUF   NoteType = 0x500 /* Metag catch buffer registers */
	NT_METAG_RPIPE  NoteType = 0x501 /* Metag read pipeline state */
	NT_METAG_TLS    NoteType = 0x502 /* Metag TLS pointer */
)

var shnStrings = []intName{
	{0x1, "NT_PRSTATUS"},
	{0x2, "NT_PRFPREG"},
	{0x3, "NT_PRPSINFO"},
	{0x4, "NT_PRXREG"},
	{0x5, "NT_PLATFORM"},
	{0x6, "NT_AUXV"},

	{0x53494749, "NT_SIGINFO"},
	{0x46494c45, "NT_FILE"},
	{0x46e62b7f, "NT_PRXFPREG"}, /* copied from gdb5.1/include/elf/common.h */

	{0x100, "NT_PPC_VMX"}, /* PowerPC Altivec/VMX registers */
	{0x101, "NT_PPC_SPE"}, /* PowerPC SPE/EVR registers */
	{0x102, "NT_PPC_VSX"}, /* PowerPC VSX registers */
	{0x103, "NT_PPC_TAR"}, /* Target Address Register */
	{0x104, "NT_PPC_PPR"}, /* Program Priority Register */
	{0x105, "NT_PPC_DSCR"}, /* Data Stream Control Register */
	{0x106, "NT_PPC_EBB"}, /* Event Based Branch Registers */
	{0x107, "NT_PPC_PMU"}, /* Performance Monitor Registers */
	{0x108, "NT_PPC_TM_CGPR"}, /* TM checkpointed GPR Registers */
	{0x109, "NT_PPC_TM_CFPR"}, /* TM checkpointed FPR Registers */
	{0x10a, "NT_PPC_TM_CVMX"}, /* TM checkpointed VMX Registers */
	{0x10b, "NT_PPC_TM_CVSX"}, /* TM checkpointed VSX Registers */
	{0x10c, "NT_PPC_TM_SPR"}, /* TM Special Purpose Registers */
	{0x10d, "NT_PPC_TM_CTAR"}, /* TM checkpointed Target Address Register */
	{0x10e, "NT_PPC_TM_CPPR"}, /* TM checkpointed Program Priority Register */
	{0x10f, "NT_PPC_TM_CDSCR"}, /* TM checkpointed Data Stream Control Register */
	{0x200, "NT_386_TLS"}, /* i386 TLS slots (struct user_desc) */
	{0x201, "NT_386_IOPERM"},
	{0x202, "NT_X86_XSTATE"},
	{0x300, "NT_S390_HIGH_GPRS"}, /* s390 upper register halves */
	{0x301, "NT_S390_TIMER"}, /* s390 timer register */
	{0x302, "NT_S390_TODCMP"}, /* s390 TOD clock comparator register */
	{0x303, "NT_S390_TODPREG"}, /* s390 TOD programmable register */
	{0x304, "NT_S390_CTRS"}, /* s390 control registers */
	{0x305, "NT_S390_PREFIX"}, /* s390 prefix register */
	{0x306, "NT_S390_LAST_BREAK"}, /* s390 breaking event address */
	{0x307, "NT_S390_SYSTEM_CALL"}, /* s390 system call restart data */
	{0x308, "NT_S390_TDB"}, /* s390 transaction diagnostic block */
	{0x309, "NT_S390_VXRS_LOW"}, /* s390 vector registers 0-15 upper half */
	{0x30a, "NT_S390_VXRS_HIGH"}, /* s390 vector registers 16-31 */
	{0x400, "NT_ARM_VFP"}, /* ARM VFP/NEON registers */
	{0x401, "NT_ARM_TLS"}, /* ARM TLS register */
	{0x402, "NT_ARM_HW_BREAK"}, /* ARM hardware breakpoint registers */
	{0x403, "NT_ARM_HW_WATCH"}, /* ARM hardware watchpoint registers */
	{0x404, "NT_ARM_SYSTEM_CALL"}, /* ARM system call number */
	{0x500, "NT_METAG_CBUF"}, /* Metag catch buffer registers */
	{0x501, "NT_METAG_RPIPE"}, /* Metag read pipeline state */
	{0x502, "NT_METAG_TLS"}, /* Metag TLS pointer */
}

func (i NoteType) String() string   { return stringName(uint32(i), shnStrings, false) }
func (i NoteType) GoString() string { return stringName(uint32(i), shnStrings, true) }

func ReadNotes(s *elf.Section, o binary.ByteOrder) ([]*Note, error) {
	if s.Type != elf.SHT_NOTE {
		return []*Note{}, fmt.Errorf("invalid section type: %v/%v", s.Name, s.Type)
	}

	notes := []*Note{}
	r := s.Open()
	for {
		note := &Note{}

		// Copyright 2015 The Go Authors. All rights reserved.
		// Use of this source code is governed by a BSD-style
		// license that can be found in the LICENSE file.
		var namesize, descsize int32
		err := binary.Read(r, o, &namesize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read namesize failed: %v", err)
		}

		err = binary.Read(r, o, &descsize)
		if err != nil {
			return nil, fmt.Errorf("read descsize failed: %v", err)
		}

		err = binary.Read(r, o, &note.Type)
		if err != nil {
			return nil, fmt.Errorf("read type failed: %v", err)
		}
		// END

		if name, err := readAligned4(r, namesize); err != nil {
			return nil, fmt.Errorf("read name failed: %v", err)
		} else {
			note.Name = strings.TrimRight(string(name), "\x00")
		}

		note.Data, err = readAligned4(r, descsize)
		if err != nil {
			return nil, fmt.Errorf("read desc failed: %v", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func ReadNoteByType(s *elf.Section, o binary.ByteOrder, search NoteType) (*Note, error) {
	if s.Type != elf.SHT_NOTE {
		return nil, fmt.Errorf("invalid section type: %v/%v", s.Name, s.Type)
	}

	note := &Note{}

	r := s.Open()
	for {
		var namesize, descsize int32
		err := binary.Read(r, o, &namesize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read namesize failed: %v", err)
		}

		err = binary.Read(r, o, &descsize)
		if err != nil {
			return nil, fmt.Errorf("read descsize failed: %v", err)
		}

		err = binary.Read(r, o, &note.Type)
		if err != nil {
			return nil, fmt.Errorf("read type failed: %v", err)
		}

		if note.Type != search {
			sz := int64(namesize + descsize)
			full := (sz + 3) &^ 3
			_, _ = r.Seek(full, io.SeekCurrent)

			continue
		}
		// END

		if name, err := readAligned4(r, namesize); err != nil {
			return nil, fmt.Errorf("read name failed: %v", err)
		} else {
			note.Name = strings.TrimRight(string(name), "\x00")
		}

		note.Data, err = readAligned4(r, descsize)
		if err != nil {
			return nil, fmt.Errorf("read desc failed: %v", err)
		}

		return note, nil
	}

	return nil, errors.New("not found")
}

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func readAligned4(r io.Reader, sz int32) ([]byte, error) {
	full := (sz + 3) &^ 3
	data := make([]byte, full)
	_, err := io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}
	data = data[:sz]
	return data, nil
}

/*
 * ELF constants and data structures
 *
 * Derived from:
 * $FreeBSD: src/sys/sys/elf32.h,v 1.8.14.1 2005/12/30 22:13:58 marcel Exp $
 * $FreeBSD: src/sys/sys/elf64.h,v 1.10.14.1 2005/12/30 22:13:58 marcel Exp $
 * $FreeBSD: src/sys/sys/elf_common.h,v 1.15.8.1 2005/12/30 22:13:58 marcel Exp $
 * $FreeBSD: src/sys/alpha/include/elf.h,v 1.14 2003/09/25 01:10:22 peter Exp $
 * $FreeBSD: src/sys/amd64/include/elf.h,v 1.18 2004/08/03 08:21:48 dfr Exp $
 * $FreeBSD: src/sys/arm/include/elf.h,v 1.5.2.1 2006/06/30 21:42:52 cognet Exp $
 * $FreeBSD: src/sys/i386/include/elf.h,v 1.16 2004/08/02 19:12:17 dfr Exp $
 * $FreeBSD: src/sys/powerpc/include/elf.h,v 1.7 2004/11/02 09:47:01 ssouhlal Exp $
 * $FreeBSD: src/sys/sparc64/include/elf.h,v 1.12 2003/09/25 01:10:26 peter Exp $
 *
 * Copyright (c) 1996-1998 John D. Polstra.  All rights reserved.
 * Copyright (c) 2001 David E. O'Brien
 * Portions Copyright 2009 The Go Authors.  All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

type intName struct {
	i uint32
	s string
}

func stringName(i uint32, names []intName, goSyntax bool) string {
	for _, n := range names {
		if n.i == i {
			if goSyntax {
				return "elf." + n.s
			}
			return n.s
		}
	}

	// second pass - look for smaller to add with.
	// assume sorted already
	for j := len(names) - 1; j >= 0; j-- {
		n := names[j]
		if n.i < i {
			s := n.s
			if goSyntax {
				s = "elf." + s
			}
			return s + "+" + strconv.FormatUint(uint64(i-n.i), 10)
		}
	}

	return strconv.FormatUint(uint64(i), 10)
}

func flagName(i uint32, names []intName, goSyntax bool) string {
	s := ""
	for _, n := range names {
		if n.i&i == n.i {
			if len(s) > 0 {
				s += "+"
			}
			if goSyntax {
				s += "elf."
			}
			s += n.s
			i -= n.i
		}
	}
	if len(s) == 0 {
		return "0x" + strconv.FormatUint(uint64(i), 16)
	}
	if i != 0 {
		s += "+0x" + strconv.FormatUint(uint64(i), 16)
	}
	return s
}