```toml
title = "Golang第十一章-测试"
date = "2017-03-28 00:39:43"
update_date = "2017-03-28 00:39:43"
author = "KDF5000"
thumb = ""
tags = ["Golang", "学习笔记"]
draft = false
```
* 命令
 * go test
 * 遍历*_test.go 文件中复合规则的函数
* 类型
 * 测试函数：Test开头用于测试程序逻辑行为的正确性
 * 基准测试函数：Benchmark开头，用于衡量函数的性能，go test会运行多次取平均的执行时间
 * 示例函数：Examole开头的函数，提供一个有编译器保证正确性的示例文档
 
## 测试函数
* 导入的包和形式如下
```
import "testing"
func TestName(t *testing){
    //...
}
```
* go test -v参数可以打印每个测试函数的名字和运行时间
* go test -run 对应一个正则表达式，只有测试函数名被它正确匹配的测试函数才会被go test测试命令运行。
* t.Error不能终止测试，t.Fatal可以终止测试
* t.Fatal必须和测试函数在同一个groutine里调用才能终止测试
* 一般测试失败的信息形式为`f(x)=y, want z`

<!--more-->

### 随机测试
### 白盒测试

* 在测试一些诸如邮件发送函数的时候，我们并不想真正的发送邮件，因此可以将发送函数作为一个包级的私有函数值，然后在测试代码里先用修改该函数值，进行同等效果的测试
```go
//fakefunc.go
package ch11
import (
    "fmt"
)
var realFunc = func(data string) {
    fmt.Printf("Real Func:%s\n", data)
}
func CheckInfo(data string) {
    realFunc(data)
}
//
//fakefunc_test
package ch11
import (
    "testing"
)
func TestFakeFunc(t *testing.T) {
    var data string
    saved := realFunc //保存原来的realfunc
    defer func() { realFunc = saved }()
    realFunc = func(d string) {
        data = d + "demo"
    }
    CheckInfo("demo")
    if data != "demo" {
        t.Errorf(`data=%s want "demo"`, data)
    }
}
```

### 测试扩展包
一个测试包想要使用，调用了该包的函数，这样讲形成一个包循环，go是不允许的，因此可以在被调用包的目录创建一个`<path>/<package>_test`的目录。告诉go test工具应该建立一个额外的包运行测试。这样就可以在该包内倒入测试代码依赖的其他的辅助包

* go list 可以查看哪些是Go源文件产品代码，哪些是包测试，哪些是测试扩展包。
 * go list -f={ {.GoFiles} } fmt。GoFiles表示产品代码对应的Go源文件列表；也就是go build命令要编译的部分
 * go list -f={ {.TestGoFiles} } fmt fmt。TestGoFiles表示包内部测试代码，以_test.go为后缀文件名，不过只是在测试时被构建
 * go list -f={ {.XTestGoFiles} } fmt。XTestGoFiles表示属于测试扩展包的测试代码，也就是fmt_test
 
* 测试覆盖率
 * go test -coverprofile=c.out
  * coverprofile参数通过在测试代码中插入生成钩子来统计覆盖率数据。也就是说在运行每个测试前。他会修改测试代码的副本，在每个词法块都会设置一个布尔标志变量。当被修改后的的被测试代码运行退出时，讲统计日志数据写入到c.out文件，并打印一部分执行的语句的一个总结。如果你需要的是摘要使用go test -cover
  * -covermod=count可以在每个代码块插入计数器。统计每个块的执行次数，从而可以知道执行的热点代码
  * go tool cover -html=c.out可以生成一个HTML报告


## 基准测试
* Benchmark开头，参数是b *testing.B,其中b有一个变量N，用于指定操作循环的次数。系统自己指定,会自己调整
```
import "testing"
func BenchmarkIsPalindrome(b *testing.B) {
    for i := 0; i < b.N; i++ {
        IsPalindrome("A man, a plan, a canal: Panama")
    }
}
```
循环放在测试函数内部，而不是放在测试框架内实现，这样可以让每个基准测试函数有机会在循环启动前执行初始化代码，这样并不会显著影响每次迭代的平均运行时间

* go test -bench=. 需要通过bench标志指定要运行的基准测试函数。支持正则，默认是空。`.`表示匹配所有基准测试函数
* -benchmem命令行标志参数可以在报告中包含内存的分配数据统计
* 可以通过go tool pprof对操作系统信息进行采样，比如cpu,内存占用.可以参考[Profiling Go Programs](https://blog.golang.org/profiling-go-programs)

## 示例函数
* 以Example开头，没有参数
```
func ExampleIsPalindrome() {
    fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
    fmt.Println(IsPalindrome("palindrome"))
    // Output:
    // true
    // false
}
```
* 作为文档，godoc会根据示例函数的后缀名部分，将一个示例函数关联到某个具体函数或包本身
 * go test会运行示例函数，如果含有上面Output的注释，那么测试工具会比较输出结果和注释是否一致
 * 提供一个真实的演练场，像[http://golang.org](http://golang.org)一样就是有godoc提供的文档服务，他使用了Go Playground提高的技术让用户可以在浏览器中在线编辑和运行每个示例函数
