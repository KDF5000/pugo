```toml
title = "对象存储比较 - Haystack"
date = "2017-06-25 20:54:24"
update_date = "2017-06-25 20:54:24"
author = "KDF5000"
thumb = ""
tags = ["Distributed System", "Object Storage", "Haystack"]
draft = false
```
经过一段时间的调研，对分布式对象系统有了一点浅显的认识，暂不谈和文件系统的区别(其实还是有点傻傻分不清)，姑且认为对象存储和文件系统最大的区别就是API，文件系统提供了完整POSIX语义，往往具有层次化的目录结构，对文件可以进行精细的操作(open, read, write, seek, delete等)，相对会复杂一些。相比之下，对象存储就简单的多，对象本义是**B**inary **L**arge **OB**ject(BLOB), 是一个大的二进制文件，是作为一个整体出现，不能对其进行修改，因此对象存储系统有一个显著的特点就是Immutable，即不变性，正是因为这个特点对象存储系统一般只提供put, get, 和delete的操作，并且都是以key/val的形式，val一般是整个blob. 具体到实际的资源，主要指的是像图片，视频，文档，源码，二进制程序等这样的文件，这些文件往往都很小，一般在100k左右，视频可能会大一些，能达到几十兆甚至好几个G。

综合上面的一些特点学术界和工业界对对象存储系统也区别于文件系统做了很多优化，其中比较有名的当属10年Facebook在OSDI公开的一个分布式对象存储系统Haystack，针对其内部的业务场景(大量的小图片存储)，创新的提出了一个小文件合并成大文件的方法去存储海量图片的需求，并且得到了很好的实践考研。除了facebook之外，Amazon, twitter， linkdin等知名厂商也都提出了自己的对象存储，其中Amazon的S3作为一个商业服务也取得了巨大的成功。

本文主要调研了Facebook， twitter, linkdin内部使用的对象存储系统，企图从其系统的构建初衷到设计进行一次全面的探索，最后再进行一个对比，方便有需求的同学选择适合自己的系统。

<!--more-->

## Fackbook
Facebook作为全球最大的图片分享社区，其图片，视频的存储量用海量形容当之无愧，2010年其第一次将内部使用的海量图片存储系统Haystack发表在了系统界的顶级会议OSDI，虽然现在看整个系统设计非常简单，但是已经发表便引起了工业界和学术界的顶礼膜拜，从此各种海量小文件存储系统无不拿其说事，其提出的通过文件合并从而减少元数据，最终减少磁盘io的思想经久不衰，认真研读了当年的文章，整篇文件架构清晰，丝毫不拖泥带水。继Haystack之后，2014年经过四年的发展，其内部系统又做了大量的改进，并且新增加了一个专门用于存储冷数据的存储系统，大大的提高了存储成本，同样也是发表在了OSDI(怎么感觉OSDI是他们家开的，想发就发，当然也可能是这些系统都是经过实践的考研，这也是系统界比较看重的)。

### Haystack
在介绍Haystack之前，先看一下Facebook面临着什么样的问题，从论文里的数据可以得知，Facebook在2010年，其内部已经存储2600亿张图片，并且每周还持续增加100万，这样的体量，恐怕一般的公司是早就那以应对了。那么在开发Haystack之前Facebook内部使用的系统是什么样的呢？当时他们使用的还是基于NAS的存储系统，存储服务器讲NAS的Volumes挂载在本地，然后每个文件放在不同的目录，每个目录存放上千张图片，这样会有什么问题呢？在传统的POSIX语义的文件系统下，每个目录对应的元数据会非常大，并且每个文件都有自己的元数据，因此元数据的访问变成了主要的瓶颈，为此他们也尝试了一些解决方案，甚至修改了操作源码，增加了一个`open_by_filehandler`的方法，将每个文件的handler保存在memcached，这样确实能够减少磁盘操作，减少文件的响应时间。但是很快他们发现，被访问的图片很多都是以前没有访问过的，也就是说很多时候文件handler都是不在缓存中的，除非所有文件handler都在缓存里，这就需要大量的内存，显然不靠谱。

所以，现在Facebook面临的主要问题是：
- 海量图片
- 元数据量大，访问元数据成为瓶颈
- 如何在低成本，高可用的情况下，保证系统的性能

但是Facebook的大神们并没有屈服，经过两年时间的研发上线运行，Haystack诞生了！Haystack解决方案的核心点就是通过减少文件系统元数据的大小，保证其能够常驻内存，从而减少磁盘的操作次数。工程师们分析了Facebook的图片数据，图片有个显著的特征：一次上传，经常读，永不修改，很少删除。针对这个特点，他们认为文件系统的元数据里有很多属性，比如权限等都是多余的，因此尽可能的减少元数据的大小，并将其保存到内存，这样可以只用一次的磁盘操作就从文件里读出图片数据，并且他们讲多个图片追加写入到一个大文件里，从而进一步的减少了元数据的数量。

Haystack的整体架构如图：
![haystack_arch.jpeg](@media/archive/blog/images/haystack_arch.jpeg?imageView/0/w/500/)
从架构图可以看出Haystack一共包含了四个核心组件：CDN, Haystack Directory, Haystack Cache, Haystack Store. 下面简略介绍一下每个组件的主要功能，最后我们再给出读写操作的详细流程。

#### Haystack Directory
前面我们提到，Haystack有一个很重要的优化就是将很多图片合并到一个大文件里，论文里讲一个大文件成为Physical Volume， 并且每一个Physical Volume都有多个副本(一般是3份)，因此在这之上有抽象出了一个Logical Volume，一个Logical Volume对应几个Physical Volume，这个对应关系以及其他一些状态信息都有Directory保存。废话不说，Haystack Directory主要有下面几个功能：
- 逻辑卷标到物理卷标的映射
- 读写负载均衡
- 决定一个读请求是通过CDN还是Cache
- 记录Logical volumes是否只读(只读的粒度是机器层次)
 
所有这些信息全部都保存在Memcached里，从而减少了访问延迟。

#### CDN
CDN就是普通的CDN，接受请求，缓存数据，如果没有命中则访问Haystack Cache。

#### Haystack Cache
Cache相当于一个内部的CDN，因为有两层的CDN所以Haystack Cache的缓存策略就不太一样，并不是来一个请求就缓存，Haystack Cache只有下面两种情况才进行缓存:
- 当请求直接来自浏览器，而不是CDN
- 图片来自write-enable的存储机器

对于第一种情况，如果一个请求来自CDN说明CDN没有缓存，那么返回给CDN后，cdn会缓存，所以Cache无需缓存；对于第二种情况，基于一种假设：新上传的图片有很大可能被读取。并且当前的设计只读或者只写的性能好，但是同时读写的性能不太好，所以最好将新写的数据缓存。

#### Haystack Store
Haystack Store是最核心的一部分，所有的图片文件全部都保存在Store Server里。其主要有下面的特点：
- 保存实际的Physical Volume(大文件)，路径可能是`/hay/haystack_<logical volume id>`
- 所有的物理卷标文件描述符都保存在内存
- 对一个图片的读取，photoid, offset, size都在内存读取，不需要磁盘操作
维护photoid -> (fd, flag, offset, size)的映射

Haystak Store除了实际的Phystical Volume外，还有一个重要的数据文件：索引文件。索引文件主要作用是可以在机器重启的时候帮助快速的重建内存mapping。索引文件在写入一个needle时候，同步的写Physical Volume和内存map， 异步写索引文件；同时，删除一个文件的时候，同步置内存map对应的flag和Pyhsical Volume的flag，但是**不跟新**索引文件信息。因此可能存在两个问题：
- needle可能不在索引文件中
- 索引文件不能反映删除了的needle

针对第一个问题，在重建mapping的时候检查不在索引文件中的needle并append到索引文件索引文件的最后一个needle记录，对应volume中的needle后的所有needle都要重新append index；针对第二个问题，从物理needle中读取needle后，检查flag，如果是删除的话更新in-memory map。

### 读写删操作
- 读操作
    先根据上传时候保存的需要读的文件的(logicalVolume id, key, alternate key, cookie)等信息请求Directory获取实际的请求URL，获取的格式如下:
    ```
    http://<CDN>/<Cache>/<Machine Id>/<Logical volume Id, PhotoId>
    ```
    然后根据该URL请求CDN服务器，服务器如果没有命中则会请求Cache服务器, 同样Cache服务器没有命中则会请求URL中的Machine节点，该节点在内存内查找相关映射，如果文件没有被删除(flag)，则从磁盘读取，然后验证cookie，checksum等，如果成功则返回相应数据，否则返回读取错误。
- 写操作
同样从从Directory获取可以写的logical volume id，以及key, alternate key ，cookie
然后同步写所有logical volume 对应的physical volume。
- 删除操作
删除操作会同步置map和volume里相应的flag，volume的数据后续compact的时候会进行删除。
- 更新操作
不允许更新，只能append相同的key的文件。如果新的needle append到不同的logical volume，则更新Directory的application metadata，保证下次请求不会读取到旧版本数据；如果append到相同的logical volume，则根据offset来更新map，最大的是最新的。

到此为止Haystack比较核心的组件以及操作都进行了有一个回顾，除了这些之外，Haystack也提供了一系列的故障恢复机制，主要包括故障检测和修复技术。同时也做了一些内存以及上传等的优化。整体来说Haystack简单有效，能够有些的解决当前的场景的问题。当然也存在一些问题，比如Directory服务器存在单点问题，以及Storage Server的内存映射如果很大的时候怎么办，总不能全部放内存吧，当然可以限制每个节点的Physical Volume节点的数量限制，不过看来就就没有那么优美了。下一篇会介绍一些2016年Facebook在OSDI发表的另一个和Haystack相辅相成的存储系统，解决了Haystack的一些问题。
