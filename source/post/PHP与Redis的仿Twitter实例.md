```toml
title = "PHP与Redis的仿Twitter实例"
date = "2015-07-24 11:02:27"
update_date = "2015-07-24 11:02:27"
author = "KDF5000"
thumb = ""
tags = ["PHP", "Redis"]
draft = false
```
译者注：
原文位于Redis官网http://redis.io/topics/twitter-clone
Redis是NoSQL数据库中一个知名数据库，在新浪微博中亦有部署，适合固定数据量的热数据的访问。
作为入门，这是一篇很好的教材，简单描述了如何使用KV数据库进行数据库的设计。新的项目www.xiayucha.com亦采用Redis + MySQL进行开发，考虑Redis文档比较少，故翻译了此文。
其他参考资料：
  
[Redis命令参考中文版(Redis Command Reference](http://redis.readthedocs.org/en/latest/index.html)
[Try Redis](http://try.redis-db.com/)

我会在此文中描述如何使用`PHP`以及仅使用`Redis`来设计实现一个简单的`Twitter`克隆。很多编程社区常认为 `KV` 储存是一个特别的数据库，在 `web` 应用中不能替代关系数据库。本文尝试证明这恰恰相反。这个`twitter`克隆名为`Retwis`，结构简单，性能优异，能很轻易地用`N``个web`服务器和`Redis`服务器以分布式架构。

在此获取源码[http://code.google.com/p/redis/downloads/list](http://code.google.com/p/redis/downloads/list)
我们使用PHP作为例子是因为它能被每个人读懂，也能使用Ruby、Python、Erlang或其他语言获取同样(或者更佳)的效果。

*注意*：`Retwis-RB`是一个由`Daniel Lucraft`用`Ruby`与`Sinatra`写的`Retwis`分支！此文全部代码在本页尾部的`Git repository`链接里。此文以PHP为例，但是`Ruby`程序员也能检出其他源码。他们很相似。注意`Retwis-J`是`Retwis`的一个分支，由`Costin Leau`以`Java`和`Spring`框架写成。源码能在`GitHub`找到，并且在`springsource.org`有综合的文档。

<!--more-->

###`Key-value` 数据库基础
`KV`数据的精髓，是能够把`value`储存在`key`里，此后该数据仅能够通过确切的`key`来获取，无法搜索一个值。确切的来讲，它更像一个大型`HASH/`字典，但它是持久化的，比如，当你的程序终止运行，数据不会消失。比如我们能用`SET`命令以`key foo` 来储存值` bar`
 > 
 SET foo bar
 
`Redis`会永久储存我们的数据，所以之后我们可以问`Redis`：“储存在`key foo`里的数据是什么？”，`Redis`会返回一个值：`bar`
 > 
 GET foo => bar
 
`KV`数据库提供的其他常见操作有:`DEL`，用于删除指定的`key`和关联的`value`；
`SET-if-not-exists` (在`Redis`上称为`SETNX `)仅会在`key`不存在的时候设置一个值；
`INCR`能够对指定的`key`里储存的数字进行自增。
> 
 SET foo 10
 INCR foo => 11
 INCR foo => 12
 INCR foo => 13
 
####原子操作
目前为止它是相当简单的，但是`INCR`有些不同。设想一下，为什么要提供这个操作？毕竟我们自己能用以下简单的命令实现这个功能：
 > 
 x = GET foo
 x = x + 1
 SET foo x
 
问题在于要使上面的操作正常进行，同时只能有一个客户端操作`x`的值。看看如果两台电脑同时操作这个值会发生什么：
> 
 x = GET foo (返回10)
 y = GET foo (返回10)
 x = x + 1 (x现在是11)
 y = y + 1 (y现在是11)
 SET foo x (foo现在是11)
 SET foo y (foo现在是11)
 
问题发生了！我们增加了值两次，本应该从`10`变成`12`，现在却停留在了`11`。这是因为用`GET`和`SET`来实现`INCR`不是一个原子操作(`atomic operation`)。所以`Redis\memcached`之类提供了一个原子的`INCR`命令，服务器会保护`get-increment-set`操作，以防止同时的操作。让`Redis`与众不同的是它提供了更多类似`INCR`的方案，用于解决模型复杂的问题。因此你可以不使用任何SQL数据库、仅用Redis写一个完整的web应用，而不至于抓狂。

####超越`Ke-Value`数据库
本节我们会看到构建一个Twitter克隆所需Redis的功能。首先需要知道的是，Redis的值不仅仅可以是字符串(String)。
Redis的值可以是列表(Lists)也可以是集合(Sets)，在操作更多类型的值时也是原子的，所以多方操作同一个KEY的值也是安全的。
让我们从一个Lists开始：
> 
 LPUSH mylist a (现在mylist含有一个元素:'a'的list)
 LPUSH mylist b (现在mylist含有元素'b,a')
 LPUSH mylist c (现在mylist含有'c,b,a')
 
LPUSH的意思是Left Push， 就是把一个元素加在列表(list)的左边(或者说头上)。在PUSH操作之前，如果mylist这个键(key)不存在，Redis会自动创建一个空的list。就像你能想到的一样，同样有RPUSH操作可以把元素加在列表(list)的右边(尾部)。这对我们复制一个twitter非常有用，例如我们可以把用户的更新储存在username:updates里。当然，我们也有相应的操作来获取数据或者信息。比如LRANGE返回列表(list)的一个范围内的元素，或者所有元素
> LRANGE mylist 0 1 => c,b

LRANGE使用从零开始的索引(zero-based indexes)，第一个元素的索引是0，第二个是1，以此类推。该命令的参数是：
> 
LRANGE key first-index last-index

参数last index可以是负数，具有特殊的意义：-1是列表(list)的最后一个元素，-2是倒数第二个，以此类推。所以，如果要获取整个list，我们能使用以下命令：
> 
 LRANGE mylist 0 -1 => c,b,a
 
其他重要的操作有LLEN，返回列表(list)的长度，LTRIM类似于LRANGE，但不仅仅会返回指定范围内的元素，而且还会原子地把列表(list)的值设置这个新的值。我们将会使用这些list操作，但是注意阅读Redis文档来浏览所有redis支持的list操作。

###数据类型：集合(set)
除了列表(list)，Redis还提供了集合(sets)的支持，是不排序(unsorted)的元素集合。
它能够添加、删除、检查元素是否存在，并且获取两个结合之间的交集。当然它也能请求获取集合（set）里一个或者多个元素。
几个例子可以使概念更为清晰。记住：SADD是往集合(set)里添元素；SREM是从集合(set)里删除元素；SISMEMBER是检测一个元素是否包含在集合里；SINTER用于显示两个集合的交集。
其他操作有，SCARD用于获取集合的基数(集合中元素的数量)；SMEMBERS返回集合中所有的元素
> 
 SADD myset a
 SADD myset b
 SADD myset foo
 SADD myset bar
 SCARD myset => 4
 SMEMBERS myset => bar,a,foo,b
 
注意SMEMBERS不会以我们添加的顺序返回元素，因为集合(Sets)是一个未排序的元素集合。如果你要储存顺序，最好使用列表(Lists)取而代之。以下是基于集合的一些操作：
> 
 SADD mynewset b
 SADD mynewset foo
 SADD mynewset hello
 SINTER myset mynewset => foo,b
 
SINTER能够返回集合之间的交集，但并不仅限于两个集合(Sets)，你能获取4个、5个甚至1000个集合(sets)的交集。最后，让我们看下SISMEMBER是如何工作的：
> 
 SISMEMBER myset foo => 1
 SISMEMBER myset notamember => 0
 
Okay，我觉得我们可以开始coding啦！

####先决条件
如果你还没下载，请前往[http://code.google.com/p/redis/downloads/list](http://code.google.com/p/redis/downloads/list)下载Retwis的源码。它包含几个PHP文件，是个简单的tar.gz文件。实现的非常简单，你会在里面找到PHP客户端(redis.php)，用于redis与PHP的交互。该库由Ludovico Magnocavallo(http://qix.it/ )编写，你可以在自己的项目中免费使用。但如果要更新库的版本请下载Redis的发行版。(注意:现在有更好的PHP库了，请检查我们的客户端页面[http://redis.io/clients](http://redis.io/clients>))你需要的另一个东西是正常运行的Redis服务器。仅需要获取源码、用make编译、用./redis-server就完工了，点儿也不须配置就可以在你的电脑上运行Retwis。
 
###数据结构规划
  当使用关系数据库的时候，这一步往往是在设计数据表、索引的表单里处理。我们没有表，那我们设计什么呢？ 我们需要确认物体使用的key以及key采用的类型。

   让我们从用户这块开始设计。当然了，首先需要展示用户的username, userid, password, followers，自己follow的用户等。第一个问题是：如何在我们的系统中标识一个用户？username是个好主意，因为它是唯一的。不过它太大了，我们想要降低内存的使用。如果我们的数据库是关系数据库，我们能关联唯一ID到每一个用户。每一个对用户的引用都通过ID来关联。做起来很简单，因为我们有我们的原子的INCR命令！当我们创建一个新用户，我们假设这个用户叫"antirez"：
> 
 INCR global:nextUserId => 1000
 SET uid:1000:username antirez
 SET uid:1000:password p1pp0
 
我们使用global:nextUserId为键(Key)是为了给每个新用户分配一个唯一ID，然后用这个唯一ID来加入其他key，以识别保存用户的其他数据。这就是kv数据库的设计模式!请牢记于心，除了已经定义的KEY，我们还需要更多的来完整定义一个用户，比如有时需要通过用户名来获取用户ID，所以我们也需要设置这么一个键(Key)
> 
 SET username:antirez:uid 1000
 
一开始看上去这样很奇怪，但请记住我们只能通过key来获取数据!这不可能告诉Redis返回包含某值的Key，这也是我们的强处。
用关系数据库方式来讲，这个新实例强迫我们组织数据，以便于仅使用primary key访问任何数据。

####关注\被关注与更新
这也是在我们系统中另一个重要需求.每个用户都有follower，也有follow的用户.对此我们有最佳的数据结构!那就是.....集合(Sets).那就让我们在结构中加入两个新字段:
> 
 uid:1000:followers => Set of uids of all the followers users
 uid:1000:following => Set of uids of all the following users
 
另一个重要的事情是我们需要有个地方来放用户主页上的更新。这个要以时间顺序排序，最新的排在旧的前面。所以，最佳的类型是列表(List)。
基本上每个更新都会被LPUSH到该用户的updates key.多亏了LRANGE，我们能够实现分页等功能。请注意更新(updates)和帖子(posts)讲的是同一个东西，实际上更新(updates)是有点小的帖子(posts)。
> 
 uid:1000:posts => a List of post ids, every new post is LPUSHed here.
 
####验证
OK，除了验证，或多或少我们已经有了关于该用户的一切东西。我们处理验证用一个简单而健壮(鲁棒)的办法:我们不使用PHP的session或者其他类似方式。
我们的系统必须是能够在不同不同服务器上分布式部署的，所以一切状态都必须保存在Redis里。所以我们所需要的一个保存在已验证用户cookie里的随机字符串。
包含同样随机字符串的一个key告诉我们用户的ID。我们需要使用两个key来保证这个验证机制的健壮性:
> 
 SET uid:1000:auth fea5e81ac8ca77622bed1c2132a021f9
 SET auth:fea5e81ac8ca77622bed1c2132a021f9 1000
 
为了验证一个用户，我们需要做一些简单的工作(login.php):
* 从登录表单获取用户的用户名和密码
* 检查是否存在一个键 username:<username>:uid
* 如果这个user id存在(假设1000)
* 检查 uid:1000:password 是否匹配，如果不匹配，显示错误信息
* 匹配则设置cookie为字符串"fea5e81ac8ca77622bed1c2132a021f9"(uid:1000:auth的值)
实例代码:

```php
    include("retwis.php");   
   
     # Form sanity checks   
    if (!gt("username") || !gt("password"))   
        goback("You need to enter both username and password to login.");   
      
    # The form is OK, check if the username is available   
    $username = gt("username");   
    $password = gt("password");   
    $r = redisLink();   
    $userid = $r->get("username:$username:id");   
    if (!$userid)   
     goback("Wrong username or password");   
    $realpassword = $r->get("uid:$userid:password");   
    if ($realpassword != $password)   
     goback("Wrong useranme or password");   
    
    # Username / password OK, set the cookie and redirect to index.php   
    $authsecret = $r->get("uid:$userid:auth");   
    setcookie("auth",$authsecret,time()+3600*24*365);   
    header("Location: index.php");   
```

每次用户登录都会运行，但我们需要一个函数isLoggedIn用于检验一个用户是否已经验证。
这些是isLoggedIn的逻辑步骤
* 从用户获取cookie里auth的值。如果没有cookie，该用户未登录。我们称这个cookie为<authcookie>
* 检查auth:<authcookie>是否存在，存在则获取值(例子里是1000)
* 为了再次确认，检查uid:1000:auth是否匹配
* 用户已验证，在全局变量$User中载入一点信息
也许代码比描述更短:

```php
 function isLoggedIn() {   
      global $User, $_COOKIE;   
    
      if (isset($User)) return true;   
    
      if (isset($_COOKIE['auth'])) {   
          $r = redisLink();   
          $authcookie = $_COOKIE['auth'];   
          if ($userid = $r->get("auth:$authcookie")) {   
               if ($r->get("uid:$userid:auth") != $authcookie) return false;   
               loadUserInfo($userid);   
               return true;   
           }   
       }   
       return false;   
   }   
     
   function loadUserInfo($userid) {   
       global $User;   
     
       $r = redisLink();   
       $User['id'] = $userid;   
       $User['username'] = $r->get("uid:$userid:username");   
       return true;   
   }   
```

把loadUserInfo作为一个独立函数对于我们的应用而言有点杀鸡用牛刀了，但是对于复杂的应用而言这是一个不错的模板。
作为一个完整的验证，还剩下logout还没实现。在logout的时候我们怎么做呢？
很简单，仅仅改变uid:1000:auth里的随机字符串，删除旧的auth:<oldauthstring>并增加一个新的auth:<newauthstring>
*重要* :logout过程解释了为什么我们不仅仅查找auth:<randomstring>而是再次检查了uid:1000:auth。真正的验证字符串是后者，auth:<randomstring>是易变的.
假设程序中有BUGs或者脚本被意外中断，那么就有可能有多个auth:<something>指向同一个用户id。
logout代码如下:(logout.php)

```php
	include("retwis.php");   
	  
	if (!isLoggedIn()) {   
	    header("Location: index.php");   
	    exit;   
	}   
	  
	$r = redisLink();   
	$newauthsecret = getrand();   
    $userid = $User['id'];   
    $oldauthsecret = $r->get("uid:$userid:auth");   
    
    $r->set("uid:$userid:auth",$newauthsecret);   
    $r->set("auth:$newauthsecret",$userid);   
    $r->delete("auth:$oldauthsecret");   
    
    header("Location: index.php");   
	   
```

以上是我们所描述过的，应该比较易于理解。
更新(Updates)
更新，或者称为帖子(posts)的实现则更为简单。为了在数据库里创建一个新的帖子，我们做了以下工作:
 INCR global:nextPostId => 10343
 SET post:10343 "$owner_id|$time|I'm having fun with Retwis"
就像你看到的一样，帖子的用户id和时间直接储存在了字符串里。
在这个例子中我们不需要根据时间或者用户id来查找帖子，所以把他们紧凑地挤在一个post字符串里更佳。
在新建一个帖子之后，我们获得了帖子的id。需要LPUSH这个帖子的id到每一个follow了作者的用户里去，当然还有作者的帖子列表。
update.php这个文件展示了这个工作是如何完成的:

```PHP
    include("retwis.php");   
	  
	if (!isLoggedIn() || !gt("status")) {   
	    header("Location:index.php");   
	    exit;   
	}   
	  
	$r = redisLink();   
	$postid = $r->incr("global:nextPostId");   
    $status = str_replace("\n"," ",gt("status"));   
    $post = $User['id']."|".time()."|".$status;   
    $r->set("post:$postid",$post);   
    $followers = $r->smembers("uid:".$User['id'].":followers");   
    if ($followers === false) $followers = Array();   
    $followers[] = $User['id']; /* Add the post to our own posts too */  
    
    foreach($followers as $fid) {   
     $r->push("uid:$fid:posts",$postid,false);   
    }   
    # Push the post on the timeline, and trim the timeline to the   
    # newest 1000 elements.   
    $r->push("global:timeline",$postid,false);   
    $r->ltrim("global:timeline",0,1000);   
    
    header("Location: index.php");   
	   
```

	


函数的核心是foreach。 通过SMEMBERS获取当前用户的所有follower，然后循环会把帖子(post)LPUSH到每一个用户的 uid:<userid>:posts里
注意我们同时维护了一个所有帖子的时间线。为此我们还需要LPUSH到global:timeline里。
面对这个现实，你是否开始觉得:SQL里面用ORDER BY来按时间排序有一点儿奇怪? 我确实是这么想的。
####分页
现在很清楚，我们能用LRANGE来获取帖子的范围，并在屏幕上显示。代码很简单:
```PHP
  function showPost($id) {   
    $r = redisLink();   
    $postdata = $r->get("post:$id");   
    if (!$postdata) return false;   
  
    $aux = explode("|",$postdata);   
    $id = $aux[0];   
    $time = $aux[1];   
    $username = $r->get("uid:$id:username");   
	$post = join(array_splice($aux,2,count($aux)-2),"|");   
	$elapsed = strElapsed($time);   
	$userlink = "<a class=\"username\" href=\"profile.php?u=".urlencode($username)."\">".utf8entities($username)."</a>";   

	echo('<div class="post">'.$userlink.' '.utf8entities($post)."<br>");   
	echo('<i>posted '.$elapsed.' ago via web</i></div>');   
	return true;   
  }   
   
 function showUserPosts($userid,$start,$count) {   
     $r = redisLink();   
     $key = ($userid == -1) ? "global:timeline" : "uid:$userid:posts";   
     $posts = $r->lrange($key,$start,$start+$count);   
     $c = 0;   
     foreach($posts as $p) {   
         if (showPost($p)) $c++;   
         if ($c == $count) break;   
     }   
     return count($posts) == $count+1;   
 }   
  
```

当showUserPosts获取帖子的范围并传递给showPost时，showPost会简单输出一篇帖子的HTML代码。
 
Following users 关注的用户
如果用户id 1000 (antirez)想要follow用户id1000的pippo，我们做到这个仅需两步SADD:
> 
SADD uid:1000:following 1001
SADD uid:1001:followers 1000

再次注意这个相同的模式:在关系数据库里的理论里follow的用户和被follow的用户是一张包含类似following_id和follower_id的单独数据表。
用查询你能明确follow和被follow的每一个用户。在key-value数据里有一点特别，需要我们分别设置1000follow了1001并且1001被1000follow的关系。
这是需要付出的代价，但是另一方面讲，获取这些数据即简单又超快。并且这些是独立的集合，允许我们做一些有趣的事情，比如使用SINTER获取两个不同用户的集合。
这样我们也许可以在我们的twitter复制品中加入一个功能:当你访问某个人的资料页时显示"你和foobar有34个共同关注者"之类的东西。
你能够在follow.php中找到增加或者删除following/folloer关系的代码。它如你所见般平常。

###使它能够水平分割
亲爱的读者，如果你看到这里，你已经是一个英雄了，谢谢你。在讲到水平分割之前，看看单台服务器的性能是个不错的主意。
Retwis让人惊讶地快，没有任何缓存。在一台非常缓慢和高负载的服务器上，以100个线程并发请求100000次进行apache基准测试，平均占用5ms。
这意味着你可以仅仅使用一台linux服务器接受每天百万用户的访问，并且慢的跟个傻猴似的，就算用更新的硬件。
虽然，就算你有一堆用户，也许也不需要超过1台服务器来跑应用，但让我们假设我们是Twitter，需要处理海量的访问量呢?该怎么做?

###Hashing the key
第一件事是把KEY进行hash运算并基于hash在不同服务器上处理请求。有大量知名的hash算法，例如ruby客户端自带的consistent hashing
大致意思是你能把key转换成数字，并除以你的服务器数量
 server_id = crc32(key) % number_of_servers
这里还有大量因为添加一台服务器产生的问题，但这仅仅是大致的意思，哪怕使用一个类似consistent hashing的更好索引算法，
是不是key就可以分布式访问了呢?所有用户数据都分布在不同的服务器上，没有inter-keys使用到(比如SINTER，否则你需要注意要在同一台服务器上进行)
这是Redis不像memcached一样强制指定索引算法的原因，需要应用来指定。另外，有几个key访问的比较频繁。
####特殊的Keys
比如每次发布新帖，我们都需要增加global:nextPostId。单台服务器会有大量增加的请求。如何修复这个问题呢?一个简单的办法是用一台专门的服务器来处理增加请求。
除非你有大量的请求，否则矫枉过正了。另一个小技巧是ID并不需要真正地增加，只要唯一即可。这样你可以使用长度为不太可能发生碰撞的随机字符串(除了MD5这样的大小，几乎是不可能)。
完工，我们成功消除了水平分割带来的问题。
另一个问题是global:timeline。这里有个不是解决办法的解决办法，你可以分别保存在不同服务器上，并且在需要这些数据时从不同的服务器上取出来，或者用一个key来进行排序。
如果你确实每秒有这么多帖子，你能够再次用一台独立服务器专门处理这些请求。请记住，商用硬件的Redis能够以100000/s的速度写入数据。我猜测对于twitter这足够了。
