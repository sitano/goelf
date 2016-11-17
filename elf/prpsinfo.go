package elf

const ELF_PRARGSZ = 80 /* Number of chars for args */

type KernelUid uint
type KernelGid uint

type KernelPid int

// Linux kernel data sizes: https://static.lwn.net/images/pdf/LDD3/ch11.pdf
// http://lxr.free-electrons.com/source/include/uapi/linux/elfcore.h#L78
type ElfPRPSInfo struct {
	State byte /* numeric process state */
	SName byte /* char for pr_state */
	Zomb byte  /* zombie */
	Nice byte  /* nice val */
	Flag uint  /* flags */

	UID KernelUid
	GID KernelGid
	PID, PPID, PGRP, SID KernelPid

	/* Lots missing */
	FName [16]byte /* filename of executable */
	PSArgs [ELF_PRARGSZ]byte /* initial part of arg list */
}
