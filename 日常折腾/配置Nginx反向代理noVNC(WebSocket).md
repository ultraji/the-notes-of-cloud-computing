# 配置Nginx反向代理noVNC（WebSocket）

最近在自己的服务器中跑了一个容器运行novnc服务，对外开放的端口为6080。通过6080端口可以直接访问noVNC服务，而我在nginx中用常规的方式反向代理了noVNC，却始终访问不通。

## WebSocket代理

要将客户端和服务器之间的连接从HTTP / 1.1转换为WebSocket，使用HTTP / 1.1中提供的协议切换机制。  

然而有一个微妙之处：由于“Upgrade”是一个逐跳的头，它不会从客户端传递到代理服务器。使用正向代理，客户可以使用该CONNECT 方法来规避这个问题。但是，这不适用于反向代理，因为客户端不知道任何代理服务器，并且需要在代理服务器上进行特殊处理。  

版本1.3.13开始，nginx实现了特殊的操作模式，如果代理服务器返回了代码101（交换协议）的响应，客户端和代理服务器之间建立隧道，客户端通过请求中的“Upgrade”请求头。  

如上所述，包括“Upgrade”和“Connection”的逐跳标题不会从客户端传递到代理服务器，因此为了让代理服务器知道客户端将协议切换到WebSocket的意图，这些标题必须明确地通过：

```txt
http {
    map $http_upgrade $connection_upgrade {
        default upgrade;
        ''      close;
    }

    server {
        listen 80; #修改监听的端口
        server_name _;
        location / {
            proxy_pass  #修改为需要被反向代理的WebSocket的IP和端口号
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
        }
    }
```