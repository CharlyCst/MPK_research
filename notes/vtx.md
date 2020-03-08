# VT-x

## VT-x extension

Virtualization extension, AMD provide a similar HW interface called SVM.

With VT-x the CPU has two modes:

- VMX root, which does not change CPU behavior and enable two new instruction for managing VT-x.
- VMX non-root, restrict CPU behavior and is intended for running virtualized guest OSes.

The `VMLAUNCH` and `VMRESUME` instructions perform a VM entry, placing the CPU in VMX non-root mode. Then when the VMM is needed a VM exit is performed, followed by a jump to the corresponding VMM entry point.

The hardware automatically saves and restores most architectural state during both transitions, using an in-memory data structure called the VM control structure, or _VMCS_. The VMCS also contains a bunch of configuration, such as the instructions that should trigger a VM exit.

Some hardware interfaces, such as the interrupt descriptor table (IDT) and privilege modes are still exposed in VMX non-root mode and does not generate VM exits.

It is also possible to manually require a VM exist with the `VMCALL` instruction.

It is not possible to expose the page table root register _%CR3_ to the guest without trusting it, thus VT-x introduces a hardware mechanism called _extended page table_, or _EPT_, which adds another level of address translation. AMD's equivalent mechanism is called _nested page table_, or _NPT_.

## Dune:

- Small kernel module to initialize virtualization and manage syscalls
- User level library to help managing hardware primitives

Process then use the `VMCALL` instruction to invoke syscalls

Dune exposes three hardware features:

- Exceptions
- Virtual memory
- Privilege modes
