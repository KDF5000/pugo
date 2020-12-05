```toml
title = "PHP常见的五种设计模式——工厂模式"
date = "2015-08-08 18:53:37"
update_date = "2015-08-08 18:53:37"
author = "KDF5000"
thumb = ""
tags = ["PHP", "设计模式"]
draft = false
```
一直对设计模式有一种敬畏之心，每次想要看设计模式的时候就会想到`Erich Gamma`，` Richard Helm` ， `Ralph Johnson`， `John Vlissides `的黑皮`《设计模式》`，基本都望而止步，要把那本书看完可不是一时半会的，而且在没有项目经验的情况下，个人感觉基本都是纸上谈兵。

今天在`IBM Developerworks`上看到一篇文章将`PHP`中常用的五种设计模式，感觉还不错，而且只有五种五种五种（重要的强调三遍）！先从简单的入手，把这五种消灭了再说。以后慢慢学习其他的设计模式。

####工厂模式（`Factory Pattern`）
工厂这个词的使用是非常形象，字面意思可以这样认为，这种模式下，我们有一个工厂，这个工厂生产很多一种或者几种产品（其实多种的情况是覆盖了一种的），但是每个产品怎么生产和包装的我们不知道，其实我们也不需要知道，知道的越多你就越迷糊，以后你的行为就受制于太多杂事，也就是我们常说的耦合度太高，因此我们就将所有的事情交给工厂负责，我们只用告诉工厂需要什么，工厂把产品交付给你就是了。一旦产品的工艺发生改变，工厂负责就好，你使用该产品的工艺不受影响。因此工厂模式可以大大的降低系统的耦合度，增强系统的稳定性，当然也会提高代码的复用率。

在实际的程序设计中，工厂相当于一个对外的接口，那么这个接口的返回类型是确定的，那么我们怎么通过这个工厂来生产不同的产品发回给客户呢？很简单，做一个所有产品的“模子”就可以，这个“模子”有每个产品的所有特征，但是不能用，需要具体的产品实现这些特性，就是我们常说的`Interface`。
使用类图表示如下：
![](@media/archive/img_Factory_Pattern.png)
#####`PHP`的实现
* 编写一个接口 `Product.php`  

```
<?php
/**
* Created by PhpStorm.
* User: Defei
* Date: 2015/8/8
* Time: 16:14
*/
 
interface Product{
    public function getName();
}
```

<!--more-->

*  设计一个产品`A`实现`Product`接口

```
<?php
/**
* Created by PhpStorm.
* User: Defei
* Date: 2015/8/8
* Time: 16:16
*/

class ProductA implements Product{
 
    public function getName(){
        // TODO: Implement getName() method.
        echo '我是产品A';
    }
}
```

*   设计产品`B`实现`Product`接口

```
<?php
/**
* Created by PhpStorm.
* User: Defei
* Date: 2015/8/8
* Time: 16:17
*/
 
class ProductB implements Product{
 
    public function getName(){
        // TODO: Implement getName() method.
        echo '我是产品B';
    }
 
}
```

* 建造一座工厂生产产品`A`和`B`

```
<?php
/**
* Created by PhpStorm.
* User: Defei
* Date: 2015/8/8
* Time: 16:18
*/
class ProductFactory{
 
    /**
     * @param $product_name
     * @return mixed
     */
    public function factory($product_name){
        return new $product_name; //PHP可以使用名字直接new一个同名的对象这个很方便
    }
 
}
```

#####测试
产品`A`和`B`已经设计好了，工厂也建好了，下一步就是测试一下这个工厂对的生产能力如何。
```
<?php
/**
* Created by PhpStorm.
* User: Defei
* Date: 2015/8/8
* Time: 16:20
*/
include 'ProductFactory.php';
include 'Product.php';
include 'ProductA.php';
include 'ProductB.php';
 
$factory = new  ProductFactory();
echo $factory->factory('ProductA')->getName().PHP_EOL;
echo $factory->factory('ProductB')->getName();
 
```
输出结果如下：
![](@media/archive/img_factory_pattern.png)

