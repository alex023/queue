# QUEUE
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/alex023/queue)](https://goreportcard.com/report/github.com/alex023/queue)
[![GoDoc](https://godoc.org/github.com/alex023/queue?status.svg)](https://godoc.org/github.com/alex023/queue)
[![Build Status](https://travis-ci.org/alex023/queue.svg?branch=master)](https://travis-ci.org/alex023/queue?branch=master)
 
# 摘要
 queue 是一个基于 mpsc 的异步队列，屏蔽业务逻辑与队列的耦合
 
# Feature
 - 提供多生产者、单消费者的消息推送
 - 提供立即关闭、业务完成关闭两种方式
 - 重复关闭无异常
 
 