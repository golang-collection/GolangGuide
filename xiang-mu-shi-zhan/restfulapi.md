# RestfulAPI

## 项目地址

{% embed url="https://github.com/golang-collection/Go-RestfulAPI" %}

## Go-RestfulAPI

 ![](https://camo.githubusercontent.com/54fdbe8888c0a75717d7939b42f3d744b77483b0/687474703a2f2f6a617977636a6c6f76652e6769746875622e696f2f73622f69636f2f617765736f6d652e737667) ![](https://camo.githubusercontent.com/1ef04f27611ff643eb57eb87cc0f1204d7a6a14d/68747470733a2f2f696d672e736869656c64732e696f2f7374617469632f76313f6c6162656c3d254630253946253843253946266d6573736167653d496625323055736566756c267374796c653d7374796c653d666c617426636f6c6f723d424334453939) [![](https://camo.githubusercontent.com/41e8e16b771d56dd768f7055354613254961d169/687474703a2f2f6a617977636a6c6f76652e6769746875622e696f2f73622f6769746875622f677265656e2d666f6c6c6f772e737667)](https://github.com/SuperSupeng) [![](https://img.shields.io/github/issues/golang-collection/Go-RestfulAPI)](https://github.com/golang-collection/Go-RestfulAPI/issues) [![](https://img.shields.io/github/forks/golang-collection/Go-RestfulAPI)](https://github.com/golang-collection/Go-RestfulAPI/network/members) [![](https://img.shields.io/github/stars/golang-collection/Go-RestfulAPI)](https://github.com/golang-collection/Go-RestfulAPI/stargazers) [![](https://img.shields.io/github/license/golang-collection/go-crawler-distributed)](https://github.com/golang-collection/Go-RestfulAPI/blob/master/LICENSE) [![](https://camo.githubusercontent.com/013c283843363c72b1463af208803bfbd5746292/687474703a2f2f6a617977636a6c6f76652e6769746875622e696f2f73622f69636f2f7765636861742e737667)](https://github.com/golang-collection/Urban-computing-papers/blob/master/wechat.md)

Github：[https://github.com/golang-collection/Go-RestfulAPI](https://github.com/golang-collection/Go-RestfulAPI)

本项目为Go语言构建Restful API服务，项目采用gin构建，包括jwt，swagger，make，bash等模块，通过docker一键部署在nginx服务器。

## 目录结构

* conf 用于存放配置文件
* config 用于读取配置文件
* docs swagger生成文档
* handler 构建具体的api服务
  * sd 健康检查
  * user user操作
* model 结构体定义
* pkg 常用功能模块
  * auth 用户加密
  * constvar 分页查询
  * db 数据库连接池
  * errno 统一错误码配置
  * file 文件操作
  * logging 统一log
  * token 生成token
  * version 获取项目版本
* router 配置路由规则
  * middleware 中间件
* runtime 存储日志
* service 用于存放服务层逻辑
* util 常用工具
* admin.sh 启动，停止服务脚本
* Makefile 编译构建项目
* main.go 项目统一入口

## 配置文件

需要定制自己的配置文件

在conf下创建config.json文件

配置样例如下

```javascript
{
  "mysql": {
    "user": "",
    "password": "",
    "host": "",
    "db_name": ""
  },
  "redis": {
    "host": ""
  },
  "rabbitmq": {
    "user": "",
    "password": "",
    "host": ""
  }
}
```

## Installation

将项目部署到本地或云端提供以下两种方式：

* Direct Deploy
* Docker\(Recommand\)

#### Pre-requisite

* Go 1.13.6
* MySQL 5.7
* 相关依赖：[go.mod](./go.mod)

### Quick Start

Please open the command line prompt and execute the command below. Make sure you have installed `docker-compose` in advance.

```text
git clone git@github.com:golang-collection/Go-RestfulAPI.git
cd Go-RestfulAPI
docker-compose up -d
```

Next, you can look into the `docker-compose.yml` \(with detailed config params\).

### Run

#### Docker

Please use `docker-compose` to one-click to start up. By doing so, you don't even have to configure RabbitMQ , Reds, MySQ,ElasticSearch. Create a file named `docker-compose.yml` and input the code below.

```text
version: '3.3'
services:
```

Then execute the command below, and the project will start up. Open the browser and enter `http://localhost:8080` to see the UI interface.

```text
docker-compose up
```

#### Direct

```bash
cd Go-RestfulAPI
go run main.go
```

## Appendix

* docker安装：[https://docs.docker.com/](https://docs.docker.com/)
* docker-compose安装：[https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

## License

[MIT](https://github.com/golang-collection/Go-RestfulAPI/blob/master/LICENSE)

Copyright \(c\) 2020 Knowledge-Precipitation-Tribe

## 参考文献

\[1\] [https://juejin.im/book/6844733730678898702](https://juejin.im/book/6844733730678898702)

