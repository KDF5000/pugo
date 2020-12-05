```toml
title = "PHP常见的五种设计模式——命令链模式"
date = "2016-05-23 21:10:28"
update_date = "2016-05-23 21:10:28"
author = "KDF5000"
thumb = ""
tags = ["PHP", "设计模式"]
draft = false
```
####命令链模式
命令链模式同样是为了解决耦合的问题，该模式将命令处理者和发号命令方分开，发号命令方不用知道谁处理该命令，以及以什么样的方式处理，它向所有已经注册的命令处理者发送命令，收到成功信息即可。

具体实现的时候，需要定义一个统一的接口，然后命令处理者实现接口，然后命令发号者用一个list保存所有注册的命令处理者，每个命令处理者会处理特定的一个或者一些命令，当发令者发送命令时，遍历list向已经注册的命令处理者发送命令，收到一个正确处理的信息即说明改命令已经有人处理。

类图如下：
![Alt text](@media/archive/blog/image/commandchain.png)

<!--more-->

php实现如下：
```
<?php
/**
 * command chain
 * @author KDF5000 
 * @since 2016-5-24
 */

interface ICommand{
	public function onCommand($name, $args);
}

class CommandChain{

	private $_commands = array();

	public function addCommand($command){
		$this->_commands[] = $command;
	}

	public function runCommand($name, $args){
		foreach ($this->_commands as $cmd) {
			# code...
			if($cmd->onCommand($name, $args)){
				return;
			}
		}
	}
}

class UserCommand implements ICommand{

	public function onCommand($name, $args){
		if($name == "addUser"){
			echo __CLASS__." reponses to addUser command!".PHP_EOL;
			return true;
		}
		return false;
	}
}

class MailCommand implements ICommand{

	public function onCommand($name, $args){
		if($name == "mail"){
			echo __CLASS__." reponses to mail command!".PHP_EOL;
			return true;
		}
		return false;
	}
}


$commandChain = new CommandChain();
$commandChain->addCommand(new UserCommand());
$commandChain->addCommand(new MailCommand());
$commandChain->runCommand("mail", null);
$commandChain->runCommand("addUser", null);
?>
```

