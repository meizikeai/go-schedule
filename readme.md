# go-schedule

### 项目构架

基于[cron](https://github.com/robfig/cron/v3)构建，支持[mysql](https://github.com/go-sql-driver/mysql)、[redis](https://github.com/go-redis/redis)查询，依赖[golang](https://golang.google.cn/dl)环境。

#### 项目地址

https://github.com/meizikeai/go-schedule.git

#### 项目结构

| 路径          | 描述               | 详情 |
| ------------- | ---------------------- | ---- |
| conf          | config                 | --   |
| libs          | lib                    | --   |
| models        | mysql / redis          | --   |
| task          | schedule               | --   |
| go.mod        | go modules             | --   |

#### 开发环境

  + 克隆项目 - `$ git clone https://github.com/meizikeai/go-schedule.git`
  + 安装依赖 - `$ cd go-schedule && go mod tidy`
  + 启动项目 - `$ fresh / go run .`

推荐[Visual Studio Code](https://code.visualstudio.com)编辑器，开发依赖[Tools](https://github.com/golang/vscode-go/blob/master/docs/tools.md)，请正确安装好后再开始。

#### 项目变量

在测试、线上环境需要设置 `GIN_MODE` 的值来区分相应环境。

需要在 `～/.zshrc` 里给定 `export GIN_MODE=release/test` 后，程序中通过 `mode := os.Getenv("GIN_MODE")` 获取。

如果 `GIN_MODE` 未设定，那么运行的是线上环境。

为了区分本地开发，增加 `export GIN_ENV=debug` 环境变量。

#### 自动部署

如果使用 GitLab 作仓库，可以使用 https://github.com/meizikeai/gitlab-golang-shell.git 跑CI/CD，项目默认有 .gitlab-ci.yml 文件，请君参考！

#### 学习资料

**Go 语言设计与实现**

  + https://draveness.me/golang

**幼麟实验室 Golang 合辑**

  + https://space.bilibili.com/567195437

**Golang Example**

  + https://gobyexample.com
  + https://gowebexamples.com

**学习 Go 语言**

  + https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/preface.md
  + https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/directory.md

**Go 操作 Redis 实战**

  + https://www.cnblogs.com/itbsl/p/14198111.html
