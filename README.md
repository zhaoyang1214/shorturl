# Short Url

## 介绍
短网址系统是使用<a href="https://github.com/zhaoyang1214/ginco">Ginco</a>框架开发的。

## 使用
1、启动项目(测试)
```shell
go run main.go start
```

2、初始化数据库
```shell
go run main.go migrate -k default
```

3、`users`表添加邮箱（JWT认证）

4、访问`/swagger/index.html`查看接口文档

5、根据接口文档创建短链

6、访问`ip:port/hashKey`，自动跳转