# bochs与gdb联调时忽略page fault信号

bochs与gdb联调，在调试内核时经常会被page fault（signal 0）信号打断，如何忽略page fault信号呢？

1. 根据以下内容，修改bochs源码中的gdbstub.cc文件

    ```txt
    Created a patch "gdbstub.cc.patch" against bochs (version CVS 20080110)
    Bochs always tries to find out the reason of an exception, so that it can generate the right signal for gdb.
    If it fails to find a reason, bochs assigns a value GDBSTUB_STOP_NO_REASON (see bochs.h), which causes
    debug_loop() (see gdbstub.cc) to generate a signal of number 0.
    Signal 0 is problematic to gdb, as gdb doesn't allow us to ignore it.
    Somehow when we simulate linux, we get tons of signal 0's that seem to be caused by page faults.
    This patch makes bochs send SIGSEGV instead of signal 0, so that we can ignore it in gdb.

    gdbstub.cc
    debug_loop()

    *** gdbstub.cc.orig Thu Oct 18 18:44:38 2007
    --- gdbstub.cc Sat Jan 12 17:25:22 2008
    *************** static void debug_loop(void)
    *** 489,494 ****
    --- 489,498 ----
            {
                write_signal(&buf[1], SIGTRAP);
            }
    +  else if (last_stop_reason == GDBSTUB_STOP_NO_REASON)
    +  {
    +    write_signal(&buf[1], SIGSEGV);
    +  }
            else
            {
                write_signal(&buf[1], 0);
    ```

2. 重新编译生成支持gdb调试的bochs二进制文件

    ```shell
    ./configure --enable-gdb-stub --enable-disasm
    make
    ```

## 附录

1. 下载bochs源码

    ```shell
    wget https://downloads.sourceforge.net/project/bochs/bochs/2.6.9/bochs-2.6.9.tar.gz -q --show-progress
    ```

2. 安装编译bochs的必要依赖

    ```shell
    sudo apt-get install -y build-essential libgtk2.0-dev libx11-dev xserver-xorg-dev xorg-dev g++ pkg-config libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev
    # 选装
    sudo apt-get install -y bochs bochs-x bochs-sdl 
    ```