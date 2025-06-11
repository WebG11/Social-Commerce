# 介绍

# 开发流程
#### gateway
1. 在 main.go 的 registerRoutes 定义接口
2. 在 hadlers 文件夹里将调用rpc将请求转发给各services
#### rpc
1. 在 idl 文件定义请求req，resp和函数
2. 使用 `make gensvc svc=...` 生成代码
#### sevices
1. 在 handler.go 里实现接口
2. 在 biz/bll 里写业务逻辑代码

# 提交流程
### 1. 创建 Issue（可选）
描述要实现的功能或者修复的bug。

### 2. 创建本地开发分支
进入你本地的Git仓库，并确保你在`main`分支上。然后，创建一个新的开发分支，例如`feature/user-auth`：

```bash
git checkout main               # 切换到 main 分支
git pull origin main             # 拉取最新的 main 分支
git checkout -b feature/user-auth # 创建并切换到开发分支
```

### 3. 开发和提交代码
在新的分支上进行功能开发。开发完成后，使用以下命令将更改提交到本地仓库：

```bash
git add .                        # 添加所有更改文件
git commit -m "Add user login and registration functionality" # 提交
git push origin feature/user-auth # 推送到远程仓库的 feature/user-auth 分支
```

### 4. 创建 Pull Request (PR)
在GitHub上，切换到你的远程仓库，进入“Pull requests”页面，点击“New pull request”：

- 选择你开发的分支（例如`feature/user-auth`。
- 在PR的标题和描述中，简要说明这次提交的功能。
- 然后点击“Create pull request”。

### 5. 等待代码审查
提交PR后，其他团队成员可以审查你的代码。如果没有问题，他们会合并你的代码。如果有修改意见，进行修改后再次提交PR。

### 6. 合并和清理
一旦PR合并完成，你可以删除本地分支：

```bash
git checkout main               # 切换回 main 分支
git branch -d feature/user-auth  # 删除本地分支
git push origin --delete feature/user-auth # 删除远程分支
```

# 目录结构
写的差不多了再改成树的形式
```
.github:        github action的文件
deploy:         部署相关代码
scripts:        被github action调用的测试逻辑
src:            后端代码
utils:         工具集
model:          各个PO的定义
rpc:            rpc相关文件

rpc/idl:        thrift 定义文件
rpc/kitex_gen:  使用 kitex + thrift 生成的go代码，包括编解码,server,client的实现

services: 各个服务的代码
services/.../biz: 业务(biz)代码
services/.../biz/bll: 业务逻辑(为了和service区分开，用的bll和dal这种C#里的概念)
services/.../biz/dal: 初始化db/redis等存储中间件
services/.../biz/dao: 操作外部数据原
services/.../biz/conf: 包含了开发和上线的yaml配置文件，以及根据当前环境读取对应文件的代码。看需不需要一个test配置
services/.../docs: 基于hertz提供的swag工具生成的swagger文档
services/.../handler.go: 本来是让kitex生成的，但没啥必要，直接从user中粘贴即可
services/.../main.go: 服务的启动代码，启动对应的服务就进对应的文件夹，然后 `go run .`
```


# 环境搭建
## go依赖
`go mod tidy`

## ~~protobuf~~ 已改为使用thrift
- 不同os:
    - Ubuntu可以: `apt install protobuf-compiler`
    - Windows: `https://github.com/protocolbuffers/protobuf/releases`
- 验证: `protoc --version`

kitex: `go install github.com/cloudwego/kitex/tool/cmd/kitex@latest`
```shell
cd rpc

kitex -module=github.com/bitdance-panic/gobuy/app/rpc idl/user.thrift
```
业务只需要写handler即可 照着user里的写就行


## swagger

1. `go install github.com/swaggo/swag/cmd/swag@latest`
2. 进入某个使用了hertz的服务
2. 参考`services/gateway/main.go`和`handler.go`的写法，为路由函数添加swagger参数
3. 在该服务的根目录运行`swag init`，生成 `docs` 文件夹和 `docs/doc.go `
4. 启动hertz, 访问`http://<ip>:<port>/swagger/index.html`

# 热重载
air: `go install github.com/air-verse/air@v1.52.3`



# go.work
允许跨模块调用,但不要在服务间调用，而是服务调 `common` 和 `rpc`
需要在各個go.mod中添加replace，可參考`gateway/go.mod`


# golangci-lint
https://golangci-lint.run/welcome/install/#install-from-sources


# gofumpt
`go install mvdan.cc/gofumpt@latest`


# 中间件版本
tidb: ...


# 本地测试方式
> 务必在提交之前进行本地测试!!!!
```
# 当前在项目根目录,与.github同级
make
# 如果不会用make，就依次执行makefile里的命令
```

# 测试服务跑起来没
```
cd services/gateway
go run .
# 再开一个终端
cd services/user
go run .
# 访问`http://localhost:8080/login?email=1234@password=1234`,请求会发到gateway，然后rpc到user，然后操作数据库
```
