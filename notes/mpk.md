# A few notes on MPK

MPK is a __userspace__ mechanism to control page table permissions. Permissions set on a per thread basis (permissions are stored inside a _protection key rights register_ `PKRU`) and apply to a group of page tables instead of a single page. 

A key can be associated to RW, R or no access, execution cannot be restricted. The final permission is the intersection of thread local key permission and process permissions on the page table. 

A note on performance: unpriviligied instruction (WRPKRU) that executes in less than 20 cycles, i.e. 10 to 100 times faster than a syscall, plus there is no TLB flush.

# libmpk

libmpk addresses these three main issues:
- Vulenerability to protection-key-use-after-free: ???
- Limited number of protection key (16)
- Incompatible with `mprotect()` 

## protection-key-use-after-free:

My current understanting is that if a page is tagged with key 6 for instance, then key 6 is freed by `pkey_free` and reallocated later, then the page is still tagged with key 6 which may cause undesired access right restriction.

# To be investigated

Group 0 has a special purpose: Which one ???
Thread or hyperthead? Is the register PKRU shared by the two threads of the hyperthread?
Linux syscalls `pkey_alloc` and `pkey_free`
