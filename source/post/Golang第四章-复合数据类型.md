```toml
title = "Golang第四章:复合数据类型"
date = "2017-02-25 12:00:21"
update_date = "2017-02-25 12:00:21"
author = "KDF5000"
thumb = ""
tags = ["Golang", "学习笔记"]
draft = false
```
主要介绍数组，slice, map和结构体
数组和结构体是聚合类型；数组有同构的元素组成，结构体有异构的元素组成，他们都是有固定内存大小的数据结构。
slice和map是动态的的数据结构，他们根据需要动态增长。

<!--more-->

### 数组
初始化
```go
var q [3]int = [3]int{1,2,3}
var q [3]int = [3]int{1,2}
q := [...]int{1,2,3} //...表示数组的长度根据初始化值得个数来计算
q := [...]int{1:1, 4:4,5:5} //提供索引和对应值的方式，长度则根据最大的索引值确定，这里是6，因此数组的长度是7
```

如果数组的元素类型是可以直接可以进行比较的，那么数组本身也是可以直接比较的，通过==进行比较，只有当两个数组的所有元素都是相等的时候数组才是相等的。

数组作为函数参数时，是值传递，会复制一份到参数里(C语言则是传递的引用)，因此在函数内部的修改不会反映在原始数组中。因此由于这种机制传递大的数组的效率是非常低的。但是可以将函数的参数显示的设置为数组指针，这样对函数内部数组的修改就可以反映在原始数组里了。

数组作为参数要指定长度，如果想要修改数组的值可以传递数组的指针
```go
func zero(ptr *[32]byte) {
    for i := range ptr {
        ptr[i] = 0 }
    }
}

func zero(ptr *[32]byte) {
    *ptr = [32]byte{}
}
//

```

### Slice
Slice(切片)代表变长的序列，序列中每个元素都有相同的类型。 
![80066410.png](@media/archive/blog/images/6F89E07747579464A98601695116E687.png)
数组和Slice的关系：一个slice是一个轻量级的数据结构，提供了访问**数组**子序列或者全部元素的功能，而且slice的底层确实引用了一个数组对象。一个slice有三个部分组成：指针，长度和容量。指针指向第一个slice元素对应的底层数组元素的地址，这个第一个元素并不一定是数组的第一个元素。长度对应slice元素的数目，长度不能超过容量，容量一般是从slice的开始位置到底层数据的结尾位置。内置的len和cap函数分别返回slice的长度和容量. 如果切片的操作超出了cap(s)的上限则会导致一个panic异常，但是超出len(s)则意味着扩展了slice，因为新的slice的长度会变大。
slice的初始化
```
s := []int{0,1,2,3,4,5,6} //slice
s := [...]int[0,1,2,3,4,5,6] //数组
```
好像和数组没很大的差别，但是仔细对比数组的初始化长度是固定的要么指定数组的长度要么...，使用初始化的参数个数。
slice不需要指明序列的长度，这会隐式的创建一个合适大小的数组，然后slice的指针指向底层的数组
```
a1 := []int{1, 2, 3, 4, 5}
a2 := [...]int{1, 2, 3, 4, 5} //数组
fmt.Println(cap(a1)) //5
fmt.Println(cap(a2)) //5
b1 := []int{}
b2 := [...]int{}
fmt.Println(cap(b1)) //0
fmt.Println(cap(b2)) //0
```

因为slice值包含指向第一个slice元素的指针，因此向函数传递一个slice将允许在函数内部修改底层数组的元素，换句话说，，复制一个slice只是对底层的数组创建了一个新的slice别名

```go
// reverse reverses a slice of ints in place.
func reverse(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}
a := [...]int{0, 1, 2, 3, 4, 5}
reverse(a[:])
fmt.Println(a) // "[5 4 3 2 1 0]"
```

slice之间不能比较，不能像数组那种==进行比较，因此必须我们自己遍历slice实现比较的功能，对于byte类型的slice，标准库提供了高度优化的bytes.Equal的比较方法。
原因见：

slice唯一合法的比较操作是和nil的比较，一个零值得slice等于nil.但是长度为0不一定等于nil, 比如：make([]int,3)[3:], 因此判断slice是否为空不能用slice == nil，要使用len(s) == 0.

使用make创建一个slice
```
make([]T, len)
make([]T, len, cap) // same as make([]T, cap)[:len]
```
在底层，make创建了一个匿名数组，然后返回一个slice.只有通过slice才能访问底层的匿名数组

```go
var s []int //slice
var s [3]int //数组
```

append函数可以扩充slice
```go
slice1 = append(slice1, x)
slice1 = append(slice1, slice2)
```
append的返回值是一个slice,因此必须将原来的slice重新赋值，在这个过程会产生内存的重新分配，下面是一个简单的模拟append的操作
```go
func appendInt(x []int, y int) []int {
    var z []int
    zlen := len(x) + 1
    if zlen <= cap(x) {
        // There is room to grow.  Extend the slice.
        z = x[:zlen]
    } else {
        // There is insufficient space.  Allocate a new array.
        // Grow by doubling, for amortized linear complexity.
        zcap := zlen
        if zcap < 2*len(x) {
            zcap = 2 * len(x)
        }
        z = make([]int, zlen, zcap)
        copy(z, x) // a built-in function; see text
    }
    z[len(x)] = y
    return z
}

func main() {
    // var nums []int //slice 长度为0
    // fmt.Println(nums)
    // nums = appendInt(nums, 1)
    // nums = appendInt(nums, 2)
    // nums = appendInt(nums, 4) //cap=4
    // nums = appendInt(nums, 4) //cap=4
    // fmt.Println(nums, cap(nums))
    var x, y []int
    for i := 0; i < 10; i++ {
        y = appendInt(x, i)
        fmt.Printf("%d cap=%d\t%v\n", i, cap(y), y)
        x = y
    }
}

///
/*
0 cap=1 [0]
1 cap=2 [0 1]
2 cap=4 [0 1 2]
3 cap=4 [0 1 2 3]
4 cap=8 [0 1 2 3 4]
5 cap=8 [0 1 2 3 4 5]
6 cap=8 [0 1 2 3 4 5 6]
7 cap=8 [0 1 2 3 4 5 6 7]
8 cap=16    [0 1 2 3 4 5 6 7 8]
9 cap=16    [0 1 2 3 4 5 6 7 8 9]
*/
```
slice可以用来实现stack


### Map
一个map就是一个哈希表的**引用**
```go
args := map[string]int {
    "alice": 12,
    "tom":21, //最好加上逗号，避免编译器在这里添加“;”
}
//2.make创建
args := make(map[string]int)
args["tom"] = 21
//删除
delete(args, "tom")
//遍历
for k,v := range args{
    fmt.Printf("%s\t%d\n", k, v)
}
```

* 如果一个key值不存在，则访问会返回value类型的零值
* map的遍历是随机的，每次的顺序都不同，这是故意的，每次使用随机的遍历顺序可以强制要求程序不会依赖具体的哈希函数实现。如果要按照顺序对key/value进行遍历，必须显式的对key进行排序，可以使用sort包的Strings函数对字符串slice排序
```go
import "sort"
var names []string
//先取出所有Key
for name := range args {
    names = append(names, name)
}
//排序
sort.Strings(names)
//有序访问
for _, name := range names {
    fmt.Printf("%s\t%d\n", name, ages[name])
}
```
* map的零值是nil,也就是没有指向任何的hash表
* map声明后如果不进行make创建让其指向一个hash表，则为零值，但是依然可以进行查找，删除，len和range操作，不过不能进行访问和存储元素，这样会造成panic
```go
var args map[string]int  //nil
fmt.Println(len(args)) //0
args["Tom"] = 32 //panic: assignment to entry in nil map
args = make(map[string]int)
args["Tom"] = 32 //Ok
```
* 使用下面的代码测试key是否存在
```go
age, ok := ages["bob"]
if !ok {
    fmt.Println("No key:  Bob,", age) //age=0
}
```
* 两个map不能进行比较，只能和nil进行比较
```go
func equal(x, y map[string]int) bool {
    if len(x) != len(y) {
        return false
    }
    for k, xv := range x {
        if yv, ok := y[k]; !ok || yv != xv { //注意这里要通过ok和相等两个条件判断
            return false
        }
    }
    return true
}
//equal(map[string]int{"A": 0}, map[string]int{"B": 42})
```

* map的vaule也可以是聚合类型，比如slice或者map，下面的graph可以表示一个图
```go
var graph = make(map[string]map[string]bool)
func addEdge(from, to string) {
    edges := graph[from]
    if edges == nil {
        edges = make(map[string]bool)
        graph[from] = edges
    }
    edges[to] = true
}
func hasEdge(from, to string) bool {
    return graph[from][to]
}
```

### 结构体
结构体就是你种类型变量的集合
```go
type  Employee struct {
    ID        int
    Name      string
    Address   string
    DoB       time.Time
    Position  string
    Salary    int
    ManagerID int
}
var dilbert Employee
```
* 结构体变量的的成员的访问可以通过点操作，与C不同的是，对于一个指针类型的结构体变量同样可以通过点操作访问，下面两个是等价的
```go
var employeeOfTheMonth *Employee = &dilbert
employeeOfTheMonth.Position += " (proactive team player)"
*employeeOfTheMonth).Position += " (proactive team player)"
```
* 结构体成员名字是以大写字母开头，那么该成员就是导出的，这是GO语言导出规则决定的
* 结构体成员不能是该结构体类型，但是可以是该类型的指针

**重要:**
> 如果指针作为函数的参数，如果传递到函数的指针在函数调用前已经初始化指向了一个数据，那么在函数内部通过指针修改变量，在函数外部是可见的，但是如果指针是在函数内部初始化，时间上传递的指针也是值传递，内部看到的变量只是外部变量的一个复制，所以外部变量的指针仍然为空

* 结构体赋值可以按照变量的顺序赋值，也可以k-v的形式，但是在外部包中不能使用顺序赋值的方式对结构体中未导出的变量赋值
```go
type Point struct{ X, Y,z int }
p := Point{1, 2, 3}
p := Point{X:1,Y:2,z:3}
//外部包
p := Point{1, 2, 3} //编译错误，compile error: can't reference z
```
* 结构体可以作为函数的参数和返回值，但是考虑到效率通常传入结构体的指针
* 如果结构体的全部成员是可以比较的，那么结构体也是可以比较的
* 可比较的结构体可以作为map的key
* 匿名类型指的是一个结构体可以包含一个结构体类型的成员变量，并且可以省略成员变量名，实际上改类型的成员会将该类型名作为变量名，在进行访问的时候可以直接访问匿名成员的成员变量，但是赋值的时候必须使用嵌套的方式进行复制
```go
type Point struct {
    X, Y int
}
type Circle struct {
    Center Point
    Radius int 
}
type Wheel struct {
    Circle Circle
    Spokes int 
}
//可以这么访问
var w Wheel
w.X = 8         // equivalent to w.Circle.Point.X = 8
w.Y = 8         // equivalent to w.Circle.Point.Y = 8
w.Radius = 5    // equivalent to w.Circle.Radius = 5
w.Spokes = 20
//不能这么赋值
w = Wheel{8, 8, 5, 20}                       // compile error: unknown fields
w = Wheel{X: 8, Y: 8, Radius: 5, Spokes: 20} // compile error: unknown fields
//这样赋值
w = Wheel{Circle{Point{8, 8}, 5}, 20}
w = Wheel{
    Circle: Circle{
        Point:  Point{X: 8, Y: 8},
        Radius: 5,
    },
    Spokes: 20, // NOTE: trailing comma necessary here (and at Radius)
}
```
* 匿名成员有个隐式的名字，因此不能同时包含两个类型相同的匿名成员，这会导致名字重涂
* 匿名成员同样有可见性的规则约束(首字母大写)
**重要:**
> 结构体的匿名成员并不要求是结构体，任何类型都可以作为匿名成员，但是为什么要嵌入一个没有任何子成员类型的匿名成员类型呢？简短的点运算语法可以用于选择匿名成员的嵌套成员，也可以访问他们的方法。实际上，外层的结构体不仅仅获得了匿名成员类型的所有成员，而且也获得了该类型导出的全部的方法。这个机制可以将一个简单行为的队形组合成复杂行为的对象。组合是Go语言中面向对象编程的核心。

### JSON
* 使用encoding/json进行json的编码解码，可以直接编码struct结构体
```go
type Movie struct {
    Title  string
    Year   int  `json:"released"`
    Color  bool `json:"color,omitempty"`
    Actors []string
}
var movies = []Movie{
    {Title: "Casablanca", Year: 1942, Color: false,
        Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
    {Title: "Cool Hand Luke", Year: 1967, Color: true,
        Actors: []string{"Paul Newman"}},
    {Title: "Bullitt", Year: 1968, Color: true,
        Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
    // ...
}
data, err := json.Marshal(movies)
if err != nil {
    log.Fatalf("JSON marshaling failed: %s", err)
}
///
[{"Title":"Casablanca","released":1942,"Actors":["Humphrey Bogart","Ingr id Bergman"]},{"Title":"Cool Hand Luke","released":1967,"color":true,"Ac tors":["Paul Newman"]},{"Title":"Bullitt","released":1968,"color":true," Actors":["Steve McQueen","Jacqueline Bisset"]}]

//解码
//解码的时候可以选择性的进行解码，比如下面只选择电影的名字解码
var titles []struct{ Title string }
if err := json.Unmarshal(data, &titles); err != nil {
    log.Fatalf("JSON unmarshaling failed: %s", err)
}
fmt.Println(titles) // "[{Casablanca} {Cool Hand Luke} {Bullitt}]"

```
* 编码时默认使用结构体成员名字作为Json的key值，但是只有导出的成员(大写)才会被编码
* Tag: 一个结构体成员Tag是和在编译节点关联到该成员的元信息字符串
```go
type Movie struct {
    Title  string
    Year   int  `json:"released"`
    Color  bool `json:"color,omitempty"`
    Actors []string
}
```
结构体成员Tag可以是任意的字符串面值，但是通常是一系列用空格隔开的key:"values"键值对序列；因为值中含有双引号字符，因此Tag一般用原生字符串面值形式书写。json头健名对应的值用于控制encoding/json包的解码和解码行为，并且encoding/..下面其他的包也遵循这个约定。成员Tag中json对应值得第一个部分指定json对象的名字，比如上面讲Year指定为released，Color指定为color,Color的Tag还带了额外的一个omitempty选项，便是当Go语言结构体成员为空或者零值得时候不生成json对象，也就是Color为false时，json中就没有color这个成员

### 文本和HTML模板
后期再看，更加适合web编程
* 有text/template和htmp/template提供一个将变量值填充到一个文本或者HTML格式的模板
* {{action}}是模板中的替换标志，action可以使一个变量也可以使表达式
