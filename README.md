# goelf

Golang specific ELF reader/parser CLI tool

## Usage

    $ go build && ./goelf --all -f ./goelf
    
## Example output

    $ go build && ./goelf --all -f ./goelf                       
    
        CLASS    |    DATA     |  VERSION   |     OSABI     | ABIVERSION |  BYTEORDER   |  TYPE   |  MACHINE  |  ENTRY    
    +------------+-------------+------------+---------------+------------+--------------+---------+-----------+----------+
      ELFCLASS64 | ELFDATA2LSB | EV_CURRENT | ELFOSABI_NONE | 0x0        | LittleEndian | ET_EXEC | EM_X86_64 | 0x4568d0  

      ID |      SECTIONS      |      TYPE       |            FLAGS            |   ADDR   |  OFFSET  |  SIZE   | LINK | INFO | ADDRALIGN | ENTSIZE  
    +----+--------------------+-----------------+-----------------------------+----------+----------+---------+------+------+-----------+---------+
      0  |                    | SHT_NULL        | 0x0                         | 0x0      | 0x0      | 0x0     | 0x0  | 0x0  | 0         | 0        
      1  | .text              | SHT_PROGBITS    | SHF_ALLOC+SHF_EXECINSTR     | 0x401000 | 0x1000   | 0xe6d6e | 0x0  | 0x0  | 16        | 0        
      2  | .plt               | SHT_PROGBITS    | SHF_ALLOC+SHF_EXECINSTR     | 0x4e7d80 | 0xe7d80  | 0x180   | 0x0  | 0x0  | 16        | 16       
      3  | .rodata            | SHT_PROGBITS    | SHF_ALLOC                   | 0x4e8000 | 0xe8000  | 0x5a6aa | 0x0  | 0x0  | 32        | 0        
      4  | .typelink          | SHT_PROGBITS    | SHF_ALLOC                   | 0x5426ac | 0x1426ac | 0x13f8  | 0x0  | 0x0  | 4         | 0        
      5  | .itablink          | SHT_PROGBITS    | SHF_ALLOC                   | 0x543aa8 | 0x143aa8 | 0x1a0   | 0x0  | 0x0  | 8         | 0        
      6  | .gosymtab          | SHT_PROGBITS    | SHF_ALLOC                   | 0x543c48 | 0x143c48 | 0x0     | 0x0  | 0x0  | 1         | 0        
      7  | .gopclntab         | SHT_PROGBITS    | SHF_ALLOC                   | 0x543c60 | 0x143c60 | 0x72d3b | 0x0  | 0x0  | 32        | 0        
      8  | .rela              | SHT_RELA        | SHF_ALLOC                   | 0x5b69a0 | 0x1b69a0 | 0x18    | 0xf  | 0x0  | 8         | 24       
      9  | .rela.plt          | SHT_RELA        | SHF_ALLOC                   | 0x5b69c0 | 0x1b69c0 | 0x228   | 0xf  | 0x2  | 8         | 24       
      10 | .gnu.version       | SHT_GNU_VERSYM  | SHF_ALLOC                   | 0x5b6c00 | 0x1b6c00 | 0x38    | 0xf  | 0x0  | 2         | 2        
      11 | .gnu.version_r     | SHT_GNU_VERNEED | SHF_ALLOC                   | 0x5b6c40 | 0x1b6c40 | 0x70    | 0xe  | 0x2  | 8         | 0        
      12 | .hash              | SHT_HASH        | SHF_ALLOC                   | 0x5b6cc0 | 0x1b6cc0 | 0x90    | 0xf  | 0x0  | 8         | 4        
      13 | .shstrtab          | SHT_STRTAB      | 0x0                         | 0x0      | 0x1b6d60 | 0x16d   | 0x0  | 0x0  | 1         | 0        
      14 | .dynstr            | SHT_STRTAB      | SHF_ALLOC                   | 0x5b6ee0 | 0x1b6ee0 | 0x1d4   | 0x0  | 0x0  | 1         | 0        
      15 | .dynsym            | SHT_DYNSYM      | SHF_ALLOC                   | 0x5b70c0 | 0x1b70c0 | 0x2a0   | 0xe  | 0x0  | 8         | 24       
      16 | .got.plt           | SHT_PROGBITS    | SHF_WRITE+SHF_ALLOC         | 0x5b8000 | 0x1b8000 | 0xd0    | 0x0  | 0x0  | 8         | 8        
      17 | .dynamic           | SHT_DYNAMIC     | SHF_WRITE+SHF_ALLOC         | 0x5b80e0 | 0x1b80e0 | 0x130   | 0xe  | 0x0  | 8         | 16       
      18 | .got               | SHT_PROGBITS    | SHF_WRITE+SHF_ALLOC         | 0x5b8210 | 0x1b8210 | 0x8     | 0x0  | 0x0  | 8         | 8        
      19 | .noptrdata         | SHT_PROGBITS    | SHF_WRITE+SHF_ALLOC         | 0x5b8220 | 0x1b8220 | 0x114a0 | 0x0  | 0x0  | 32        | 0        
      20 | .data              | SHT_PROGBITS    | SHF_WRITE+SHF_ALLOC         | 0x5c96c0 | 0x1c96c0 | 0x8010  | 0x0  | 0x0  | 32        | 0        
      21 | .bss               | SHT_NOBITS      | SHF_WRITE+SHF_ALLOC         | 0x5d16e0 | 0x1d16e0 | 0x1ad90 | 0x0  | 0x0  | 32        | 0        
      22 | .noptrbss          | SHT_NOBITS      | SHF_WRITE+SHF_ALLOC         | 0x5ec480 | 0x1ec480 | 0x4f60  | 0x0  | 0x0  | 32        | 0        
      23 | .tbss              | SHT_NOBITS      | SHF_WRITE+SHF_ALLOC+SHF_TLS | 0x0      | 0x0      | 0x8     | 0x0  | 0x0  | 8         | 0        
      24 | .debug_abbrev      | SHT_PROGBITS    | 0x0                         | 0x5f2000 | 0x1d2000 | 0xff    | 0x0  | 0x0  | 1         | 0        
      25 | .debug_line        | SHT_PROGBITS    | 0x0                         | 0x5f20ff | 0x1d20ff | 0x2df7f | 0x0  | 0x0  | 1         | 0        
      26 | .debug_frame       | SHT_PROGBITS    | 0x0                         | 0x62007e | 0x20007e | 0x1b34c | 0x0  | 0x0  | 1         | 0        
      27 | .debug_pubnames    | SHT_PROGBITS    | 0x0                         | 0x63b3ca | 0x21b3ca | 0x208b1 | 0x0  | 0x0  | 1         | 0        
      28 | .debug_pubtypes    | SHT_PROGBITS    | 0x0                         | 0x65bc7b | 0x23bc7b | 0xfd0c  | 0x0  | 0x0  | 1         | 0        
      29 | .debug_aranges     | SHT_PROGBITS    | 0x0                         | 0x66b987 | 0x24b987 | 0x30    | 0x0  | 0x0  | 1         | 0        
      30 | .debug_gdb_scripts | SHT_PROGBITS    | 0x0                         | 0x66b9b7 | 0x24b9b7 | 0x35    | 0x0  | 0x0  | 1         | 0        
      31 | .debug_info        | SHT_PROGBITS    | 0x0                         | 0x66b9ec | 0x24b9ec | 0x74dbc | 0x0  | 0x0  | 1         | 0        
      32 | .interp            | SHT_PROGBITS    | SHF_ALLOC                   | 0x400fe4 | 0xfe4    | 0x1c    | 0x0  | 0x0  | 1         | 0        
      33 | .note.go.buildid   | SHT_NOTE        | SHF_ALLOC                   | 0x400fac | 0xfac    | 0x38    | 0x0  | 0x0  | 4         | 0        
      34 | .symtab            | SHT_SYMTAB      | 0x0                         | 0x0      | 0x2c1000 | 0x192c0 | 0x23 | 0x7a | 8         | 24       
      35 | .strtab            | SHT_STRTAB      | 0x0                         | 0x0      | 0x2da2c0 | 0x1c809 | 0x0  | 0x0  | 1         | 0        

           PROGS       |   FLAGS   |   OFF    |  VADDR   |  PADDR   | FILESZ  |  MEMSZ  | ALIGN   
    +------------------+-----------+----------+----------+----------+---------+---------+--------+
      PT_PHDR          | PF_R      | 0x40     | 0x400040 | 0x400040 | 0x230   | 0x230   | 0x1000  
      PT_INTERP        | PF_R      | 0xfe4    | 0x400fe4 | 0x400fe4 | 0x1c    | 0x1c    | 0x1     
      PT_NOTE          | PF_R      | 0xfac    | 0x400fac | 0x400fac | 0x38    | 0x38    | 0x4     
      PT_LOAD          | PF_X+PF_R | 0x0      | 0x400000 | 0x400000 | 0xe7f00 | 0xe7f00 | 0x1000  
      PT_LOAD          | PF_R      | 0xe8000  | 0x4e8000 | 0x4e8000 | 0xcf360 | 0xcf360 | 0x1000  
      PT_LOAD          | PF_W+PF_R | 0x1b8000 | 0x5b8000 | 0x5b8000 | 0x196e0 | 0x393e0 | 0x1000  
      PT_DYNAMIC       | PF_W+PF_R | 0x1b80e0 | 0x5b80e0 | 0x5b80e0 | 0x130   | 0x130   | 0x8     
      PT_TLS           | PF_R      | 0x0      | 0x0      | 0x0      | 0x0     | 0x8     | 0x8     
      PT_LOOS+74769745 | PF_W+PF_R | 0x0      | 0x0      | 0x0      | 0x0     | 0x0     | 0x8     
      PT_LOOS+84153728 | 0x2a00    | 0x0      | 0x0      | 0x0      | 0x0     | 0x0     | 0x8     

          IMPORTED SYMBOLS      |   VERSION   |     LIBRARY      
    +---------------------------+-------------+-----------------+
      __stack_chk_fail          | GLIBC_2.4   | libc.so.6        
      stderr                    | GLIBC_2.2.5 | libc.so.6        
      fwrite                    | GLIBC_2.2.5 | libc.so.6        
      __vfprintf_chk            | GLIBC_2.3.4 | libc.so.6        
      fputc                     | GLIBC_2.2.5 | libc.so.6        
      abort                     | GLIBC_2.2.5 | libc.so.6        
      pthread_create            | GLIBC_2.2.5 | libpthread.so.0  
      strerror                  | GLIBC_2.2.5 | libc.so.6        
      __fprintf_chk             | GLIBC_2.3.4 | libc.so.6        
      pthread_mutex_lock        | GLIBC_2.2.5 | libpthread.so.0  
      pthread_cond_wait         | GLIBC_2.3.2 | libpthread.so.0  
      pthread_mutex_unlock      | GLIBC_2.2.5 | libpthread.so.0  
      pthread_cond_broadcast    | GLIBC_2.3.2 | libpthread.so.0  
      free                      | GLIBC_2.2.5 | libc.so.6        
      malloc                    | GLIBC_2.2.5 | libc.so.6        
      pthread_attr_init         | GLIBC_2.2.5 | libpthread.so.0  
      pthread_attr_getstacksize | GLIBC_2.2.5 | libpthread.so.0  
      pthread_attr_destroy      | GLIBC_2.2.5 | libpthread.so.0  
      __errno_location          | GLIBC_2.2.5 | libpthread.so.0  
      sigfillset                | GLIBC_2.2.5 | libc.so.6        
      pthread_sigmask           | GLIBC_2.2.5 | libpthread.so.0  
      mmap                      | GLIBC_2.2.5 | libc.so.6        
      setenv                    | GLIBC_2.2.5 | libc.so.6        
      unsetenv                  | GLIBC_2.2.5 | libc.so.6        

          LIBRARY      
    +-----------------+
      libpthread.so.0  
      libc.so.6        

                                                                       SYM                                                                   | INFO | OTHER |  SECTION  |  OFFSET  |  SIZE   
    +----------------------------------------------------------------------------------------------------------------------------------------+------+-------+-----------+----------+--------+
      go.go                                                                                                                                  | 0x4  | 0x0   | SHN_ABS   | 0x0      | 0       
      runtime.text                                                                                                                           | 0x2  | 0x0   | 1         | 0x401000 | 0       
      runtime.etext                                                                                                                          | 0x2  | 0x0   | 1         | 0x4e7d6e | 0       
      runtime.end                                                                                                                            | 0x1  | 0x0   | 22        | 0x5f13e0 | 0       
      $f64.8000000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa98 | 8       
      $f64.3eb0000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa40 | 8       
      $f32.40d00000                                                                                                                          | 0x1  | 0x0   | 3         | 0x53fa2c | 4       
      runtime.data                                                                                                                           | 0x1  | 0x0   | 20        | 0x5c96c0 | 0       
      $f64.403a000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa78 | 8       
      $f64.bfe62e42fefa39ef                                                                                                                  | 0x1  | 0x0   | 3         | 0x53faa0 | 8       
      $f64.4059000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa80 | 8       
      $f64.3ff0000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa60 | 8       
      $f64.43e0000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa90 | 8       
      $f64.3fd0000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa48 | 8       
      $f64.3fe0000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa50 | 8       
      $f64.3fee666666666666                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa58 | 8       
      $f64.4024000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa70 | 8       
      $f64.4014000000000000                                                                                                                  | 0x1  | 0x0   | 3         | 0x53fa68 | 8       
      runtime.memhash_varlen.args_stackmap                                                                                                   | 0x1  | 0x0   | 3         | 0x5402a0 | 16      
      runtime.reflectcall.args_stackmap                                                                                                      | 0x1  | 0x0   | 3         | 0x540198 | 12      
      runtime.cgocallback_gofunc.args_stackmap                                                                                               | 0x1  | 0x0   | 3         | 0x540188 | 12      
      runtime.publicationBarrier.args_stackmap                                                                                               | 0x1  | 0x0   | 3         | 0x53fe98 | 8       
      runtime.asmcgocall.args_stackmap                                                                                                       | 0x1  | 0x0   | 3         | 0x540280 | 16      
      runtime.call32.args_stackmap                                                                                                           | 0x1  | 0x0   | 3         | 0x5400b8 | 12      
      runtime.call64.args_stackmap                                                                                                           | 0x1  | 0x0   | 3         | 0x540138 | 12      
      ...
