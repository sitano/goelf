package elf

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
type ElfPRStatus struct {
									   //#if 0
									   //        long    pr_flags /* XXX Process flags */
									   //        short   pr_why   /* XXX Reason for process halt */
									   //        short   pr_what  /* XXX More detailed reason */
									   //#endif

	Info ElfSigInfo                    /* Info associated with signal */
	CurSig int16                       /* Current signal */
	SigPend uint                       /* Set of pending signals */
	SigHold uint                       /* Set of held signals */

									   //#if 0
									   //        struct sigaltstack pr_altstack; /* Alternate stack info */
									   //        struct sigaction pr_action;     /* Signal action for current sig */
									   //#endif

	PID                     KernelPid
	PPID                    KernelPid
	PGRP                    KernelPid
	SID                     KernelPid

	PR_UTime                TimeVal    /* User time */
	PR_STime                TimeVal    /* System time */
	PR_CUTime               TimeVal    /* Cumulative user time */
	PR_CSTime               TimeVal    /* Cumulative system time */

									   //#if 0
									   //        long    pr_instr;               /* Current instruction */
									   //#endif

    pr_reg                  ElfGRegSet /* GP registers */

									   // #ifdef CONFIG_BINFMT_ELF_FDPIC
									   /* When using FDPIC, the loadmap addresses need to be communicated
										* to GDB in order for GDB to do the necessary relocations.  The
										* fields (below) used to communicate this information are placed
										* immediately after ``pr_reg'', so that the loadmap addresses may
										* be viewed as part of the register set if so desired.
										*/
	PR_Exec_FDPic_LoadMap   uint
	PR_Interp_FDPic_LoadMap uint
									   // #endif

    PRFPValid               int32      /* True if math co-processor being used.  */
}
