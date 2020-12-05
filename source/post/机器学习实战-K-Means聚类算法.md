```toml
title = "机器学习实战:K-Means聚类算法"
date = "2015-11-07 23:11:32"
update_date = "2015-11-07 23:11:32"
author = "KDF5000"
thumb = ""
tags = ["机器学习", "算法"]
draft = false
```
#### k-means聚类算法

> 
> 优点: 容易实现
> 缺点: 可能收敛到局部最小值,在大规模数据集上收敛较慢
> 使用数据类型: 数值型护具

k-均值是发现给定数据集的k个簇的算法.k有用户决定.每一个簇通过旗质心,即簇中所有点的中心描述.
工作流程: 首先随机确定k个初始点作为其质心.然后讲护具集中的每个点分配到一个簇中,也就是分配到距其最近的质心对应的簇.这一步完成时后,每个簇的质心更新为该簇所有点的平均值.
上述过程的伪代码如下:
```
创建k个点作为起始质心(随机选择)
当任意一个点的簇分配结果改变时:
    对数据集中的每个数据点
         对每个质心
              计算质心到数据点之间的距离
         讲数据点分配到距其最近的簇
    对每一个簇,计算簇中所有点的均值幷讲均值作为质心
```

<!--more-->

##### 起始质心
随机生成指定个数的起始质心,一般可能采取选择数据点中的几个点,本文使用所提供的数据点的各个维度的最大值和最小值随机生成基于最大和最小之间的数值,代码如下:
```
# 随机取k个中心
def randCent(dataSet, k):
    n = shape(dataSet)[1]  # 列数
    centroids = mat(zeros((k, n))) # k行n列的矩阵 也就是取k个n维向量
    for j in range(n):
        minJ = min(dataSet[:, j])
        maxJ = max(dataSet[:, j])
        rangeJ = float(maxJ - minJ)
         # 生成j列向量
        centroids[:, j] = minJ + rangeJ * random.rand(k, 1) 

    return centroids
```
##### 计算两点的距离
计算两点有时候是两个向量(可以认为是高维的点)之间的距离有很多方法,在机器学习或者数据挖掘中经常需要计算两个向量的相似度,实际上也是计算两个向量的距离.计算距离的方法有很多,比如欧氏距离,曼哈顿距离,夹角余弦等等,本文采用的是欧式距离
```
# 计算欧式距离
def distEclud(vecA, vecB):
    return sqrt(sum(power(vecA-vecB, 2)))

```
##### k-means算法的实现
通过`randCent`随机选择k个质心,然后计算每个数据点与各个质心的距离,分配到距离最小的质心所在的簇,然后队每个簇根据其均值重新计算质心,然后在队每个点进行距离起算,聚类直到所有点分配结果不在改变为止.
```
# k-means算法
def kMeans(dataSet, k, distMeas=distEclud, createCent=randCent):
    """
    创建k个点作为起始质心(随机选择)
    当任意一个点的簇分配结果改变时:
        对数据集中的每个数据点
             对每个质心
                  计算质心到数据点之间的距离
             讲数据点分配到距其最近的簇
        对每一个簇,计算簇中所有点的均值幷讲均值作为质心
    """
    m = shape(dataSet)[0]
    # 第一列记录最近簇的索引,第二咧是距离
    clusterAssment = mat(zeros((m, 2)))  
    centroids = createCent(dataSet, k)
    clusterChanged = True
    while clusterChanged:
        clusterChanged = False
        for i in range(m):
            minDist = inf
            minIndex = -1
            for j in range(k):
                distJI = distMeas(centroids[j, :], dataSet[i, :])
                if distJI < minDist:
                    minDist = distJI
                    minIndex = j
            if clusterAssment[i, 0] != minIndex:
                clusterChanged = True
            clusterAssment[i, :] = minIndex, minDist ** 2
        # 更新质心的位置
        for cent in range(k):
            ptsInClust = dataSet[nonzero(clusterAssment[:, 0].A == cent)[0]]
            centroids[cent, :] = mean(ptsInClust, axis=0)
    return centroids, clusterAssment
```
##### 测试
给定一个测试文件,看看分类的效果,测试文件的内容如下,每行有两个值,一个是x坐标一个是y坐标,使用上面的算法看看聚类效果如何
```
1.658985	4.285136
-3.453687	3.424321
4.838138	-1.151539
-5.379713	-3.362104
0.972564	2.924086
-3.567919	1.531611
0.450614	-3.302219
-3.487105	-1.724432
2.668759	1.594842
-3.156485	3.191137
...
```
使用python的图形库matplotlib将测试文件的数据点绘制再二维平面中
```
if __name__ == "__main__":
    dataMat = mat(loadDataSet('testSet.txt'))
    centroids, clusterAssment = kMeans(dataMat, 4)
    pl.plot(centroids[:, 0], centroids[:, 1], 'r')
    pl.plot(dataMat[:, 0], dataMat[:, 1], 'bo')
    pl.show()
```
结果如下:
![](@media/archive/img_figure_1.png)
图中红色的点是四个簇的质心,可以看出k-means的聚类效果还是挺好的.

[源码下载](https://github.com/KDF5000/MLPractice/tree/master/ch10)
