```toml
title = "使用svn和markdown管理项目文档"
date = "2015-07-31 22:07:38"
update_date = "2015-07-31 22:07:38"
author = "KDF5000"
thumb = ""
tags = ["PHP", "项目管理", "SVN", "Markdown"]
draft = false
```
文档管理，是每个项目非常重要又非常让人头疼的问题，每份文档可能会有不同的版本，而且可能还有有不同的人来共同编辑。除此还要考虑怎么才能方便的分享给项目组其他成员查看。也正是这些原因很多笔记记录软件推出了协同办公的功能，比如印象笔记，为知笔记等。但是把东西放在他们那里似乎不是很安全呵，毕竟是项目文档，都是比较私密的东西。也有一些团队协作平台提供了这个功能，不过他们都比较简单，只是上传下载文件。这两种方式似乎都不是很令人满意。

最近做项目，内部使用`svn`管理代码。就想是不是也可以使用svn来管理文档（当然有的公司确实是使用svn管理文档的），等等，这只是解决了文档的多人协作和版本控制的功能。那么怎么才能实现大家更方便的查看文档呢？借鉴现在比较流行的`Jekyll`，`hexo`搭建博客的思路，使用`markdown`写博客，然后生成`html`。因此借鉴这个思路，使用`svn`管理项目文档，项目文档尽可能使用`markdown`编写，然后写一个脚本定时解析`md`文件生成`html`，搭建一个服务器，将生成的`html`部署到服务器，这样同一个局域网的组员就可以访问这些文档了。

####搭建svn
`svn` 分为服务端和客户端，搭建详细教程参考[`windows下svn服务器的搭建`](http://www.jb51.net/article/29005.htm)
####解析`md`文件生成`html`
markdown本来就是一种轻量级的标记语言，有很多开源的解析markdown文件的项目。本文使用一个开源的`PHP`的开源项目[ParseDown](https://github.com/erusev/parsedown)，用来解析markdown生成html文件，在此项目的基础上，封装了一个类，实现解析指定目录下所有markdown文件，生成html文件（包括子目录下的markdown文档，按照原来的目录存放）。

下面是指定目录下所有的`md`文件生成`html`文件的源码：

<!--more-->

```
<?php
/**
* 解析指定目录下所有的md文件生成html文件
**/
require_once 'Parsedown.php';
class Parse2Html{
	private $files; //指定目录下的所有文件

	function __construct($dir,$extension='md'){
		if($dir==NULL ||!is_dir($dir)){
			echo '请指定一个合法目录';
			exit;
		}
		$this->files = self::getAllFiles($dir);
	}

	/**
	 * 解析目录下面的所有md文件输出到html文件
	 */
	public function parse2Html(){
		foreach ($this->files as $file ) {
			$down_text = file_get_contents($file);
			$file_name = pathinfo($file,PATHINFO_FILENAME);
			$title = iconv('gb2312','utf-8',$file_name);
			$html_content = $this->html($down_text,$title);

			file_put_contents(dirname($file).DIRECTORY_SEPARATOR.$file_name.'.html',$html_content);
			echo $file_name.'export successfully!'.PHP_EOL;
		}
	}

	/**
	 * 将markdown内容转换为html
	 * @param $down_text
	 * @param $title
	 * @return string
	 */
	private function html($down_text, $title){
		$parsedown = Parsedown::instance();
		$parsedown->setBreaksEnabled(true);
		$html_text = $parsedown->text($down_text);
		$html_content = <<<CONTAINER
<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv='Content-Type' content='text/html; charset=utf-8' />
	<meta charset="UTF-8">
	<title>$title</title>
</head>
<body>
	$html_text
</body>
</html>
CONTAINER;
		return $html_content;
	}
	/**
	 * 获取指定目录下的所有文件，返回一个数组，数组元素为文件的实际路径
	 * @param $dir
	 * @return array
	 */
	private static function getAllFiles($dir,$extension='md'){
		$files = array();
		if(is_file($dir)){
			if(pathinfo($dir,PATHINFO_EXTENSION) == $extension){
				array_push($files,$dir);
			}
			return $files;
		}
		$folder = new DirectoryIterator($dir);
		foreach($folder as $file){
			if($file->isFile()){
				if(pathinfo($file,PATHINFO_EXTENSION) == $extension){
					array_push($files,$file->getRealPath());
				}
				continue ;
			}
			if(!$file->isDot()){
				$sub_files = self::getAllFiles($file->getRealPath(),$extension);
				$files = array_merge($files,$sub_files);
			}
		}
		return $files;
	}
}
?>
```
 调用Parse2Html类的方法如下(index.php):
```
<?php
require_once 'Parse2Html.php';
$dir = dirname(__FILE__).DIRECTORY_SEPARATOR.'Doc';
$parse2Html = new Parse2Html($dir);
$parse2Html->parse2Html();
?>
```
####`svn`建立文档仓库
新建一个`svn`仓库，所有编写的`markdown`文档都放进去，然后**将在上一步的`index.php`里指定该仓库的目录为解析的目录**。然后写一个`bat`文件，调用`index.php`文件生成`html`文件，内容如下:
```
php F:\Boss\Doc\index.php  ；假设已经将php可执行程序添加到环境变量
```
此时运行`bat`文件就可以解析`svn`目录下的所有`md`文件生成`html`文档，也可以将该`bat`文件添加到系统任务里，设定每天固定时间执行。
执行结果如下:
![](@media/archive/img_parse2htmll2.png)
其中`产品概述.html`是由`产品概述.md`文档生成的
####部署到`apache`服务器
这里使用的是`Apache`服务器，当然也可以是其他服务器，比如`IIS`，只要可以解析html文件就可以。将网站的目录设置为`svn`文档仓库的根目录，在根目录下编写一个`index.php`文件，用于罗列该目录下所有想要显示的文件。内容如下
```
<?php
header("content-type:text/html;charset=gb2312");
/**
 * 获取指定目录下的所有文件，返回一个数组，数组元素为文件的实际路径
 * @param $dir
 * @return array
 */
function getAllFiles($dir,$extension){
    $files = array();
    if(is_file($dir)){
        if(in_array(pathinfo($dir,PATHINFO_EXTENSION),$extension)){
            array_push($files,$dir);
        }
        return $files;
    }
    $folder = new DirectoryIterator($dir);
    foreach($folder as $file){
        if($file->isFile()){
            if(in_array(pathinfo($file,PATHINFO_EXTENSION) ,$extension) ){
                array_push($files,$file->getRealPath());
            }
            continue ;
        }
        if(!$file->isDot()){
            $sub_files = getAllFiles($file->getRealPath(),$extension);
            $files = array_merge($files,$sub_files);
        }
    }
    return $files;
}
$files = getAllFiles(dirname(__FILE__),array('html','doc','docx'));
$root_len = strlen(dirname(__FILE__));
foreach ($files as $file) {
    $name = substr($file,$root_len+1,strlen($file)-strlen(pathinfo($file,PATHINFO_EXTENSION))-$root_len-1-1);
    echo "<a href='".substr($file,strlen($_SERVER['DOCUMENT_ROOT']))."'/>".$name."</a></br>";
}
?>
```
此时访问大家好的网站，就可以看到如下结果:
![](@media/archive/img_parse2html_1.png)
![](@media/archive/img_parse2html_3.png)


