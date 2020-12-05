```toml
title = "分布式对象存储系统SeaweedFS"
date = "2017-06-13 19:39:31"
update_date = "2017-06-13 19:39:31"
author = "KDF5000"
thumb = ""
tags = ["Distributed System", "Object Storage", "SeaweedFS"]
draft = false
```
SeaweedFS是Facebook的海量图片存储系统Haystack的一个开源实现，其目标是：
- 存储数亿张图片
- 快速响应的文件服务

seaweedfs提供一个Key-> file的k/v服务，不支持POSIX的文件语义。虽然是Haystack的开源实现，但是也有所不同，比如Haystack论文里提到的使用一个中心节点保存所以文件的metadate。Seaweedfs实现中，中心master是负责文件卷标的管理，具体文件的metadata有相应的valume服务器负责。这样可以大大减缓中心节点的并发的压力，并且提供了快速的文件访问**(一次磁盘操作)**。
<!--more -->
**每个文件的metadata只有50个字节，因此可以在O(1)的时间内从磁盘读取出来。**

额外的功能:
* 可以选择不同层次的副本策略：不备份，不同机架，不同数据中心
* 中心节点自动故障恢复 - 无单点问题
* 根据文件mime类型自动选择Gzip压缩
* 更新或者删除文件之后自动compaction
* 同一个集群的服务器可以有不同的磁盘大小，操作系统，文件系统
* 增加或者移除节点不会导致任何数据的re-balance
* 可选的Filer Server提供"normal"的目录和文件服务
* 可选的文件方向调整修复
* 支持Etag, Accept-Range, Last-Modified等
* 支持in-memory/leveldb/boltdb/btree 模式，以便调整性能/内存的平衡

## 相关术语
### Save File Id
fid的格式是`volumeId, filekey+fileCookie`, 其中volumeId是一个32位的无符号整数。fileKey是一个64位的无符号整数。fileCookie是一个64位的无符号整数，用来防止url暴力破解。volumeID, fileKey和fileCookie被编码成十六进制，**获得的fid需要保存到应用服务器**。

### WriteFile
写文件分为两步：
* 从master获取可以保存的volume信息
```shell
curl -X POST http://localhost:9333/dir/assign
```
可以返回volume node的相关信息，以及相应的fid
```shell
{"fid":"11,026bfba733","url":"127.0.0.1:8080","publicUrl":"127.0.0.1:8080","count":1}
```
* 保存文件到相应的volume server
```shell
curl -X PUT -F file=@/Users/KDF5000/Documents/2017/Coding/ObjectStorage/SeaweedFS/seaweedfs/weed/vim-key.png http://127.0.0.1:8080/11,026bfba733 
```
成功的话返回上传的文件名和文件大小
```
{"name":"vim-key.png","size":236702}
```

### ReadFile
读文件可以根据fid到任何一个节点(master,volumeserver)去读，他会自动的跳转到实际存储的节点。
请求的格式可以至此下面的几种格式：
```
 http://localhost:8080/3/01637037d6/my_preferred_name.jpg
 http://localhost:8080/3/01637037d6.jpg
 http://localhost:8080/3,01637037d6.jpg
 http://localhost:8080/3/01637037d6
 http://localhost:8080/3,01637037d6
```
也可以对图片进行伸缩
```
http://localhost:8080/3/01637037d6.jpg?height=200&width=200
http://localhost:8080/3/01637037d6.jpg?height=200&width=200&mode=fit
http://localhost:8080/3/01637037d6.jpg?height=200&width=200&mode=fill
```

## 架构
seaweedfs主要有两种服务类型，master server和volume server。

master server可以部署多个节点，使用raft维护一致性，主要负责记录volumeId到volume server的映射关系。

Volume server 是主要的存储节点，有一系列的volume文件(默认32G),每个volume存储很多的小文件，与之对应还有有一个所以文件，记录volume中文件的编号，偏移，大小等信息。Volume server维护本地所有volume的元信息，初始化后保存在内存中，这样当读一个文件的时候只用一次磁盘操作即读出文件的实际内容。

### 读写文件
写文件的时候需要首先与master server通信，获得可以保存文件的volume信息，master server会返回(volume id, file key, file cookie, volume node url)格式的信息用于保存文件。然后client在请求返回的volume node url以及文件路径保存文件到volume server。

读文件时候client可以根据之前记录的volume server信息，根据(volumeId, file key, file cookie)直接读取文件信息，也可以向master节点查询volume server的url, 然后去读取文件信息。

### Replication
在上传图片的时候，获取fid时候可以指定是否需要保存副本，保存的类型是什么
```
curl -X POST http://localhost:9333/dir/assign?replication=001
```
seaweedfs一共有6种副本类型：
* 000: 没有副本(默认)
* 001: 在同一个机架保存一份副本
* 010: 在同一个数据中心的不同机架保存一份副本
* 100：在不同数据中心保存一份副本
* 200：在两个不同的数据中心保存两个副本
* 110：一个副本在不同机架，一个在不同数据中心
