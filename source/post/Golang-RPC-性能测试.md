```toml
title = "Golang RPC 性能测试"
date = "2017-03-28 16:44:27"
update_date = "2017-03-28 16:44:27"
author = "KDF5000"
thumb = ""
tags = ["Golang", "RPC", "Benchmark", "性能测试"]
draft = false
```
最近刚好要使用Golang的RPC，因此对Golang标准库的RPC进行了一下测试，看看其性能到底如何。RPC服务端和客户端的实现完全使用RPC的`net/rpc`标准库，没有经过特殊的优化，主要针对下面三个场景进行测试。测试之前需要先说明一下，Go的rpc连接是支持并发请求的，就是说一个一个连接可以并发的发送很多个请求，不像http协议一问一答的模式。
### 测试环境
操作系统：Centos 6.8 (Linux 2.6.32)
内存：32G
核数：双CPU, 一共12核
CPU型号：Intel(R) Xeon(R) CPU E5645  @ 2.40GHz
Golang: 1.7.4

<!--more-->

### 场景
测试的场景主要是下面两个指标
* QPS指标
 * 单个个连接保证一个并发，随着该并发请求数增加，QPS的变化
 * 单个连接(Client), 单个并发请求10w, 随着并发数的增加，QPS的变化
 * 单个连接并发数固定(第一个测试的最优值)，增加连接数，QPS的变化
 
* 单机Server的并发数(同时连接数)
 * 单机Server, 测试所能接收的连接数

QPS指标测试中，第一个设置是为了测试单个连接的并发数
### 实现
Server端的实现使用tcp协议，监听4200端口，循环等待连接，每当检测到请求时，启动一个goroutine去处理该连接，注册的服务执行一个简单的乘法操作。
```Go
//Service
type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func (t *Arith) Divide(args Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("Divided by zero!")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}
//////////////////////////////////////////////////
//Server
runtime.GOMAXPROCS(4)
arith := new(service.Arith)
server := rpc.NewServer()
log.Printf("Register service:%v\n", arith)
server.Register(arith)

log.Printf("Listen tcp on port %d\n", 4200)
l, e := net.Listen("tcp", ":4200")

if e != nil {
    log.Fatal("Listen error:", e)
}
log.Println("Ready to accept connection...")
conCount := 0
go func() {
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal("Accept Error:,", err)
            continue
        }
        conCount++
        log.Printf("Receive Client Connection %d\n", conCount)
        go server.ServeConn(conn)
    }
}()
```
Client端的创建多个连接(Client), 然后每个连接指定不同的并发数，每个并发启动一个goroutine发送指定数量的request，该请求执行一个简单的乘法操作，最后统计整个过程的QPS。

### 部署
QPS指标的测试使用两台上述配置的服务器，然后设置相关的内核参数，主要是允许最大打开的文件数，可使用的端口范围，tcp缓存大小等,操作如下：
```shell
sysctl -w fs.file-max=10485760 #系统允许的文件描述符数量10m
sysctl -w net.ipv4.tcp_rmem=1024 #每个tcp连接的读取缓冲区1k，一个连接1k
sysctl -w net.ipv4.tcp_wmem=1024 #每个tcp连接的写入缓冲区1k

sysctl -w net.ipv4.ip_local_port_range='1024 65535' #修改默认的本地端口范围
sysctl -w net.ipv4.tcp_tw_recycle=1  #快速回收time_wait的连接
sysctl -w net.ipv4.tcp_tw_reuse=1
sysctl -w net.ipv4.tcp_timestamps=1

#用户单进程的最大文件数，用户登录时生效
echo '* soft nofile 1048576' >> /etc/security/limits.conf 
echo '* hard nofile 1048576' >> /etc/security/limits.conf 
ulimit -n 1048576 #用户单进程的最大文件数 当前会话生效
```
对于Server连接数的测试，使用多台机器，尽量保证同时进行连接同一个server

### 测试结果
#### QPS指标
* 单个连接保证一个并发，随着该一个并发请求数增加，QPS的变化
![单个连接请求数.png](@media/archive/blog/images/单个连接请求数.png)
这个测试可以评估单个连接单并发的情况下，能过处理的最多请求数。随着请求数的增加，QPS一直下降，但是下降到6k左右的时候讲保持这个值不再改变，从图中可以大概在请求量大概在10k以后qps将不再改变，也就是说按照现在这个频率(无休止持续请求),client的处理能力大概就是6k/s，可以根据这个指标控制请求的频率。

* 单个连接(Client), 单个并发请求10w, 随着并发数的增加，QPS的变化
![单个连接并发数.png](@media/archive/blog/images/单个连接并发数.png)
这个测试可以指导客户端在设计单个连接的并发请求数时候怎么选择最佳的并发数。从测试结果可以看出，单个连接随着并发数的增加，QPS(RPS)并不是线性增长的，基本上增加到100个并发数就基本不再增加了，说明对于单个连接最大的QPS大概是6w多点。而且对于单个连接来说，并发数最好不要超过150。

* 单个连接并发数固定(第一个测试的最优值)，增加连接数，QPS的变化
![固定并发连接数.png](@media/archive/blog/images/固定并发连接数.png)
这里设置单个连接数的并发数为100，每个并发的请求量位10w, 从结果可以看出随着连接数的增加，Qps也是在不断增加的,但是当增加到100个连接数后，Qps基本不再变化，维持在47w左右，这个值相对还是比较大的，这说明Go的RPC client和Server的性能还是不错的。

#### 单机Server的并发数(同时连接数)
* 单机Server, 测试同时能接收的连接数
 因为这个涉及到多台部署了Client的服务器进行同时请求，暂时手里服务器资源还不够，过段时间腾出来一些机器了补充起来，**如果有单机可以模拟(单机局限于端口数)的方案望转告**。

### 总结
上面的前三个实验主要关注Client端的性能，因为自己对服务端的压测指标，业务场景还不是非常明确，所以对Server端的压测试还不够充分，以后会逐步的补充起来。从测试结果来看，Go RPC的Client的性能还是不错的。

* 单连接的满负载Qps大概是6k
* 单连接的并发数最好不要超过150
* 单个机器的的连接数最好不要超过100
