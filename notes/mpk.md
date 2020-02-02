# A few notes on MPK

MPK is a __userspace__ mechanism to control page table permissions. Permissions set on a per thread basis (permissions are stored inside a _protection key rights register_ `PKRU`) and apply to a group of page tables associated to a key instead of a single page. 

A key can be associated to RW, R or no access, execution cannot be restricted. The final permission is the intersection of thread local key permission and process permissions on the page table. 
Protections keys are 4 bits long (because there is 16 groups).

A note on performance: unpriviligied instruction (WRPKRU) that executes in less than 20 cycles, i.e. 10 to 100 times faster than a syscall, plus there is no TLB flush.

PKRU is a 32 bits register, each protection key is associated with two bits:
- Bit _2i_ block any data access if set to 1 (access disable bit, or AD)
- Bit _2i+1_ block any write if set to 1 (write disable bit, or WD)

Two ASM instructions are available to interact with PKRU:
- WRPKRU: overwrite PKRU with EAX, the two registers ECX and EDX must be filled with 0.
- RDPKRU: return the current value of PKRU inside EAX, the ECX register must be filled with 0 while EDX will be overwritted with 0.
The actual usage of ECX and EDX is undocumented.

Three system calls directly interact with pkeys:
```C
int pkey_alloc(unsigned long flags, unsigned long ini_access_rights);
int pkey_free(int pkey);
int pkey_mprotect(unsigned long start, size_t len, unsigned long protection, int pkey);
```

Because kernel mode is required to modify page table entries (PTE), the system call `pkey_mprotect` is required to tag memory pages with a key group. `pkey_mprotect` is an extension of `mprotect`: it allows to update protection right of page while setting the key group.

# libmpk

The aim of libmpk is to prevent an adversary from reading from or writing in sensitive pages through memory corruption vulnerabilities. 

A program relying on libmpk should only use MPK through libmpk.

## Limitation of MPK

libmpk addresses these three main issues:
- Vulenerability to protection-key-use-after-free: ???
- Limited number of protection key (16)
- Incompatible with `mprotec)` 

### protection-key-use-after-free:

My current understanting is that if a page is tagged with key 6 for instance, then key 6 is freed by `pkey_free` and reallocated later, then the page is still tagged with key 6 which may cause undesired access right restriction.

### Limited number of protection key

Only 16 keys per process, any attempt to alocate extra keys with `pkey_alloc` will fail.

### Incompatible with `mprotect`

`mprotect` changes protection for ALL THE THREADS of the calling process, whereas MPK update protections on a thread basis.

## Proposed solution

### API

- `mpk_init(evict_rate)`
- `mpk_mmap(vkey, addr, len, prot, flags, fd, offset)`
- `mpk_munmap(vkey)`
- `mpk_begin(vkey, prot)`
- `mpk_end(vkey)`
- `mpk_mprotect(vkey, prot)`
- `mpk_malloc(vkey, size)`
- `mpk_free(size)`

`mpk_init` must be called at the beginning of the program, it allocates all the virtual key from the kernel.

Virtual keys are integers chosen by the developer.

libmpk maintains the mappings between virtual keys and pages to avoid scanning through pages with `mpk_munmap`.

libmpk offer a proper execute only permission (synchronized among threads, contrary to `mprotect` for execute only) and better performances overall.

### Key virtualization

libmpk enables an application to use more than 16 page groups by virtualizing hardware keys. The application is prohibited from manipulating hardware keys.

Virtual keys are mapped either to a physical key or to nothing (null). When the application call libmpk with a virtual key, if it is mapped to a physical key permissions can be updated through MPK, otherwise either libmpk evict another key (and update the mappings with `pkey_mprotect`) or does not update the cache and simply call `mprotect`. The frequency of eviction is determined by the eviction rate.

libmpk always use a physical key for thread logal permissions (`mpk_begin`), thus the call to `mpk_begin` may rise an error if no physical keys are available (all used by `mpk_begin`).

libmpk reserves one key for execute only memory (which the `mprotect` also achieve through MPK), and use it to map all execute only memory. This key is not freed until at lest one page is execute only.

### Key synchronisation

Because MPK is thread local, in order to emulate `mprotect` it is necessary to synchronise the value of PKRU among threads.

To increase performances PKRU is updates lazily: if a thread is not currently scheduled it does not need the up to date PKRU. Synchronization is forced for thread in userspace while a hook is created for sleeping and kernel mode thread, it will be invoked right before jumping to user space again.

In linux this can be achieved by sending a rescheduling interupt and adding a callback (with `task_work_add()`). 

### Security

libmpk maps its metadata into two virtual pages: one with read right for user code and one with write access for libmpk kernel code.

# To be investigated

Group 0 has a special purpose: Which one ???
Thread or hyperthread? Is the register PKRU common to both thread of an hyperthread?
The PKRU register is XSAVE-managed and can thus be read or written by instruction in the XSAVE feature set - What does that mean???
Is there a reason why libmpk insist on the fact that virtual keys should be constant?
How to ensure that the application can not use WRPKRU and RDPKRU?
How to run part of libmpk into kernel mode?
How to create hooks?
