```toml
title = "机器学习实战:二分k-均值聚类算法"
date = "2015-11-08 22:50:54"
update_date = "2015-11-08 22:50:54"
author = "KDF5000"
thumb = ""
tags = ["机器学习", "算法"]
draft = false
```
前面介绍过k-means聚类算法,通过不断的更新簇质心直到收敛为止,但是这个收敛是局部收敛到了最小值,并没有考虑全局的最小值.

那么一个聚类算法怎么才能称得上效果好呢?要想评价一个算法的好坏,首先需要有一个标准,这也是我们设计算法的时候要首先考虑的,我们设计算法的目的是什么,设计的算法要达到的什么效果,这样才能明确设计算法的目标.一种度量聚类算法效果的指标是SSE(Sum of Squard Error, 误差平方和),也就是前面介绍的计算数据点与质心的欧氏距离的平方和.SSE越小说明数据点越接近它们的质心,聚类效果越好.

<!--more-->

#### 二分k均值算法

为了克服k-均值算法的局部最小值问题,有人提出了二分k均值算法.该算法首先将所有点作为一个簇,然后讲该簇一分为二.之后选择其中一个簇继续划分,选择哪一个簇进行划分取决于对其划分是否可以最大程度降低SSE的值.上述基于SSE的划分过程不断重复,直到得到用户指定的簇数目为止.

二分k均值算法的伪代码形式如下:
>> 
将所有点看成一个簇
当簇数目小于k时
    对于每一个簇
       计算总误差
       在给定的簇上进行2均值聚类
       计算讲该簇一分为二之后的总误差
    选择使得总误差最小的那个簇进行划分操作

下面是python的实现
```
def binKMeans(dataSet, k, distMeas=distEclud):
    """
    :param dataSet:
    :param k:
    :param distMeas:
    :return:
    选择一个初始的簇中心(取均值),加入簇中心列表
    计算每个数据点到簇中心的距离
    当簇的个数小于指定的k时
          对已经存在的每个簇进行2-均值划分,并计算其划分后总的SSE,找到最小的划分簇
          增加一个簇幷更新数据点的簇聚类
    """
    m = shape(dataSet)[0]
    clusterAssment = mat(zeros((m, 2)))
    # 创建一个初始簇, 取每一维的平均值
    centroid0 = mean(dataSet, axis=0).tolist()[0]
    centList = [centroid0]  # 记录有几个簇
    for j in range(m):
        clusterAssment[j, 1] = distMeas(mat(centroid0), dataSet[j, :]) ** 2
    while len(centList) < k:
        lowestSSE = inf
    #     # 找到对所有簇中单个簇进行2-means可以是所有簇的sse最小的簇
        for i in range(len(centList)):
            # 属于第i簇的数据
            ptsInCurrCluster = dataSet[nonzero(clusterAssment[:, 0].A == i)[0], :]
            # print ptsInCurrCluster
            # 对第i簇进行2-means
            centroidMat, splitClustAss = kMeans(ptsInCurrCluster, 2, distMeas)
            # 第i簇2-means的sse值
            sseSplit = sum(splitClustAss[:, 1])
            # 不属于第i簇的sse值
            sseNotSplit = sum(clusterAssment[nonzero(clusterAssment[:, 0].A != i), 1])
            if (sseSplit + sseNotSplit) < lowestSSE:
                bestCentToSplit = i
                bestNewCents = centroidMat
                bestClustAss = splitClustAss.copy()
                lowestSSE = sseNotSplit + sseSplit
        # 更新簇的分配结果
        #新增的簇编号
        bestClustAss[nonzero(bestClustAss[:, 0].A == 1)[0], 0] = len(centList)
        #另一个编号改为被分割的簇的编号
        bestClustAss[nonzero(bestClustAss[:, 0].A == 0)[0], 0] = bestCentToSplit  #
        # 更新被分割的的编号的簇的质心
        centList[bestCentToSplit] = bestNewCents[0, :].tolist()[0]
        # 添加新的簇质心
        centList.append(bestNewCents[1, :].tolist()[0])
        # 更新原来的cluster assment
        clusterAssment[nonzero(clusterAssment[:, 0].A == bestCentToSplit)[0], :] = bestClustAss
    return mat(centList), clusterAssment
```
#### 测试
给定一个测试文件,看看分类的效果,测试文件的内容如下,每行有两个值,一个是x坐标一个是y坐标,使用上面的算法看看聚类效果如何
```
3.275154	2.957587
-3.344465	2.603513
0.355083	-3.376585
1.852435	3.547351
-2.078973	2.552013
-0.993756	-0.884433
2.682252	4.007573
-3.087776	2.878713
-1.565978	-1.256985
2.441611	0.444826
-0.659487	3.111284
...
```
使用python的图形库matplotlib将测试文件的数据点绘制再二维平面中
```
if __name__ == "__main__":
    dataMat = mat(loadDataSet('testSet2.txt'))
    centroids, clusterAssment = binKMeans(dataMat, 4)
    pl.plot(centroids[:, 0], centroids[:, 1], 'ro')
    pl.plot(dataMat[:, 0], dataMat[:, 1], 'bo')
    pl.show()
```
结果如下:
![](@media/archive/img_figure_2.png)
图中红色的点是四个簇的质心,可以看出2分k均值算法的聚类效果还是挺好的.

[源码下载](https://github.com/KDF5000/MLPractice/tree/master/ch10)
    
   
