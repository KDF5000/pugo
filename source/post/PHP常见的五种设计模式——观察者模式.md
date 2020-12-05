```toml
title = "PHP常见的五种设计模式——观察者模式"
date = "2016-05-22 18:59:32"
update_date = "2016-05-22 18:59:32"
author = "KDF5000"
thumb = ""
tags = ["PHP", "设计模式"]
draft = false
```
####观察者模式
观察者模式是解决组件(对象)之间紧耦合的一种方式，顾名思义，既然叫观察者模式，涉及到两个对象，一个观察者，一个被观察者。观察者模式想要实现的效果就是当被观察者发生改变的时候，要主动通知观察者，自己发生了改变，至于观察者得知消息后做什么操作，被观察者无从得知，也不需要知道。

要实现这种模式，被观察的对象，要保存需要通知的观察者，这里称之为注册。然后当被观察者发生某个需要通知观察者的改变的时候，遍历已经注册的观察者，通知他们使用函数调用的方式。具体实现需要定义两个接口，用于观察者和被观察者实现，类图如下：
![Alt text](@media/archive/blog/image/观察者模式.png)

<!--more-->

#####PHP 实现
* 定义一个IObserver, 用于观察者实现
```
<?php
interface IObserver{
	public function onChange($sender, $args);
}

interface IObservable{
	public function addObserver($observer);
}

class ObservedObj implements IObservable{

	private $_observer = array();

	public function addObserver($observer){
		$this->_observer []= $observer;
	}

	public function changeSomeThing($str){
		foreach ($this->_observer as $ob) {
			# code...
			$ob->onChange($this, $str);
		}
	}
}

class Observer implements IObserver{

	private $_name = null;

	function __construct($name){
		$this->_name = $name;
	}
	function onChange($sender, $args){
		echo $this->_name." received msg:changed ".$args.PHP_EOL;
	}
}

$observedOb = new ObservedObj();
$observedOb->addObserver(new Observer("ob1"));
$observedOb->addObserver(new Observer("ob2"));
$observedOb->changeSomeThing("hello");
?>
```
