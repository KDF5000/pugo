```toml
title = "Golang第七章:接口"
date = "2017-04-24 16:34:48"
update_date = "2017-04-24 16:34:48"
author = "KDF5000"
thumb = ""
tags = ["Golang", "学习笔记"]
draft = false
```
* 接口更像是一种约定，范式满足接口约定的形式的类型都可以作为该接口的实例
* 接口类型具体描述了一系列方法的集合，一个实现了这些方法的具体类型是这个接口类型的实例
* 接口类型可以进行组合成新的接口类型，这种也叫接口内嵌
```go
package io
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Closer interface {
    Close() error
}
//
type ReadWriter interface {
    Reader
Writer }
type ReadWriteCloser interface {
    Reader
Writer
Closer }
```

<!--more-->

* 一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口
* 表达一个类型属于某个接口只要这个类型实现了这个接口
```go
var w io.Writer
w = os.Stdout
w = new(bytes.Buffer)
w = time.Second
// OK: *os.File has Write method
// OK: *bytes.Buffer has Write method
// compile error: time.Duration lacks Write method
var rwc io.ReadWriteCloser
rwc = os.Stdout         // OK: *os.File has Read, Write, Close methods
rwc = new(bytes.Buffer) // compile error: *bytes.Buffer lacks Close method
```
* 一个具体类型赋值给一个接口类型的变量后，即使该具体类型有其他方法也不能调用
* interface{}空类型对实现它的类型没有任何要求，所以可以将任意一个值赋给空接口类型
* flag.Value可以为自己的数据类型位flag添加新的标记符号


### sort.Interface接口
* go语言的sort函数不会对具体的序列和它的元素做任何的假设
* 排序算法通常需要知道的三个要素
 - 序列的长度 --Len() int
 - 两个元素比较的结果 - Less(i int, j int)bool
 - 交换两个元素的方式 -- Swap(i,j int)
 
 ```go
 ///
 package sort
type Interface interface {
    Len() int
    Less(i, j int) bool // i, j are indices of sequence elements
    Swap(i, j int)
}
///////
type StringSlice []string
func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] 
//调用
sort.Sort(StringSlice(names))
 ```
 
 
 ### http.Handler接口
* 在package http里
```go
package http
type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}
func ListenAndServe(address string, h Handler) error
```
* 一个最简单的例子
```go
func main() {
    db := database{"shoes": 50, "socks": 5}
    log.Fatal(http.ListenAndServe("localhost:8000", db))
}
type dollars float32
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }
type database map[string]dollars
func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) { for item, price := range db {
        fmt.Fprintf(w, "%s: %s\n", item, price)
    }
}
```
* 如果要想对不同url执行不同的操作，则可以在ServeHttp里通过req.URL.Path获得请url然后通过switch等方式判断执行不同的操作
* go也提供了请求多路ServeMux来简化URL和handler的对应
```go
func main() {
    db := database{"shoes": 50, "socks": 5}
    mux := http.NewServeMux()
    mux.Handle("/list", http.HandlerFunc(db.list))
    mux.Handle("/price", http.HandlerFunc(db.price))
    log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
type database map[string]dollars
func (db database) list(w http.ResponseWriter, req *http.Request) {
    for item, price := range db {
        fmt.Fprintf(w, "%s: %s\n", item, price)
    }
}
func (db database) price(w http.ResponseWriter, req *http.Request) {
    item := req.URL.Query().Get("item")
    price, ok := db[item]
    if !ok {
        w.WriteHeader(http.StatusNotFound) // 404
        fmt.Fprintf(w, "no such item: %q\n", item)
        return
}
    fmt.Fprintf(w, "%s\n", price)
}
```
其中http.HandlerFunc其实是一个实现了http.Handler ServeHttp方法的函数类型
```go
package http
type HandlerFunc func(w ResponseWriter, r *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```
* 如果有很多handler，有可能需要放在不同的文件中，这个时候可以使用http包提供的默认的ServeMux示例DefaultServeMux和包级http.Handle和http.HandleFunc函数， 这个时候ListenAndServe函数的handler值为空就可以
```go
func main() {
    db := database{"shoes": 50, "socks": 5}
    http.HandleFunc("/list", db.list)
    http.HandleFunc("/price", db.price)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
```
 
### error接口
* error接口
```go
type error interface {
    Error() string
}
```
* errors包的error的一个实现
```go
package errors
func New(text string) error { return &errorString{text} }
type errorString struct { text string }
func (e *errorString) Error() string { return e.text }
```
### 类型断言
* x.(T)判断x是否是T类型，如果成功则返回x的的动态值，如果失败则跑出panic
