# Powershell

## PowerShell因为在此系统上禁止运行脚本

首次在计算机上启动Windows PowerShell时，默认设置的执行策略很可能是Restricted：
    
- Restricted 执行策略不允许任何脚本运行。 
- AllSigned 和 RemoteSigned 执行策略可防止 Windows PowerShell 运行没有数字签名的脚本。

具体操作：

1. 获取计算机上的现用执行策略，使用命令`get-executionpolicy`；

2. （需管理员权限运行powershell）将执行策略更改为`RemoteSigned`， 使用命令`set-executionpolicy remotesigned`。