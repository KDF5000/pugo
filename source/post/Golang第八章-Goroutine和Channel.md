```toml
title = "Golang第八章: Goroutine和Channel"
date = "2017-05-07 20:36:18"
update_date = "2017-05-07 20:36:18"
author = "KDF5000"
thumb = ""
tags = ["Golang", "学习笔记"]
draft = false
```

并发编程模型
 - 顺序通信进程(Communicating Sequential Processes) CSP
  值会在不同的实例(goroutine)中传递
 - 多线程共享内存
 
 <!--more -->

### Goroutines
* Go语言中每一个并发的执行单元叫做一个goroutine
* Go语言中主函数即在一个单独的goroutine中运行
* 所有的goroutine在主函数返回时，都会直接打断，程序退出

简单的clock程序[clock.go](ch8/clock.go)
```go
package main

import (
    "flag"
    "fmt"
    "io"
    "log"
    "net"
    "time"
)

var (
    host string
    port int
)

func init() {
    flag.StringVar(&host, "host", "localhost", "clock server host")
    flag.IntVar(&port, "port", 8000, "clock server port")
    flag.Parse()
}

func handleCon(c net.Conn) {
    defer c.Close()
    // for {
    _, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
    if err != nil {
        return
    }
    time.Sleep(1 * time.Second)
    // }
}

func main() {
    server := fmt.Sprintf("%s:%d", host, port)
    listener, err := net.Listen("tcp", server)
    if err != nil {
        log.Fatal(err)
    }
    for {
        conn, err := listener.Accept() // 会阻塞
        if err != nil {
            log.Print(err)
            continue
        }
        go handleCon(conn)
    }
}
```
[client.go](ch8/netcat.go)：
```go
package main

import (
    "io"
    "log"
    "net"
    "os"
    "sync"
)

func main() {
    servers := []string{"localhost:8001", "localhost:8002", "localhost:8003"}
    var wg sync.WaitGroup
    for _, s := range servers {
        wg.Add(1)
        go func(serv string) {
            conn, err := net.Dial("tcp", serv)
            if err != nil {
                log.Fatal(err)
            }

            if _, err := io.Copy(os.Stdout, conn); err != nil {
                log.Fatal(err)
            }
            conn.Close()
            wg.Done()
        }(s)
    }
    wg.Wait()
}
```
### Channels
* goroutine是并发单元，可以通channels实现他们之间的通信
* 创建一个channel
```go
ch := make(chan int)
```
channel是一个底层数据结构的引用，当作为参数传递是，是传递的一个引用
* 两个channel可以用==运算比较，如果两个channel引用的是相同的对象，那么比较的结果为真，channel的零值是nil
* 关闭一个channel后，继续发送会产生panic，但是可以继续接收，此时将不会阻塞，立即返回一个零值
* 没有办法直接测试一个channel是否被关闭，但是接收操作有一个变体形式，它可以多接收一个结果，多接收的第二个结果是一个布尔值Ok, true表示成功从channel接收到了值，false便是已经关闭 并且里面没有值可以接收
    ```go
    x,ok := <- ch
    if !ok{
        //do something
    }
    //使用range简化
    for x := range ch{
        //do something
    }
    ```
* channel 可以不用显式的关闭，当他没有被引用时go语言的垃圾回收会自动回收
* 单向channel
 - chan<- int 只用来发送。 向channel发送数据
 - <- chan int 只接收。接收channel的数据
 * 关闭channel用于断言不再向channel发送数据，因此对于一个只接收的channel调用close将是一个编译错误
* 可以将一个双向的channel赋值给单向channel变量，会做隐式的转换。但是不能将一个单向的channel转换为双向的channel

### 基于select的多路复用
* select会等待case中能够执行的case时去执行，当条件满足时，才会去通信并执行case之后的语句
* 一个没有任何case的select语句写作select{}，会永远等待下去
* 如果多个case同时就绪，select会随机选择一个执行，这样来保证每一个channel都有相等的被select执行的机会,下面的例子中channel的缓冲设为1，这样每次循环的时候ch的状态为空或者满，偶数的时候恰好为空，奇数时候为满，所以输出0，2，4，6，8。如果缓冲区设为大于1的数，那么select就会随机选择，结果就不确定了。
    ```go
    func main() {
        ch := make(chan int, 1)
        for i := 0; i < 10; i++ {
            select {
            case x := <-ch:
                fmt.Println(x)
            case ch <- i:
                //
            }
        }
    }
    ```
* Time.Tick函数表现的像是它创建一个在循环中调用time.Sleep的goroutine每次被唤醒时发送一个事件，依然会不断的尝试向channel中发送值，如果没有接受方去接受，那么就会造成goroutine泄露，因此只有当程序整个生命周期都需要这个时间时我们使用它才比较合适，否则建议使用下面的模式：
    ```go
    ticker := time.NewTicker(1 * time.Second)
    <-ticker.C    // receive from the ticker's channel
    ticker.Stop() // cause the ticker's goroutine to terminate
    ```
* 对一个nil的channel发送和接受操作会永远阻塞，在select语句中操作nil的channel永远都不会被select到


### 实例 - 并发的目录遍历
下面是一个并发遍历目录的实例
```go
package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
    "time"
)

func walkDir(dir string, wg *sync.WaitGroup, filesize chan<- int64) {
    defer wg.Done()
    for _, entry := range dirents(dir) {
        if entry.IsDir() {
            wg.Add(1)
            subdir := filepath.Join(dir, entry.Name())
            go walkDir(subdir, wg, filesize)
        } else {
            filesize <- entry.Size()
        }
    }
}

var sema = make(chan struct{}, 20) //最多打开20个目录

func dirents(dir string) []os.FileInfo {
    sema <- struct{}{}
    defer func() { <-sema }()
    entries, err := ioutil.ReadDir(dir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "du1: %v\n", err)
        return nil
    }
    return entries
}

var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
    flag.Parse()
    roots := flag.Args()
    if len(roots) == 0 {
        roots = []string{"."}
    }
    filesizes := make(chan int64)
    var wg sync.WaitGroup
    for _, dir := range roots {
        wg.Add(1)
        go walkDir(dir, &wg, filesizes)
    }

    go func() {
        wg.Wait()
        close(filesizes)
    }()
    var tick <-chan time.Time
    if *verbose {
        tick = time.Tick(500 * time.Millisecond)
    }
    var nfiles, nsize int64
loop:
    for {
        select {
        case size, ok := <-filesizes:
            if !ok {
                break loop
            }
            nfiles++
            nsize += size
        case <-tick:
            fmt.Printf("%d files %1.fG\n", nfiles, float64(nsize)/1e9)
        }
    }
}
```

### 并发的退出
* Go语言不提供在一个goroutine里终止另一个goroutine的方法，由于这样会导致goroutine之间共享变量落在未定义的状态上。
* 可以通过关闭一个channel进行广播关闭所有goroutine

### 一个简单的聊天程序
[聊天服务器](https://github.com/KDF5000/gopl/blob/master/src/ch8/chat.go)+[客户端](https://github.com/KDF5000/gopl/blob/master/src/ch8/echoclient.go)
