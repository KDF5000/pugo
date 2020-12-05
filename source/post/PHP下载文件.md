```toml
title = "PHP下载文件"
date = "2015-07-24 10:59:48"
update_date = "2015-07-24 10:59:48"
author = "KDF5000"
thumb = ""
tags = ["PHP"]
draft = false
```
###php实现下载功能
```
<?php
$file = NAC_ROOT.'/ident/theme/'.$theme.'/images/logo_src.gif';
		$fileTmp = pathinfo($file);
		$fileExt = $fileTmp['extension'];
		$saveFileName = ($this->themes[$theme].'.'.$fileExt);
		$fp=fopen($file,"r");
		$file_size=filesize($file);
		
		//下载文件需要用到的头
		Header("Content-type: application/octet-stream"); 
		Header("Accept-Ranges: bytes"); 
		Header("Accept-Length:".$file_size); 
		Header("Content-Disposition: attachment; filename=".$saveFileName); 
		$buffer=1024;
		$file_count=0;
		//向浏览器返回数据
		while(!feof($fp) && $file_count<$file_size){
			$file_con=fread($fp,$buffer);
			$file_count+=$buffer;
			echo $file_con;
		}
		fclose($fp);
		exit;
 ?>
```
<!--more-->
