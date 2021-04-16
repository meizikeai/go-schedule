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

#### 学习资料

  + https://crontab.guru
  + https://github.com/robfig/cron
