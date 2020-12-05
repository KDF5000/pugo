```toml
title = "Golang第六章:方法"
date = "2017-03-15 20:58:12"
update_date = "2017-03-15 20:58:12"
author = "KDF5000"
thumb = ""
tags = ["Golang", "学习笔记"]
draft = false
```
* 对象其实就是属性和方法的集合，可以认为方法就是意义和特殊类型关联的函数
* 在函数声明时，在其名字之前放上一个变量，即是一个方法。这个附加的参数会将该函数附加到这种类型上，即相当于这种类型定义了一个独占的方法。这个附加的参数叫接收器(Receiver)
* 方法名不能和对象的字段名一样，编译会报错
* **除了常见的struct类型可以定义方法，slice类型也可以，其实任何类型都可以。这个和很多其他面向对象的语言不同**，看个例子
```Go
package main
import "math"
import "fmt"
type Point struct{ X, Y float64 }
// traditional function
// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}
type Path []Point
func (p Path) Distance() float64 {
    sum := 0.0
    for i := range p {
        if i > 0 {
            sum += p[i-1].Distance(p[i])
        }
    }
    return sum
}
func main() {
    perim := Path{
        {1, 1},
        {5, 1},
        {5, 4},
        {1, 1},
    }
    fmt.Println(perim.Distance()) // "12"
}
```

<!--more-->

* 方法相比函数有一个明显的好处就是可以很简短，尤其是在外部包使用的时候，方法可以省略包名，直接通过对象变量调用
```Go
import "gopl.io/ch6/geometry"
perim := geometry.Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
fmt.Println(geometry.PathDistance(perim)) // "12", standalone function
fmt.Println(perim.Distance())             // "12", method of geometry.Path
```

### 基于指针对象的方法
如果一个函数需要更新一个变量，或者函数的其中一个参数实在太大我们希望能够避免进行这种默认的拷贝，这种情况下就需要用到指针了。对应到我们用来更新接收器的对象的方法，当这个接受者变量本身较大，我们就可以用其指针而不是用对象来声明方法。
```go
func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}
```
* 一般约定如果某个类有一个指针作为接收器的方法，那么所有其他的方法都必须有一个指针接收器。
* 只有类型和指向他们的指针，才能做接收器。
* 如果一个类型名本身就是一个指针的话，是不允许出现在接收器中的。
```go
type P *int
func (P) f() { /* ... */ } // compile error: invalid receiver type
```
* 要想调用指针类型的方法，只要提供一个该类型的指针即可
```go
r := &Point{1, 2}
r.ScaleBy(2)
fmt.Println(*r) // "{2, 4}"
//
p := Point{1, 2}
pptr := &p
pptr.ScaleBy(2)
fmt.Println(p) // "{2, 4}"
//
p := Point{1, 2}
(&p).ScaleBy(2)
fmt.Println(p) // "{2, 4}"
```
* 但是go语言也是允许直接使用该类型的变量调用方法,编译器会隐式的帮我们用&p去调用, 但是这种简写也是有限制的,这种方法只适用于变量，也就是说能够保证编译器能够取到地址。
```go
p := Point{1, 2}
p.ScaleBy(2)
fmt.Println(p) // "{2, 4}"
//
 Point{1, 2}.ScaleBy(2) // compile error: can't take address of Point literal
```
* 在每一个合法的方法调用表达式中，下面三种任意一种都是可以的
 - 接收器的实际参数和其接受器的形式参数相同，比如两者都是类型T或者类型*T
 - 接收器的形参是T，实参是*T, 编译器会隐式的位我们解引用，取到指针指向的实际变量
 - 接收器的形参是*T, 实参是T，编译器会隐式的取变量的地址
 
* 如果类型T的所有方法都是类型T，而不是*T作为接收器，那么拷贝这种类型的实例就是安全的，调用他的任何一个方法也就会产生一个值得拷贝，如果是*T就要避免对其进行拷贝，这样可能会破坏该类型内部的不变性
* 总结
 - 不管你的方法的接收器是指针类型还是非指针类型，都是可以通过指针/非指针类型进行调用，编译器会帮你做类型转换
 - 在声明一个方法的接收器是指针还是非指针类型时， 需要考虑：第一，这个对象是不是特别大，如果声明为非知怎类型，调用就会产生一次拷贝，；第二，如果选择指针类型作为接收器，那么这种类型时钟执行同一块内存地址，拷贝只是一个别名
 
### nil作为接收器
对类型是map, slice这样的nil也可以作为接收器,比如标准库`net/url`
```go
// Values maps a string key to a list of values.
type Values map[string][]string
// Get returns the first value associated with the given key,
// or "" if there are none.
func (v Values) Get(key string) string {
    if vs := v[key]; len(vs) > 0 {
        return vs[0]
    }
    return "" 
}
// Add adds the value to key.
// It appends to any existing values associated with key.
func (v Values) Add(key, value string) {
    v[key] = append(v[key], value)
}
```
 
### 通过嵌入结构体来扩展类型
* 一个结构体嵌套一个结构体，是对结构体的扩展，可以认为结构体拥有了嵌入结构体的属性和方法，可以当成自己的成员一样直接调用但是必须是作为匿名字段嵌入的类型，否则必须制定变量名调用。通过这种内嵌可以使我们定义字段特别多的复杂类型，我们可以将字段先按小类型分组，然后定义小类型的方法，之后再把他们组合起来。**很像继承**
```go
import "image/color"
type Point struct{ X, Y float64 }
type ColoredPoint struct {
    Point //必须匿名，这样下面才可以直接想自己的成员一样调用
    Color color.RGBA
}
var cp ColoredPoint
cp.X = 1
fmt.Println(cp.Point.X) // "1"
cp.Point.Y = 2
fmt.Println(cp.Y) // "2"
//
ed := color.RGBA{255, 0, 0, 255}
blue := color.RGBA{0, 0, 255, 255}
var p = ColoredPoint{Point{1, 1}, red}
var q = ColoredPoint{Point{5, 4}, blue}
fmt.Println(p.Distance(q.Point)) // "5" //必须显示的知名Point
p.ScaleBy(2)
q.ScaleBy(2)
fmt.Println(p.Distance(q.Point)) // "10"
```
* 从实现的角度看，内嵌字段会指导编译器去生成额外的包装方法来委托已经声明好的方法
```go
func (p ColoredPoint) Distance(q Point) float64 {
    return p.Point.Distance(q)
}
func (p *ColoredPoint) ScaleBy(factor float64) {
    p.Point.ScaleBy(factor)
}
```
* 内嵌的字段可以是指针类型的，这种情况下字段和方法会被间接地引入到当前的类型中，访问通过该指针指向的对象。这样就可以在不同的**对象之间共享通用的数据结构**。
* 当解析器解析一个选择器方法时，，首先找直接定义在这个类型的方法，然后内嵌字段引入的方法，然后去找匿名的内嵌字段引入的方法，一直这么递归向下找，如果选择器有二义性的话就会报错，也就是不能再同一级里有两个同名的方法。
* 如果两个匿名字段含有相同的方法名，不管其参数类型是否相同都是不能编译通过的，编译器无法确定是哪个方法

* 匿名struct
```go
var cache = struct {
    sync.Mutex
    mapping map[string]string
}{
    mapping: make(map[string]string),
}
func Lookup(key string) string {
    cache.Lock() //因为Mutex是匿名的
    v := cache.mapping[key]
    cache.Unlock()
    return v
}
```

### 方法值和方法表达式
* 方法值是绑定了**特定对象的方法**
```go
p := Point{1, 2}
q := Point{4, 6}
distanceFromP := p.Distance
fmt.Println(distanceFromP(q)
```

* 方法表达式：当T表示一个类型，方法表达式可能会写作T.f或者(*T).f, 会返回一个函数值，这种函数会将其第一个参数用作接收器
```go
p := Point{1, 2}
q := Point{4, 6}
distance := Point.Distance   // method expression
fmt.Println(distance(p, q))  // "5"
fmt.Printf("%T\n", distance) // "func(Point, Point) float64"
//
scale := (*Point).ScaleBy
scale(&p, 2)
fmt.Println(p)            // "{2 4}"
fmt.Printf("%T\n", scale) // "func(*Point, float64)"
```

### Bit数组
使用bit数组数显的非负整数集合: [intset.go](https://github.com/KDF5000/gopl/blob/master/src/ch6/intset.go)

### 封装
* 这里的封装是指一个对象的变量或者方法对调用方是不可见的，也叫信息隐藏。
* 大写字母的标识符从定义他们的包中导出，小写不会导出
* 封装的好处
 - 调用方不能直接修改对象的变量值，其只需关注少量的语句并且只要弄懂少量变量的可能的值即可
 - 隐藏实现细节，可以防止调用方不能直接修改对象的具体实现，这样可以保值不破坏api的情况下灵活的改变实现方式
 - 阻止外部调用方对对象内部的值任意的修改
