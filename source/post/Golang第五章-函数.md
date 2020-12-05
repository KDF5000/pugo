```toml
title = "Golang第五章:函数"
date = "2017-03-11 20:58:00"
update_date = "2017-03-11 20:58:00"
author = "KDF5000"
thumb = ""
tags = ["Golang", "学习笔记"]
draft = false
```
* 函数的类型被称为函数的标识符，如果两个函数的形参列表和返回值列表中的变量类型一一对应，那么这个函数被认为有相同的类型和标识符。形参和返回值的变量名不影响含糊是标识符，也不影响他们是否可以省略参数类型的形式表示。
* Go语言没有默认参数值
* Go可以有有名返回值，和形参一样作为局部变量，被存储在相同的词法块，一般使用带有变量名的返回值主要为了表示返回值的含义,也可以直接在函数体对这些变量赋值，返回时候可以省略return的操作数
```go
func Size(rect image.Rectangle) (width, height int)
```
* Go的参数传递都是值传递，因此对形参的修改不会影响实参，但是如果实参是指针，slice,map, function,channel等类型，则可以间接的影响实参值。
* 下面没有函数体的函数声明，表明该函数不是Go事先的，这样的声明定义了函数标识符
```go
package math
func Sin(x float64) float //implemented in assembly language
```

## 错误处理
* error类，是个接口类型，有nil和non-nil两种
* 错误处理策略
 - 传播错误
    - 将错误err传递给调用者，然后决定怎么处理
    - 错误信息里最好包含上下文信息，通常是连式组合的形式，所以也要避免大写和换行符
 - 重试
  设置超时时间，重新尝试操作
 - 结束进程
  如果错误会导致程序无法继续运行，则输出错误结束程序,os.Exit(1)
 - Log记录错误信息，继续运行
 - 忽略错误              
* io包提供一个错误io.EOF,保证任何文件结束引起的读取失败错误。

## 函数值
* Go中函数被看做第一类值(first-class values)： 函数像其他值一样，拥有类型，可以赋值给其他变量，传递给函数，从函数返回
```go
func square(n int) int { return n * n  }
func negative(n int) int { return -n  }
func product(m, n int) int { return m * n  }
f := square
fmt.Println(f(3)) // "9"
f = negative
fmt.Println(f(3))     // "-3"
fmt.Printf("%T\n", f) // "func(int) int"
f = product // compile error: can't assign func(int, int) int to func(int) int
```
* 函数类型的零值为nil, 调用值为nil的函数值会引起panic错误
* 函数值不能进行比较，不能作为map的key
* 函数可以作为函数的参数
```go
func add(a int , b int) int{return a+b}
func sub(a int , b int) int {return a-b}
func compute(a, b int, fun func(int,int)int){
    return fun(a,b)

}
```

## 匿名函数
* 匿名函数在Go中称为Function literal, 声明方式和函数类似，但是没有函数名，是一种表达式，他的值被称为匿名函数
* 匿名函数可以在使用时定义
* 匿名函数可以访问完整的词法环境，也就是说在一个函数内部定义匿名函数，在匿名函数内部可以访问函数的变量
```Go
func squares() func() int {
  var x int
      return func() int {
              x++
                      return x * x
                          
      }

}
func main() {
    f := squares()           //f指向一个匿名函数，该匿名函数里的X是一个变量
        fmt.Println(f())         //1
            fmt.Println(f())         //4
                fmt.Println(f())         //9
                    fmt.Println(squares()()) //1
                        fmt.Println(squares()()) //1

}
```
对squares的一次调用会生成一个局部变量x并返回一个匿名函数，每次调用改匿名函数，该函数都会使x得值加1， 上面f指向匿名函数，是同一个局部变量x。第一次调用squares（）,会生成第二个x变量。并返回一个新的匿名函数

这个例子证明，函数值不仅仅是一串代码，还记录了状态。在square中定义的匿名内部函数可以访问和更新square中的局部变量，这意味着匿名函数和square中，存在变量引用。这就是**函数值属于引用类型和函数值不可比较**的原因。Go使用闭包技术实现函数值，Go程序员也把函数值叫做闭包。

### 捕获迭代变量
先看下面的程序
```go
var rmdirs []func()
    for _, d := range tempDirs() {
        os.MkdirAll(dir, 0755) // creates parent directories too
            rmdirs = append(rmdirs, func() {
                    os.RemoveAll(d)
                        
                    })

    }
    // ...do some work...
for _, rmdir := range rmdirs {
    rmdir() // clean up

}
```
上面程序的执行能否达到预期的结果呢？答案是否定的。for循环内部创建一个新的目录的时候，我们使用一个匿名的函数，在循环内部引用了变量d, d其实只是指向变量的内存地址，每一次循环的过程都在发生改变，因此for循环结束后指向的是最后一个变量的地址，所以后面for循环的时候删除的都是同一个目录。解决方法就是在for循环内部对循环变量进行一个copy即可，见下面代码
```go
var rmdirs []func()
    for _, d := range tempDirs() {
        dir := d // NOTE: necessary!
            os.MkdirAll(dir, 0755) // creates parent directories too
            rmdirs = append(rmdirs, func() {
                    os.RemoveAll(dir)
                        
                    })

    }
    // ...do some work...
for _, rmdir := range rmdirs {
    rmdir() // clean up

}
```

### 可变参数
声明可变参数函数时，需要在参数列表的最后一个参数类型之前加上省略号"...",这表示该函数会接受任意数量的该类型参数。
```go
func sum(vals...int) int {
    total := 0
               for _, val := range vals {
                       total += val

               }
                   return total

}
```
对于可变参数，可以直接传递一个切片，然后后面加省略号
```go
values := []int{1, 2, 3, 4}
fmt.Println(sum(values...)) // "10"
```

interfac{}表示函数的最后一个参数可以接受任意类型

### Deferred函数
* defer可以在异常判断的过程中多次执行诸如文件关闭的操作，使用defer可以在return的时候自动执行
* 调试复杂程序时，defer机制也常被用于记录何时进入和退出函数
```go
func bigSlowOperation() {
    defer trace("bigSlowOperation")() // don't forget the
        extra parentheses
            // ...lots of work...
                time.Sleep(10 * time.Second) // simulate slow
                    operation by sleeping

}
func trace(msg string) func() {
    start := time.Now()
        log.Printf("enter %s", msg)
        return func() {
                log.Printf("exit %s (%s)", msg,time.Since(start))
                    
        }

}
```
* defer语句中的函数会在return语句更新返回值后再执行，又因为函数中定义的匿名函数可以访问该函数包括返回值变量在内的所有变量，所以，对匿名函数采用defer机制可以使其观察函数的返回值. 
```go
func double(x int) (result int) {
    defer func() { fmt.Printf("double(%d) = %d\n", x,result)  }()
        return x + x

}
```
* 被延迟执行的匿名函数甚至可以修改函数返回给调用者的返回值
```go
func triple(x int) (result int) {
    defer func() { result += x  }()
        return double(x)

}
fmt.Println(triple(4)) // "12"
```

* 对于读写文件信息的操作可能要慎用defer操作，通过os.Create打开的文件进行写入，在关闭文件时，不能采用defer机制。因为对于许多文件系统，尤其是NFS， 写入文件时发生的错误会被延迟到文件关闭时反馈。如果没有检查文件关闭时的反馈信息，可能会导致数据丢失，而我们还误以为写入操作成功。
* 延迟函数的调用在释放堆栈信息之前

### Panic异常
* panic函数接受任何值作为参数，当某些不应该发生的场景发生时，我们就应该调用painc。
* 如果在deferred函数中调用了内置函数recover, 并且定义改defer语句的函数发生了panic异常，recover会使程序从panic中恢复，并返回panic value. 导致panic异常的函数不会继续运行，但能整场 返回。在未发生panic时调用recover， recover会返回nil
```go
func Parse(input string) (s *Syntax, err error) {
    defer func() {
        if p := recover(); p != nil {
                    err = fmt.Errorf("internal error: %v", p)

        } 
    }()
        // ...parser...

}
```
