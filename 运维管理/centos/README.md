# centos（centos8）

## 安装centos后无网络

修改配置文件：```vim /etc/sysconfig/network-script/ifcfg-ens33```

```txt
# 开机启用本配置
ONBOOT=yes
```

## 设置静态IP地址

修改配置文件：```vim /etc/sysconfig/network-script/ifcfg-ens33```

```txt
BOOTPROTO="static"      # dhcp改为static 
IPADDR=192.168.7.106    # 静态IP
GATEWAY=192.168.7.1     # 默认网关
NETMASK=255.255.255.0   # 子网掩码
DNS1=192.168.7.1        # DNS 配置
```

## bash报错："Failed to set locale, defaulting to C.UTF-8"

使用如下命令

```shell
dnf install langpacks-en glibc-all-langpacks -y  # langpacks-zh
localectl set-locale LANG=en_US.UTF-8  # LANG=zh_CN.UTF-8
```