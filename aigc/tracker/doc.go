// Package tracker 指标和日志采集，初始化 Client 后，使用 Background 进行日志和指标采集，main 函数使用 defer b.SendAll() 发送所有缓存中的数据
package tracker
