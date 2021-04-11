## ELF文件结构


真正了不起的程序员对自己的程序的每一个字节都了如指掌。

现在PC平台流行的可执行文件格式（Executable）主要是Windows的PE（Portable Executable）和Linux的ELF（Executable Linkable Format），它们都是COFF（Common file format）格式的变种。

ELF文件标准里面把系统中采用ELF格式的文件归为4类。

| ELF文件类型 | 说明 | 实例 |
| ---------- | ---- | --- |
| 可重定位文件（Relocalable File） | 这类文件包含代码和数据，可以被用来链接成可执行文件或共享目标文件，静态文件也可以归为该类 | linux的`.o` windows的`.obj` |
| 可执行文件（Executable File） | 直接可以运行的程序，linux下一般没有扩展名 | 例如`/bin/bash`，windows的`.exe` |
| 共享目标文件（Shared Object File） | 这类文件使用的两种情况：跟其他的可重定位文件和共享文件链接，产生新的目标文件；或由动态链接器加载，作为进程映像的一部分来运行 | linux的`.so`、windows的`.dll` |
| 核心转储文件（Core Dump File） | 当进程意外终止时的核心转储文件 | linux下的core dump |

linux下可通过file命令查看文件格式，例如

```shell
$file /bin/bash
/bin/bash: ELF 64-bit LSB shared object, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/l, for GNU/Linux 3.2.0, BuildID[sha1]=12f73d7a8e226c663034529c8dd20efec22dde54, stripped
```

## C代码文件与目标文件

这个`SimpleSection.c`是一个具有代表性的C代码文件，通过`gcc -c SimpleSection.c`可以生成对应的目标文件。

```c
int printf( const char* format, ... );
int global_init_var = 84;
int global_uninit_var;

void func1( int i ) 
{
    printf( "%d\n",  i );
}

int main(void)
{
    static int static_var = 85;
    static int static_var2;
    int a = 1;
    int b;
    func1( static_var + static_var2 + a + b );
    return a;
}
```

一般C代码文件编译后
- 执行语句都会被编译成机器代码保存在`.text`段；
- 已初始化的全局变量和局部静态变量都保存在`.data`段；
- 未初始化或初始化为0的全局变量和局部静态变量一般放在`.bss`段。
    - 程序运行的时候，未初始化的全局变量和局部静态变量是要占内存空间的，并且可执行文件必须记录它们的大小总和，记为`.bss`段。
    - `.bss`段只是为未初始化的全局变量和局部静态变量预留位置而已，它并没有内容，所以它在文件中也不占据空间。

![c_code_elf](./.assets/c_code_elf.png)

除此之外，还可以看到一个“文件头”，它描述了整个文件的文件属性，包括文件是否可执行、是静态链接还是动态链接及入口地址（如果是可执行文件）、目标硬件、目标操作系统等信息，文件头还包括一个段表（Section Table）。段表其实是一个描述文件中各个段的数组，描述了文件中各个段在文件中的偏移位置及段的属性等。

**为什么要区分代码段和数据段？**

- **权限**：当程序被装载后，数据和指令分别被映射到两个虚存区域。对于程序而言，数据段是可读写的，而代码段是只读的，据此可以设置两个虚存区域的权限；
- **局部性原理**：现代的CPU来说，它们有着极为强大的缓存（Cache）体系，且一般都被设计成数据缓存和指令缓存分离。指令区和数据区的分离有利于提高程序的局部性。
- **共享**：有利于指令的共享。当系统中运行着多个该程序的副本时，它们的指令都是一样的，所以内存中只须要保存一份该程序的指令部分。

## 与程序运行密切相关的段结构

在Linux下，我们可以使用binutils的工具objdump来查看object内部的结构，

```shell
$objdump -h SimpleSection.o

SimpleSection.o：     文件格式 elf64-x86-64

节：
Idx Name          Size      VMA               LMA               File off  Algn
  0 .text         00000057  0000000000000000  0000000000000000  00000040  2**0
                  CONTENTS, ALLOC, LOAD, RELOC, READONLY, CODE
  1 .data         00000008  0000000000000000  0000000000000000  00000098  2**2
                  CONTENTS, ALLOC, LOAD, DATA
  2 .bss          00000004  0000000000000000  0000000000000000  000000a0  2**2
                  ALLOC
  3 .rodata       00000004  0000000000000000  0000000000000000  000000a0  2**0
                  CONTENTS, ALLOC, LOAD, READONLY, DATA
  4 .comment      0000002c  0000000000000000  0000000000000000  000000a4  2**0
                  CONTENTS, READONLY
  5 .note.GNU-stack 00000000  0000000000000000  0000000000000000  000000d0  2**0
                  CONTENTS, READONLY
  6 .eh_frame     00000058  0000000000000000  0000000000000000  000000d0  2**3
                  CONTENTS, ALLOC, LOAD, RELOC, READONLY, DATA
```

先不管其他段，只研究`.text`、`.data`、`.bss`段，每个段的第2行内容为段的属性。“CONTENTS”表示该段在文件中存在。`.bss`段没有该属性，表示它在ELF文件中不存在内容。根据段的长度（Size）和段所在的位置（File Offset），我们可以画出ELF文件的大致结构

![elf_struct](./.assets/elf_struct.jpg)

size命令可以用来查看ELF文件的代码段、数据段和BSS段的长度（**注**：size默认是运行在“Berkeley compatibility mode”下。在这种模式下，会将不可执行的拥有“ALLOC”属性的只读段归到`.text`段下，在这里就是`.rodata`段和`.eh_frame`段。如果你使用`-A`选项，那么size会运行在“System V compatibility mode”，此时就跟`objdump -h`显示的`.text`段的大小差不多了）。

```shell
$size -A SimpleSection.o
SimpleSection.o  :
section           size   addr
.text               87      0
.data                8      0
.bss                 4      0
.rodata              4      0
.comment            44      0
.note.GNU-stack      0      0
.eh_frame           88      0
Total              235
```

通过`objdump -s -d SimpleSection.o`打印所有段的信息

```shell
SimpleSection.o：     文件格式 elf64-x86-64

Contents of section .text:
 0000 554889e5 4883ec10 897dfc8b 45fc89c6  UH..H....}..E...
 0010 488d3d00 000000b8 00000000 e8000000  H.=.............
 0020 0090c9c3 554889e5 4883ec10 c745f801  ....UH..H....E..
 0030 0000008b 15000000 008b0500 00000001  ................
 0040 c28b45f8 01c28b45 fc01d089 c7e80000  ..E....E........
 0050 00008b45 f8c9c3                      ...E...         
Contents of section .data:
 0000 54000000 55000000                    T...U...        
Contents of section .rodata:
 0000 25640a00                             %d..            
Contents of section .comment:
 0000 00474343 3a202855 62756e74 7520372e  .GCC: (Ubuntu 7.
 0010 342e302d 31756275 6e747531 7e31382e  4.0-1ubuntu1~18.
 0020 30342e31 2920372e 342e3000           04.1) 7.4.0.    
Contents of section .eh_frame:
 0000 14000000 00000000 017a5200 01781001  .........zR..x..
 0010 1b0c0708 90010000 1c000000 1c000000  ................
 0020 00000000 24000000 00410e10 8602430d  ....$....A....C.
 0030 065f0c07 08000000 1c000000 3c000000  ._..........<...
 0040 00000000 33000000 00410e10 8602430d  ....3....A....C.
 0050 066e0c07 08000000                    .n......        

Disassembly of section .text:

0000000000000000 <func1>:
   0:   55                      push   %rbp
   1:   48 89 e5                mov    %rsp,%rbp
   4:   48 83 ec 10             sub    $0x10,%rsp
   8:   89 7d fc                mov    %edi,-0x4(%rbp)
   b:   8b 45 fc                mov    -0x4(%rbp),%eax
   e:   89 c6                   mov    %eax,%esi
  10:   48 8d 3d 00 00 00 00    lea    0x0(%rip),%rdi        # 17 <func1+0x17>
  17:   b8 00 00 00 00          mov    $0x0,%eax
  1c:   e8 00 00 00 00          callq  21 <func1+0x21>
  21:   90                      nop
  22:   c9                      leaveq 
  23:   c3                      retq   

0000000000000024 <main>:
  24:   55                      push   %rbp
  25:   48 89 e5                mov    %rsp,%rbp
  28:   48 83 ec 10             sub    $0x10,%rsp
  2c:   c7 45 f8 01 00 00 00    movl   $0x1,-0x8(%rbp)
  33:   8b 15 00 00 00 00       mov    0x0(%rip),%edx        # 39 <main+0x15>
  39:   8b 05 00 00 00 00       mov    0x0(%rip),%eax        # 3f <main+0x1b>
  3f:   01 c2                   add    %eax,%edx
  41:   8b 45 f8                mov    -0x8(%rbp),%eax
  44:   01 c2                   add    %eax,%edx
  46:   8b 45 fc                mov    -0x4(%rbp),%eax
  49:   01 d0                   add    %edx,%eax
  4b:   89 c7                   mov    %eax,%edi
  4d:   e8 00 00 00 00          callq  52 <main+0x2e>
  52:   8b 45 f8                mov    -0x8(%rbp),%eax
  55:   c9                      leaveq 
  56:   c3                      retq 
```

### 代码段

`Contents of section .text:`就是.text的数据以十六进制方式打印出来的内容，共0x57字节。
`Disassembly of section .text:`就是.text的数据的反编译结果。

### 数据段和只读数据段

.data段保存的是那些已经初始化了的全局静态变量和局部静态变量。SimpleSection.c代码里面一共有两个这样的变量，分别是`global_init_varabal`与`static_var`。这两个变量每个4个字节，一共刚好8个字节。

我们在调用“printf”的时候，用到了一个字符串常量“%d\n”，它是一种只读数据，所以它被放到了`.rodata`段。可以看到`.rodata`段的4个字节刚好是这个字符串常量的ASCII码，最后以`\0`结尾。

### BSS段

`global_uninit_var`和`static_var2`就是被存放在`.bss`段，其实更准确的说法是.bss段为它们预留了空间。但是我们可以看到该段的大小只有4个字节，这与`global_uninit_var`和`static_var2`的大小的8个字节不符。

通过符号表（Symbol Table）看到，只有`static_var2`被存放在了.bss段，而`global_uninit_var`却没有被存放在任何段，只是一个未定义的“COMMON符号”（这其实是跟不同的语言与不同的编译器实现有关，有些编译器会将全局的未初始化变量存放在目标文件.bss段，有些则不存放，只是预留一个未定义的全局变量符号，等到最终链接成可执行文件的时候再在.bss段分配空间）。

```shell
$objdump -x -s -d SimpleSection.o

...
SYMBOL TABLE:
0000000000000000 l    df *ABS*  0000000000000000 SimpleSection.c
0000000000000000 l    d  .text  0000000000000000 .text
0000000000000000 l    d  .data  0000000000000000 .data
0000000000000000 l    d  .bss   0000000000000000 .bss
0000000000000000 l    d  .rodata        0000000000000000 .rodata
0000000000000004 l     O .data  0000000000000004 static_var.1802
0000000000000000 l     O .bss   0000000000000004 static_var2.1803
0000000000000000 l    d  .note.GNU-stack        0000000000000000 .note.GNU-stack
0000000000000000 l    d  .eh_frame      0000000000000000 .eh_frame
0000000000000000 l    d  .comment       0000000000000000 .comment
0000000000000000 g     O .data  0000000000000004 global_init_var
0000000000000004       O *COM*  0000000000000004 global_uninit_var
0000000000000000 g     F .text  0000000000000024 func1
0000000000000000         *UND*  0000000000000000 _GLOBAL_OFFSET_TABLE_
0000000000000000         *UND*  0000000000000000 printf
0000000000000024 g     F .text  0000000000000033 main
...
```

### 其他段

| 段名 | 说明 |
| --- | --- |
| .rodata1 | Read only Data，只读数据，与..rodata一样 |
| .comment | 编译器版本信息，如“.GCC: (Ubuntu 7.4.0-1ubuntu1~18.04.1) 7.4.0.” |
| .eh_frame | 是一个记录序列，每个记录可以是CIE（Common Information Entry，公共信息条目）或FDE（Frame Description Entry，帧描述条目） |
| .debug | 调试信息 |
| .dynamic | 动态链接信息 |
| .hash | 符号哈希表 |
| .line | 调试时的行号表，即源代码行号与编译后指令的对应表 |
| .note | 额外的编译器信息，如程序的公司名、发布版本号等 |
| .strtab | String Table 字符串表，用于存放ELF文件中用到的各种字符串 |
| .symtab | Symbol Table 符号表 |
| .shstrtab | Section String Table 段名表 |
| .plt .got | 动态链接的跳转表和全局入口表 |
| .init .fini | 程序初始化与终结代码段 |

应用程序也可以**自定义段名**，但自定义的段名不能使用“.”作为前缀，否则容易跟系统保留段名冲突。另外，GCC提供了一个扩展机制，我们在全局变量或函数之前加上`__ attribute__((section("name")))`属性就可以把相应的变量或函数放到以“name”作为段名的段中。例如

```c
__attribute__((section("FOO"))) int global = 42;
```

---

## 其他重要部分

ELF文件结构中其他重要的部分，比如ELF文件头、段表、重定位表和字符串表等。

![elf_struct](./.assets/elf_struct_1.png)

### ELF文件头

ELF文件头（ELF Header）包含了描述整个文件的基本属性，比如ELF文件版本、目标机器型号、程序入口地址等。可以用readelf命令来详细查看ELF文件。

```shell
$readelf -h SimpleSection.o
ELF Header:
  Magic:   7f 45 4c 46 02 01 01 00 00 00 00 00 00 00 00 00 
  Class:                             ELF64
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              REL (Relocatable file)
  Machine:                           Advanced Micro Devices X86-64
  Version:                           0x1
  Entry point address:               0x0
  Start of program headers:          0 (bytes into file)
  Start of section headers:          1104 (bytes into file)
  Flags:                             0x0
  Size of this header:               64 (bytes)
  Size of program headers:           0 (bytes)
  Number of program headers:         0
  Size of section headers:           64 (bytes)
  Number of section headers:         13
  Section header string table index: 12
```

从上面输出的结果可以看到，ELF的文件头中定义了ELF魔数、文件机器字节长度、数据存储方式、版本、运行平台、ABI版本、ELF重定位类型、硬件平台、硬件平台版本、入口地址、程序头入口和长度、段表的位置和长度及段的数量等。

ELF文件头结构及相关常数被定义在`/usr/include/elf.h`里，因为ELF文件在各种平台下都通用。ELF文件有32位版本和64位版本，所以它的文件头结构也有`Elf32_Ehdr`和`Elf64_Ehdr`两种版本。它们的成员种类是一样的，只不过有些成员的大小不一样。为了对每个成员的大小做出明确的规定以便于在不同的编译环境下都拥有相同的字段长度，“elf.h”定义了一套自己的变量体系。

```c
/* Type for a 16-bit quantity.  */
typedef uint16_t Elf32_Half;
typedef uint16_t Elf64_Half;

/* Types for signed and unsigned 32-bit quantities.  */
typedef uint32_t Elf32_Word;
typedef	int32_t  Elf32_Sword;
typedef uint32_t Elf64_Word;
typedef	int32_t  Elf64_Sword;

/* Types for signed and unsigned 64-bit quantities.  */
typedef uint64_t Elf32_Xword;
typedef	int64_t  Elf32_Sxword;
typedef uint64_t Elf64_Xword;
typedef	int64_t  Elf64_Sxword;

/* Type of addresses.  */
typedef uint32_t Elf32_Addr;
typedef uint64_t Elf64_Addr;

/* Type of file offsets.  */
typedef uint32_t Elf32_Off;
typedef uint64_t Elf64_Off;
```

以64位版本`Elf64_Ehdr`为例，其结构为

```c
#define EI_NIDENT (16)

typedef struct
{
  unsigned char	e_ident[EI_NIDENT];	/* 魔数与其他内容 */
  Elf64_Half	e_type;			/* 文件类型 */
  Elf64_Half	e_machine;		/* 目标CPU体系结构 */
  Elf64_Word	e_version;		/* 文件版本 */
  Elf64_Addr	e_entry;		/* 程序入口地址 */
  Elf64_Off	e_phoff;		/* 程序头的偏移量 */
  Elf64_Off	e_shoff;		/* 段表在文件中的偏移量 */
  Elf64_Word	e_flags;		/* 特定于处理器的标志 */
  Elf64_Half	e_ehsize;		/* 文件头大小 */
  Elf64_Half	e_phentsize;		/* 程序头中每个结构的大小 */
  Elf64_Half	e_phnum;		/* 程序头中有多少个结构 */
  Elf64_Half	e_shentsize;		/* 段表中每个结构的大小 */
  Elf64_Half	e_shnum;		/* 段表中有多少个结构，也就是有多少段 */
  Elf64_Half	e_shstrndx;		/* 段表字符串表在段表中的下标 */
} Elf64_Ehdr;
```

从这里可以看出ELF文件头结构跟前面readelf输出的ELF文件头信息一一对应。其中e_ident的内容，对应`Magic`、`Class`、`Data`、`Version`、`OS/ABI`、`ABI Version`。

![magic_num](./.assets/magic_num.png)

前4个字节为文件魔数，对应`DEL`、`E`、`L`、`F`的ASCII码。很多文件的魔数由来都有它的历史背景，就如同“马屁股决定航天飞机”的故事。

这里不详细介绍所有的结构体成员。`e_type`指明了文件类型，就如之前介绍的可重定位文件、可执行文件、共享目标文件等。
```c
/* Legal values for e_type (object file type).  */
...
#define ET_REL		1		/* 可重定位文件 */
#define ET_EXEC	2		/* 可执行文件 */
#define ET_DYN		3		/* 共享目标文件 */
#define ET_CORE	4		/* Core文件 */
...
```

### 段表

**段表（Section Header Table）**就是保存这些段的基本属性的结构。段表是ELF文件中除了文件头以外最重要的结构，它描述了ELF的各个段的信息，比如每个段的段名、段的长度、在文件中的偏移、读写权限及段的其他属性。也就是说，ELF文件的段结构就是由段表决定的，编译器、链接器和装载器都是依靠段表来定位和访问各个段的属性的。段表在ELF文件中的位置由ELF文件头的`e_shoff`成员决定，即1104（0x450）。

`objdump -h`命令只能显示了ELF文件中关键的段，使用`readelf`工具来查看ELF文件真正的段结构。

```shell
$readelf -S SimpleSection.o
There are 13 section headers, starting at offset 0x450:

Section Headers:
  [Nr] Name              Type             Address           Offset
       Size              EntSize          Flags  Link  Info  Align
  [ 0]                   NULL             0000000000000000  00000000
       0000000000000000  0000000000000000           0     0     0
  [ 1] .text             PROGBITS         0000000000000000  00000040
       0000000000000057  0000000000000000  AX       0     0     1
  [ 2] .rela.text        RELA             0000000000000000  00000340
       0000000000000078  0000000000000018   I      10     1     8
  [ 3] .data             PROGBITS         0000000000000000  00000098
       0000000000000008  0000000000000000  WA       0     0     4
  [ 4] .bss              NOBITS           0000000000000000  000000a0
       0000000000000004  0000000000000000  WA       0     0     4
  [ 5] .rodata           PROGBITS         0000000000000000  000000a0
       0000000000000004  0000000000000000   A       0     0     1
  [ 6] .comment          PROGBITS         0000000000000000  000000a4
       000000000000002c  0000000000000001  MS       0     0     1
  [ 7] .note.GNU-stack   PROGBITS         0000000000000000  000000d0
       0000000000000000  0000000000000000           0     0     1
  [ 8] .eh_frame         PROGBITS         0000000000000000  000000d0
       0000000000000058  0000000000000000   A       0     0     8
  [ 9] .rela.eh_frame    RELA             0000000000000000  000003b8
       0000000000000030  0000000000000018   I      10     8     8
  [10] .symtab           SYMTAB           0000000000000000  00000128
       0000000000000198  0000000000000018          11    11     8
  [11] .strtab           STRTAB           0000000000000000  000002c0
       000000000000007c  0000000000000000           0     0     1
  [12] .shstrtab         STRTAB           0000000000000000  000003e8
       0000000000000061  0000000000000000           0     0     1
Key to Flags:
  W (write), A (alloc), X (execute), M (merge), S (strings), I (info),
  L (link order), O (extra OS processing required), G (group), T (TLS),
  C (compressed), x (unknown), o (OS specific), E (exclude),
  l (large), p (processor specific)
```

对于SimpleSection.o来说，段表是有13个元素的数组。ELF段表的这个数组的第1个元素是无效的段描述符，它的类型为`NULL`，除此之外每个段描述符都对应一个段，也就是说共有12个有效的段。其中，段描述符的结构为

```c
typedef struct
{
  Elf64_Word	sh_name;		/* 段名 (段名字符串在.shstrtab中的偏移) */
  Elf64_Word	sh_type;		/* 段的类型 */
  Elf64_Xword	sh_flags;		/* 段的标志位 */
  Elf64_Addr	sh_addr;		/* 段虚拟地址 */
  Elf64_Off	sh_offset;	/* 该段在文件中的偏移 */
  Elf64_Xword	sh_size;		/* 段的长度 */
  Elf64_Word	sh_link;		/* 段的链接信息 */
  Elf64_Word	sh_info;		/* 段的链接信息 */
  Elf64_Xword	sh_addralign;	/* 段地址对齐 */
  Elf64_Xword	sh_entsize;	/* 项的长度，一些段中会包含固定长度的项 */
} Elf64_Shdr;
```

段的名字只是在链接和编译过程中有意义，但它不能真正地表示段的类型。我们也可以将一个数据段命名为“.text”，对于编译器和链接器来说，主要决定段的属性的是段的类型（sh_type）和段的标志位（sh_flags）。

- 段的类型（sh_type）：

     ```c
     #define SHT_NULL	  0		/* Section header table entry unused */
     #define SHT_PROGBITS  1		/* Program data */
     #define SHT_SYMTAB	  2		/* Symbol table */
     #define SHT_STRTAB	  3		/* String table */
     #define SHT_RELA	  4		/* Relocation entries with addends */
     #define SHT_HASH	  5		/* Symbol hash table */
     #define SHT_DYNAMIC	  6		/* Dynamic linking information */
     #define SHT_NOTE	  7		/* Notes */
     ...
     ```

- 段的标志位（sh_flags）：表示该段在进程虚拟地址空间中的属性，比如是否可写，是否可执行等。

     ```c
     #define SHF_WRITE	     (1 << 0)	/* Writable */
     #define SHF_ALLOC	     (1 << 1)	/* Occupies memory during execution */
     #define SHF_EXECINSTR	(1 << 2)	/* Executable */
     #define SHF_MERGE	     (1 << 4)	/* Might be merged */
     ```

- 段的链接信息（sh_link 、 sh_info）：如果段的类型是与链接相关的（不论是动态链接或静态链接），比如重定位表、符号表等；对于其他类型的段，这两个成员没有意义。

### 重定位表

SimpleSection.o中有一个叫做`.rela.text`的段，它的类型（sh_type）为`SHT_RELA`，也就是说它是一个重定位表（Relocation Table）。链接器在处理目标文件时，须要对目标文件中某些部位进行重定位，即代码段和数据段中那些对绝对地址的引用的位置。这些重定位的信息都记录在ELF文件的重定位表里面，对于每个须要重定位的代码段或数据段，都会有一个相应的重定位表。`.rela.text`就是针对`.text`段的重定位表，因为`.text`段中至少有一个绝对地址的引用（即对`printf`函数的调用）。

重定位表的`sh_link`表示符号表的下标，`sh_info`表示它作用于哪个段。比如`.rela.text`作用于`.text`段（`.text`段的下标为1），那么`.rel.text`的`sh_info`为1。

### 字符串表

ELF文件中用到了很多字符串，比如段名、变量名等。因为字符串的长度往往是不定的，所以用固定的结构来表示它比较困难。一种**很不错的做法**是把字符串集中起来存放到一个表，然后使用字符串在表中的偏移来引用字符串，如下所示。

![string_table](./.assets/string_table.png)

通过这种方法，在ELF文件中引用字符串只须给出一个数字下标即可。字符串表在ELF文件中也以段的形式保存，常见为`.strtab`字符串表（String Table）或`.shstrtab`段表字符串表（Section Header String Table）。

---

*参考 《程序员的自我修养 --链接、装载与库》*