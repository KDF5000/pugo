```toml
title = "Golang第三章:基础数据类型"
date = "2017-02-25 11:59:30"
update_date = "2017-02-25 11:59:30"
author = "KDF5000"
thumb = ""
tags = ["Golang", "学习笔记"]
draft = false
```
### 整型
有符号：int8, int16, int32, int64, int
无符号: uint8, unit16, uint32, uint64, uint

```go
var x uint8 = 1<<1 | 1<<5 //00100010
var y uint8 = 1<<1 | 1<<2 //00000110
fmt.Printf("%08b\n", x)   // "00100010", the set {1, 5}
fmt.Printf("%08b\n", y)   // "00000110", the set {1, 2}

fmt.Printf("%08b\n", x&y)  // "00000010", the intersection {1}
fmt.Printf("%08b\n", x|y)  // "00100110", the union {1, 2, 5} //或
fmt.Printf("%08b\n", x^y)  // "00100100", the symmetric difference {2, 5}不同为1
fmt.Printf("%08b\n", x&^y) // "00100000", the difference {5} 如果y对应的位为1，则为0；否则为x的值
```

<!-- more -->

下面的例子
```go
medals := []string{"gold", "silver", "bronze"}
for i := len(medals) - 1; i >= 0; i-- {
    fmt.Println(medals[i]) // "bronze", "silver", "gold"
}
```
如果i是一个无符号数，则i--永远不会小于0，因此将会出现out ot index的错误。所以len函数设计时候返回的是有符号的int

```go
o := 0666
fmt.Printf("%d %[1]o %#[1]o\n", o) // "438 666 0666"
x := int64(0xdeadbeef)
fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", x)
// Output:
// 3735928559 deadbeef 0xdeadbeef 0XDEADBEEF
```

### 浮点数
两种：float32和float64
max.MaxFloat32: 表示float32能表示的最大数。大约3.4e38
max.MaxFloat64 大约1.8e308

### 布尔类型
布尔值不会隐式的转换为0或者1，需要自己转换。同样不像C语言中，大于0的整数就是true，必须显示的进行转换。

### 字符串
* 字符串的值不能改变，也就是说字符串只读
* 各种编码，ASCII, Unicode, UTF-8
* 字符串和Byte切片

### 常量
常量表达式的值在编译期计算，不是在运行期，常量的值不可修改。
```go
const pi = 3.14159
const(
    e = 2.3232322222242424
    p = 3.1415926533232323
)
```
常量间的所有算术运算，逻辑运算和比较结果也是常量，对常量的类型转换或者以下函数的调用结果也是返回常量结果：len, cap, real, imag, complex和unsafe.Sizeof

在一个const声明语句中，在第一个声明的常量所在行，iota将会被置0，然后在每一个有常量声明的行加一

