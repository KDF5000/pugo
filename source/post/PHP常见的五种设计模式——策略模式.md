```toml
title = "PHP常见的五种设计模式——策略模式"
date = "2016-05-23 23:11:04"
update_date = "2016-05-23 23:11:04"
author = "KDF5000"
thumb = ""
tags = ["PHP", "设计模式"]
draft = false
```

####策略模式
在此模式中，算法是从复杂类提取的，因而可以方便地替换。例如，如果要更改搜索引擎中排列页的方法，则策略模式是一个不错的选择。思考一下搜索引擎的几个部分 —— 一部分遍历页面，一部分对每页排列，另一部分基于排列的结果排序。在复杂的示例中，这些部分都在同一个类中。通过使用策略模式，您可将排列部分放入另一个类中，以便更改页排列的方式，而不影响搜索引擎的其余代码.

类图如下：
![Alt text](@media/archive/blog/image/strategy.png)

<!--more-->

php实现:
```
<?php
/**
 * strategy 
 * @author KDF5000
 * @since 2016-5-24 22:51
 */

interface IStrategy{
	public function filter($name);
}

class FindAfterStrategy implements IStrategy{

	private $_name;
	public function __construct($name){
		$this->_name = $name;
	}

	public function filter($name){
		return strcmp($this->_name, $name) <= 0;
	}

}

class RandomStrategy implements IStrategy{

	public function filter($name){
		return rand(0,1) >= 0.5;
	}
}

class UserList{

	private $_list = array();

	public function __construct($names){
		if($names != null){
			foreach ($names as $name) {
				# code...
				$this->_list[] = $name;
			}
		}
	}

	public function add($name){
		$this->_list[] = $name;
	}

	public function find($filter){
		$res = array();
		foreach ($this->_list as $name) {
			# code...
			if($filter->filter($name)){
				$res[] = $name;
			}
		}
		return $res;
	}
}

$user_list = new UserList(array('Jack','Tom', 'Mike', 'Devin'));
$res_1 = $user_list->find(new FindAfterStrategy("M"));
$res_2 = $user_list->find(new RandomStrategy());

var_dump($res_1);
var_dump($res_2);

?>
```

