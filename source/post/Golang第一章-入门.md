```toml
title = "Golang第一章:入门"
date = "2017-02-25 11:59:10"
update_date = "2017-02-25 11:59:10"
author = "KDF5000"
thumb = ""
tags = ["Golang", "学习笔记"]
draft = false
```
### os.Args获取参数
os.Args 是一个字符串的切片，它的第一个元素os.Args[0]是命令本身的名字，其他元素则是程序启动时传给它的参数。

i++ 给i加1，是语句，而在c语言中则是表达式，因此在golang中j = i++ 是非法的。

### for循环
```
for initialization; condition; post {
    // zero or more statements
}
// a traditional "while" loop
for condition {
    // ...
}
```

<!--more-->

### 字符串
string 类型可以看成一种特殊的slice, 因此可以使用len获取长度，同时支持切片操作，但是对于单个元素，如a[0]的结果是一个byte，输出来是asiic码，需要string(a[0])这样转换，但是可以通过切片操作获取子串，如a[2:]

#### 是否存在某子串，子串出现次数
```
//contains 和containsAny都是调用Index来判断子串是否出现在字符串中
    //空格隔开子串
    fmt.Println(strings.ContainsAny("hello", "s e")) //true
    fmt.Println(strings.ContainsAny("hello", ""))    //false
    fmt.Println(strings.ContainsAny("hello", "lo"))  //true

    //count也就是字符串匹配实现的是Rabin-Karp算法，Count 是计算子串在字符串中出现的无重叠的次
    fmt.Println(strings.Count("fivevev", "ve"))  //2
    fmt.Println(strings.Count("fivevev", ""))    //8 utf8.RuneCountInString(s) + 1 也就是长度+1
    fmt.Println(strings.Count("fivevev", "vev")) //1
```

#### 字符串的分割
六个三组函数：Fields 和 FieldsFunc、Split 和 SplitAfter、SplitN 和 SplitAfterN
##### Fileds和FiledsFunc
```
func Fields(s string) []string
func FieldsFunc(s string, f func(rune) bool) []string
```
Fields 用一个或多个连续的空格分隔字符串 s,返回子字符串的数组（slice）
由于是用空格分隔，因此结果中不会含有空格或空子字符串，例如：
```
fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))
```
输出：
```
Fields are: ["foo" "bar" "baz"]
```
FieldsFunc 用这样的Unicode代码点 c 进行分隔：满足 f(c) 返回 true。该函数返回[]string。如果字符串 s 中所有的代码点(unicode code points)都满足f(c)或者 s 是空，则 FieldsFunc 返回空slice。
也就是说，我们可以通过实现一个回调函数来指定分隔字符串 s 的字符。比如上面的例子，我们通过 FieldsFunc 来实现：
```
fmt.Println(strings.FieldsFunc("  foo bar  baz   ", unicode.IsSpace))
```
实际上，Fields 函数就是调用 FieldsFunc 实现的：
```
func Fields(s string) []string {
    return FieldsFunc(s, unicode.IsSpace)
}
```
##### Split 和 SplitAfter、 SplitN 和 SplitAfterN
之所以将这四个函数放在一起讲，是因为它们都是通过一个同一个内部函数来实现的。它们的函数签名及其实现：
```
func Split(s, sep string) []string { return genSplit(s, sep, 0, -1) }
func SplitAfter(s, sep string) []string { return genSplit(s, sep, len(sep), -1) }
func SplitN(s, sep string, n int) []string { return genSplit(s, sep, 0, n) }
func SplitAfterN(s, sep string, n int) []string { return genSplit(s, sep, len(sep), n) }
```

Split和SplitAfter的区别是，After会保留分隔符
```
//分割字符串
fmt.Printf("%q\n", strings.Split("foo,bar,baz", ","))      //["foo" "bar" "baz"]
fmt.Printf("%q\n", strings.SplitAfter("foo,bar,baz", ",")) //["foo," "bar," "baz"]
```
    带 N 的方法可以通过最后一个参数 n 控制返回的结果中的 slice 中的元素个数，当 n < 0 时，返回所有的子字符串；当 n == 0 时，返回的结果是 nil；当 n > 0 时，表示返回的 slice 中最多只有 n 个元素，其中，最后一个元素不会分割，比如：

```
fmt.Printf("%q\n", strings.SplitN("foo,bar,baz", ",", 2))
//["foot" "bar,baz"]
```

##### 字符串数组(或slice)的连接
```
func Join(a []string, sep string) string
```
标准库的是实现方式
```
func Join(a []string, sep string) string {
    if len(a) == 0 {
        return ""
    }
    if len(a) == 1 {
        return a[0]
    }
    n := len(sep) * (len(a) - 1)
    for i := 0; i < len(a); i++ {
        n += len(a[i])
    }
    b := make([]byte, n)
    bp := copy(b, a[0])
    for _, s := range a[1:] {
        bp += copy(b[bp:], sep)
        bp += copy(b[bp:], s)
    }
    return string(b)
}
```
标准库的实现没有用 bytes 包，当然也不会简单的通过 + 号连接字符串。Go 中是不允许循环依赖的，标准库中很多时候会出现代码拷贝，而不是引入某个包。这里 Join 的实现方式挺好，我个人猜测，不直接使用 bytes 包，也是不想依赖 bytes 包（其实 bytes 中的实现也是 copy 方式）。

##### 字符串替换
* 
```
// 用 new 替换 s 中的 old，一共替换 n 个。
// 如果 n < 0，则不限制替换次数，即全部替换
func Replace(s, old, new string, n int) string
```
只能替换一种string
* Replacer
这是一个结构，没有导出任何字段，实例化通过 func NewReplacer(oldnew ...string) *Replacer 函数进行，其中不定参数 oldnew 是 old-new 对，即进行多个替换。
解决上面说的替换一种的问题：
```
r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
fmt.Println(r.Replace("This is <b>HTML</b>!"))
```
另外，Replacer 还提供了另外一个方法：
func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)
它在替换之后将结果写入 io.Writer 中。
