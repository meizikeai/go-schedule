# go-schedule

## 项目构架

基于[cron](https://github.com/robfig/cron/v3)构建，支持[mysql](https://github.com/go-sql-driver/mysql)、[redis](https://github.com/go-redis/redis)查询，依赖[golang](https://golang.google.cn/dl)环境。

## 项目地址

https://github.com/meizikeai/go-schedule.git

## 项目结构

| 路径    | 描述          | 详情 |
| ------- | ------------- | ---- |
| config  | config        | --   |
| crontab | crontab       | --   |
| libs    | lib           | --   |
| models  | mysql / redis | --   |
| go.mod  | go modules    | --   |

## 开发环境

  + 克隆项目 - `$ git clone https://github.com/meizikeai/go-schedule.git`
  + 安装依赖 - `$ cd go-schedule && go mod tidy`
  + 启动项目 - `$ go run .`

推荐[Visual Studio Code](https://code.visualstudio.com)编辑器，开发依赖[Tools](https://github.com/golang/vscode-go/blob/master/docs/tools.md)，请正确安装好后再开始。

## 项目变量

在测试、线上环境需要设置 `GO_MODE` 的值来区分相应环境。

需要在 `～/.zshrc` 里给定 `export GO_MODE=release/test` 后，程序中通过 `mode := os.Getenv("GO_MODE")` 获取。

如果 `GO_MODE` 未设定，那么运行的是 `test` 环境。

为确保正确写日志，本地开发添加 `export GO_ENV=debug` 环境变量。

## 自动部署

如果使用 GitLab 作仓库，可以使用 https://github.com/meizikeai/gitlab-golang-shell.git 跑CI/CD，项目默认有 .gitlab-ci.yml 文件，请君参考！

## 帮助文档

 `通过以下命令行执行`

```sh
# 第一步
$ cd ~/go-schedule
$ GOOS=linux GOARCH=amd64 go build -o go-schedule main.go

# 后台运行，关掉终端会停止运行
$ ~/go-schedule/go-schedule &

# 后台运行，关掉终端也会继续运行
$ nohup ~/go-schedule/go-schedule > go-schedule.log 2>&1 &

# 第二步
$ 执行方法见 帮助文档
$ 项目配置件 /etc/supervisor/conf.d
$ 执行项目名 go-schedule
```

```sh
# program 为 [program:go-schedule] 里配置的值
# start、restart、stop、remove、add 都不会载入最新的配置文件

# start      启动程序
# status     查看程序状态
# stop       关闭程序
# tail       查看进程日志
# update     重启配置文件修改过的程序
# reload     停止程序，重新加载所有程序
# reread     读取有更新（增加）的配置文件，不会启动新添加的程序
# restart    重启程序

# 执行某个进程
$ supervisorctl restart program

# 一次性执行全部进程
$ supervisorctl restart all

# 载入最新的配置文件，停止原有进程并按新的配置启动所有进程
$ supervisorctl reload

# 根据最新的配置文件，启动新配置或有改动的进程，配置没有改动的进程不重启
$ supervisorctl update

# 查看运行状态
$ supervisorctl status
```

## 学习资料

**Go 语言设计与实现**

  + https://draveness.me/golang

**幼麟实验室 Golang 合辑**

  + https://space.bilibili.com/567195437

**Golang Example**

  + https://gobyexample.com
  + https://gowebexamples.com
