# Notes on dune:

When compiling the kernel with a recent version of gcc and error _code model kernel does not support PIC mode_ may occur. To fix it update the kernel Makefile:

```makefile
# Old version:
KBUILD_CFLAGS += $(CFLAGS_EXTRA)

# Change for:
EXTRA_CFLAGS += $(CFLAGS_EXTRA) -fno-pie
```

If some includes under _/asm_ are missing, cd into kernel includes and create a symlink from asm to asm-generic:

```bash
ln -s asm-generic asm
```

Other imports may be missing, for instance _stdarg.h_ from gcc, just add the corresponding symlinks:

```bash
sudo ln -s /usr/lib/gcc/x86_64-linux-gnu/7.5.0/include/stdarg.h stdarg.h
```
