# A few notes on MPK

MPK is a __userspace__ mechanism to control page table permissions. Permissions set on a per thread basis (permissions are stored inside a _protection key rights register_ `PKRU`) and apply to a group of page tables instead of a single page. 

A key can be associated to RW, R or no access, execution cannot be restricted. The final permission is the intersection of thread local key permission and process permissions on the page table. 

A note on performance: unpriviligied instruction (WRPKRU) that executes in less than 20 cycles, i.e. 10 to 100 times faster than a syscall, plus there is no TLB flush.

PKRU is a 32 bits register, each protection key is associated with two bits:
- Bit _2i_ block any data access if set to 1 (access disable bit)
- Bit _2i+1_ block any write if set to 1 (write disable bit)

When using RDPKRU and WRPKRU to read and write PKRU the register ECX must be set to 0.

Protections keys are 4 bits long (because there is 16 groups).

Three system calls directly interact with pkeys:
- int pkey_alloc(unsigned long flags, unsigned long ini_access_rights)
- int pkey_free(int pkey)
- int pkey_mprotect(unsigned long start, size_t len, unsigned long protection, int pkey)


# libmpk

libmpk addresses these three main issues:
- Vulenerability to protection-key-use-after-free: ???
- Limited number of protection key (16)
- Incompatible with `mprotec)` 

## protection-key-use-after-free:

My current understanting is that if a page is tagged with key 6 for instance, then key 6 is freed by `pkey_free` and reallocated later, then the page is still tagged with key 6 which may cause undesired access right restriction.

## Limited number of protection key

Only 16 keys per process, any attempt to alocate extra keys with `pkey_alloc` will fail.

## Incompatible with `mprotect`

`mprotect` changes protection for ALL THE THREADS of the calling process, whereas MPK update protections on a thread basis.

# To be investigated

Group 0 has a special purpose: Which one ???
Linux syscalls `pkey_alloc` and `pkey_free`
Thread or hyperthread? Is the register PKRU common to both thread of an hyperthread?
The PKRU register is XSAVE-managed and can thus be read or written by instruction in the XSAVE feature set - What does that mean???
