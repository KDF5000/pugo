```toml
title = "snort规则入门"
date = "2016-09-07 16:50:17"
update_date = "2016-09-07 16:50:17"
author = "KDF5000"
thumb = ""
tags = ["Tech"]
draft = false
```
### snort规则

一种简单的，轻量级的规则描述语言，包括连个逻辑部分：格则头和规则选项。

**规则头**：规则的动作，协议，源和目标ip与网络掩码，源和目标端口信息

**规则选项**: 报警信息内容和要检查的包的具体部分

```
alert tcp any any -> 192.168.1.0/24 111 (content:"|00 01 86 a5|"; msg: "mountd access";)
```

 第一个括号前的部分是规则头（rule header），包含的括号内的部分是规则选项（rule options）。规则选项部分中冒号前的单词称为选项关键字（option keywords）。注意，不是所有规则都必须包含规则选项部分，选项部分只是为了使对要收集或报警，或丢弃的包的定义更加严格。组成一个规则的所有元素 对于指定的要采取的行动都必须是真的。当多个元素放在一起时，可以认为它们组成了一个逻辑与（AND）语句。同时，snort规则库文件中的不同规则可以 认为组成了一个大的逻辑或（OR）语句

<!-- more -->

#### 规则头

规则的头包含了定义一个包的who，where和what信息，以及当满足规则定义的所有属性的包出现时要采取的行动

* 规则动作

  snort中有五种动作： alert, log, pass, activate. dynamic

  1、Alert-使用选择的报警方法生成一个警报，然后记录（log）这个包。
  2、Log-记录这个包。
  3、Pass-丢弃（忽略）这个包。
  4、activate-报警并且激活另一条dynamic规则。
  5、dynamic-保持空闲直到被一条activate规则激活，被激活后就作为一条log规则执行。

* 协议

  四种： tcp, ump, icmp, ip

* ip地址

  表示处理一个给定规则的ip地址和端口信息。 any用来定义任何地址。snort没有提供ip地址查询域名的机制。地址直接有数字型ip地址和一个cidr块组成。Cidr块只是作用在规则地址和要检查的进入任何包的网络掩码。/24表示c类网络， /16表示b类网络，/32表示一个特定的机器的地址。

  例如，192.168.1.0/24代表从192.168.1.1到192.168.1.255的 地址块。在这个地址范围的任何地址都匹配使用这个192.168.1.0/24标志的规则。这种记法给我们提供了一个很好的方法来表示一个很大的地址空间。

  有一个操作符可以应用在ip地址上，它是否定运算符（negation operator）。这个操作符告诉snort匹配除了列出的ip地址以外的所有ip地址。否定操作符用"！"表示。下面这条规则对任何来自本地网络以外的流都进行报警。

  ```
  alert tcp !192.168.1.0/24 any -> 192.168.1.0/24 111 (content: "|00 01 86 a5|"; msg: "external mountd access";)
  ```

  这个规则的ip地址代表"任何源ip地址不是来自内部网络而目标地址是内部网络的tcp包"。
  你也可以指定ip地址列表，一个ip地址列表由逗号分割的ip地址和CIDR块组成，并且要放在方括号内“[”，“]”。此时，ip列表可以不包含空格在ip地址之间。下面是一个包含ip地址列表的规则的例子。

  ```
  alert tcp ![192.168.1.0/24,10.1.1.0/24] any -> [192.168.1.0/24,10.1.1.0/24] 111 (content: "|00 01 86 a5|"; msg: "external mountd access";)
  ```

* 端口号

  端口号可以用几种方法表示，包括"any"端口、静态端口定义、范围、以及通过否定操作符。"any"端口是一个通配符，表示任何端口。静态端口定 义表示一个单个端口号，例如111表示portmapper，23表示telnet，80表示http等等。端口范围用范围操作符"："表示。范围操作符 可以有数种使用方法，如下所示：

  >  log udp any any -> 192.168.1.0/24 1:1024

  记录来自任何端口的，目标端口范围在1到1024的udp流

  > log tcp any any -> 192.168.1.0/24 :6000

  记录来自任何端口，目标端口小于等于6000的tcp流

  > log tcp any :1024 -> 192.168.1.0/24 500:

  记录来自任何小于等于1024的特权端口，目标端口大于等于500的tcp流


  端口否定操作符用"！"表示。它可以用于任何规则类型（除了any，这表示没有，呵呵）。例如，由于某个古怪的原因你需要记录除x windows端口以外的所有一切，你可以使用类似下面的规则：

  > log tcp any any -> 192.168.1.0/24 !6000:6010

#### 方向操作符

方向操作符"->"表示规则所施加的流的方向。方向操作符左边的ip地址和端口号被认为是流来自的源主机，方向操作符右边的ip地址和端口信 息是目标主机，还有一个双向操作符"<>"。它告诉snort把地址/端口号对既作为源，又作为目标来考虑。这对于记录/分析双向对话很方 便，例如telnet或者pop3会话。用来记录一个telnet会话的两侧的流的范例如下：

> log !192.168.1.0/24 any <> 192.168.1.0/24 23

#### Activate和dynamic规则

注：Activate 和 dynamic 规则将被tagging 所代替。在snort的将来版本，Activate 和 dynamic 规则将完全被功能增强的tagging所代替。

Activate 和 dynamic 规则对给了snort更强大的能力。你现在可以用一条规则来激活另一条规则，当这条规则适用于一些数据包时。在一些情况下这是非常有用的，例如你想设置一 条规则：当一条规则结束后来完成记录。Activate规则除了包含一个选择域：activates外就和一条alert规则一样。Dynamic规则除 了包含一个不同的选择域：activated_by 外就和log规则一样，dynamic规则还包含一个count域。

Actevate规则除了类似一条alert规则外，当一个特定的网络事件发生时还能告诉snort加载一条规则。Dynamic规则和log规则类似，但它是当一个activate规则发生后被动态加载的。把他们放在一起如下图所示：

> activate tcp !$HOME_NET any -> $HOME_NET 143 (flags: PA; content: "|E8C0FFFFFF|/bin"; activates: 1; msg: "IMAP buffer overflow!";)
>
> dynamic tcp !$HOME_NET any -> $HOME_NET 143 (activated_by: 1; count: 50;)

#### 规则选项

规则选项组成了snort入侵检测引擎的核心，既易用又强大还灵活。所有的snort规则选项用分号"；"隔开。规则选项关键字和它们的参数用冒号"："分开。按照这种写法，snort中有42个规则选项关键字。

msg - 在报警和包日志中打印一个消息。
logto - 把包记录到用户指定的文件中而不是记录到标准输出。
ttl - 检查ip头的ttl的值。
tos 检查IP头中TOS字段的值。
id - 检查ip头的分片id值。
ipoption 查看IP选项字段的特定编码。
fragbits 检查IP头的分段位。
dsize - 检查包的净荷尺寸的值 。
flags -检查tcp flags的值。
seq - 检查tcp顺序号的值。
ack - 检查tcp应答（acknowledgement）的值。
window 测试TCP窗口域的特殊值。
itype - 检查icmp type的值。
icode - 检查icmp code的值。
icmp_id - 检查ICMP ECHO ID的值。
icmp_seq - 检查ICMP ECHO 顺序号的值。
content - 在包的净荷中搜索指定的样式。
content-list 在数据包载荷中搜索一个模式集合。
offset - content选项的修饰符，设定开始搜索的位置 。
depth - content选项的修饰符，设定搜索的最大深度。
nocase - 指定对content字符串大小写不敏感。
session - 记录指定会话的应用层信息的内容。
rpc - 监视特定应用/进程调用的RPC服务。
resp - 主动反应（切断连接等）。
react - 响应动作（阻塞web站点）。
reference - 外部攻击参考ids。
sid - snort规则id。
rev - 规则版本号。
classtype - 规则类别标识。
priority - 规则优先级标识号。
uricontent - 在数据包的URI部分搜索一个内容。
tag - 规则的高级记录行为。
ip_proto - IP头的协议字段值。
sameip - 判定源IP和目的IP是否相等。
stateless - 忽略刘状态的有效性。
regex - 通配符模式匹配。
distance - 强迫关系模式匹配所跳过的距离。
within - 强迫关系模式匹配所在的范围。
byte_test - 数字模式匹配。
byte_jump - 数字模式测试和偏移量调整。





-----

#### 参考

Snort 中文手册: [http://man.chinaunix.net/network/snort/Snortman.htm ](http://man.chinaunix.net/network/snort/Snortman.htm)
