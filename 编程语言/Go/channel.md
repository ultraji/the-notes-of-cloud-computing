# channel

## channel的使用场景

1. 消息传递、消息过滤
2. 信号广播
3. 事件订阅与广播
4. 请求、响应转发
5. 任务分发
6. 结果汇总
7. 并发控制
8. 同步与异步
9. ...

## 注意事项

channel存在**3种状态**：

1. nil，未初始化的状态，只进行了声明，或者手动赋值为nil
2. active，正常的channel，可读或者可写
3. closed，已关闭，千万不要误认为关闭channel后，channel的值是nil

channel可进行**3种操作**：

1. 读
2. 写
3. 关闭

把这3种操作和3种channel状态可以组合出9种情况：

| 操作 | nil的channel | 正常channel | 已关闭channel |
| --- | --- | --- | --- |
| <- ch | 阻塞 | 成功或阻塞 | 读到零值 |
| ch <- | 阻塞 | 成功或阻塞 | panic |
| close(ch) | panic | 成功 | panic |

## 常用操作

### 使用`for range`读channel

- **场景**：当需要不断从channel读取数据时
- **原理**：使用for-range读取channel，这样既安全又便利，当channel关闭时，for循环会自动退出，无需主动监测channel是否关闭，可以防止读取已经关闭的channel，造成读到数据为通道所存储的数据类型的零值。
- **用法**：
    ```go
    for x := range ch{
        fmt.Println(x)
    }
    ```

### 使用`_, ok`判断channel是否关闭

- **场景**：读channel，但不确定channel是否关闭时
- **原理**：读已关闭的channel会得到零值，如果不确定channel，需要使用ok进行检测。ok的结果和含义：
    - `true`：读到数据，并且通道没有关闭。
    - `false`：通道关闭，无数据读到。
- **用法**：
    ```go
    if v, ok := <- ch; ok {
        fmt.Println(v)
    }
    ```

### 使用select处理多个channel

- **场景**：需要对多个通道进行同时处理，但只处理最先发生的channel时
- **原理**：select可以同时监控多个通道的情况，只处理未阻塞的case。当通道为nil时，对应的case永远为阻塞，无论读写。特殊关注：普通情况下，对nil的通道写操作是要panic的。
- **用法**：
    ```go
    // 分配job时，如果收到关闭的通知则退出，不分配job
    func (h *Handler) handle(job *Job) {
        select {
        case h.jobCh<-job:
            return 
        case <-h.stopCh:
            return
        }
    }
    ```

### 使用channel的声明控制读写权限

- **场景**：协程对某个通道只读或只写时
- **目的**：A. 使代码更易读、更易维护，B. 防止只读协程对通道进行写数据，但通道已关闭，造成panic。
- **用法**：
    如果协程对某个channel只有写操作，则这个channel声明为只写。
    如果协程对某个channel只有读操作，则这个channe声明为只读。
    ```go
    // 只有generator进行对outCh进行写操作，返回声明
    // <-chan int，可以防止其他协程乱用此通道，造成隐藏bug
    func generator(int n) <-chan int {
        outCh := make(chan int)
        go func(){
            for i:=0;i<n;i++{
                outCh<-i
            }
        }()
        return outCh
    }

    // consumer只读inCh的数据，声明为<-chan int
    // 可以防止它向inCh写数据
    func consumer(inCh <-chan int) {
        for x := range inCh {
            fmt.Println(x)
        }
    }
    ```

### 为操作加上超时

- **场景**：需要超时控制的操作
- **原理**：使用select和time.After，看操作和定时器哪个先返回，处理先完成的，就达到了超时控制的效果
- **用法**：
    ```go
    func doWithTimeOut(timeout time.Duration) (int, error) {
        select {
        case ret := <-do():
            return ret, nil
        case <-time.After(timeout):
            return 0, errors.New("timeout")
        }
    }

    func do() <-chan int {
        outCh := make(chan int)
        go func() {
            // do work
        }()
        return outCh
    }
    ```






## 附录

参考资料：

- [总结了才知道，原来channel有这么多用法！](https://segmentfault.com/a/1190000017958702)