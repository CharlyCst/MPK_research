.globl  main

.text 
main:
        movq    $len, %rdx # msg len
        movq    $msg, %rcx # msg pointer
        movq    $1, %rbx   # fd stdout
        movq    $4, %rax   # syscall number 
        int     $0x80      #  syscall

	pushq   $65
	call    print

        movl    $0, %ebx
        movl    $1, %eax
        int     $0x80

print:
	movq    $4, %rdx
	movq    %rbp, %rcx
	addq    $1, %rcx
	movq    $1, %rbx
	movq    $4, %rax
	int     $0x80
	ret
.data
msg:
        .ascii  "Hello, world!\n"
        len =   . - msg
