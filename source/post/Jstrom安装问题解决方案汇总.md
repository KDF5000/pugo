```toml
title = "Jstrom安装问题解决方案汇总"
date = "2016-06-20 16:35:09"
update_date = "2016-06-20 16:35:09"
author = "KDF5000"
thumb = ""
tags = ["JStorm"]
draft = false
```
#### jstrom 2.1.1安装错误
#####**the hostname which  supervisor get is localhost**
查看[源码(https://github.com/alibaba/jstorm/blob/master/jstorm-core/src/main/java/com/alibaba/jstorm/schedule/FollowerRunnable.java)81行
```
this.data = data;
        this.sleepTime = sleepTime;
        boolean isLocaliP;
        if (!ConfigExtension.isNimbusUseIp(data.getConf())) {
            this.hostPort = NetWorkUtils.hostname() + ":" + String.valueOf(Utils.getInt(data.getConf().get(Config.NIMBUS_THRIFT_PORT)));
            isLocaliP = NetWorkUtils.hostname().equals("localhost");
        } else {
            this.hostPort = NetWorkUtils.ip() + ":" + String.valueOf(Utils.getInt(data.getConf().get(Config.NIMBUS_THRIFT_PORT)));
            isLocaliP = NetWorkUtils.ip().equals("127.0.0.1");
        }
        try {
            if (isLocaliP) {
                throw new Exception("the hostname which Nimbus get is localhost");
            }
        } catch (Exception e1) {
            LOG.error("get nimbus host error!", e1);
            throw new RuntimeException(e1);
        }
```
可以看出Nimbus的ip设置不能为localhost和127.0.0.1，而启动nimbus server的脚本start.sh如下，使用的是`hostname -i`然后匹配storm.yaml文件找到用户设置的地址，而`hostname -i`输出的时候主机名对应的ip地址，因为不能用127.0.0.1和localhost,因此只能将主机名映射到本机ip地址，所以就要在/etc/hosts文件添加映射关系，例如本机的主机名为kdf5000，ip地址为10.30.5.64，　则需要添加下面信息
```
10.30.5.64　kdf5000
```
重新启动即可。
