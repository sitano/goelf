package main

import (
	"os"

	"golang.org/x/debug/elf"
	"golang.org/x/debug/dwarf"
)

type Process struct {
	path string

	efd *elf.File
	dwf *dwarf.Data
}

func New(path string) (*Process, error) {
	var err error

	p := &Process{}
	if p.efd, err = Open(path); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Process) DWARF() (*dwarf.Data, error) {
	var err error

	if p.dwf == nil {
		if p.dwf, err = p.efd.DWARF(); err != nil {
			return nil, err
		}
	}

	return  p.dwf, err
}

func Open(path string) (*elf.File, error) {
	fd, err := os.OpenFile(path, 0, os.ModePerm)
	if err != nil {
		return nil, err
	}

	efd, err := elf.NewFile(fd)
	if err != nil {
		return nil, err
	}

	return efd, nil
}
