# Goby 电商平台

Goby 是一个基于微服务架构的现代化电商平台，采用 Go 语言开发后端，Next.js 开发前端，提供完整的电商解决方案。

## 项目结构
```
goby/
├── app/                   # 后端应用
│   ├── clients/           # 客户端代码
│   ├── common/            # 公共代码
│   ├── consts/            # 常量定义
│   ├── models/            # 数据模型
│   ├── rpc/               # RPC服务
│   ├── services/          # 业务服务
│   └── utils/             # 工具函数
├── docs/                  # 项目文档
└── delpoy/                # docker部署

goby-frontend/                 # 前端应用
├── app/                   # Next.js应用
├── components/            # React组件
├── public/                # 静态资源
└── styles/                # 样式文件

```

## 快速开始

### 环境要求
- Go 1.20+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 后端启动
1. 安装依赖
```bash
cd goby
go mod tidy
```

2. 配置数据库
```bash
# 修改 config/config.yaml 中的数据库配置
```

3. 启动docker
```bash
cd delpoy
docker compose-up
```

4. 启动服务
```bash
make run svc=...
```

### 前端启动
1. 安装依赖
```bash
cd goby-frontend
npm install
```

2. 启动开发服务器
```bash
npm run dev
```

## 部署说明

### 后端部署
1. 编译
```bash
cd app
go build -o goby
```

2. 运行
```bash
./goby
```

### 前端部署
1. 构建
```bash
cd web
npm run build
```

2. 启动
```bash
npm start
```

## 开发指南

### 代码规范
- 遵循 Go 标准代码规范
- 使用 ESLint 进行前端代码检查
- 提交前进行代码格式化

### 分支管理
- main: 主分支，用于生产环境
- develop: 开发分支，用于开发环境
- feature/*: 功能分支，用于新功能开发
- hotfix/*: 修复分支，用于紧急bug修复

### 提交规范
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 代码重构
test: 测试用例
chore: 构建过程或辅助工具的变动
```

## 测试

### 单元测试
```bash
# 后端测试
cd app
go test ./...

# 前端测试
cd web
npm test
```

### 集成测试
```bash
# API测试
cd tests/api
npm install
npm test

# 数据库测试
cd app
go test -v ./tests/database
```

### 端到端测试
```bash
# 安装依赖
cd tests/e2e
npm install

# 运行测试
npm test
```

### 性能测试
```bash
# 负载测试
ab -n 1000 -c 100 http://localhost:8080/api/products
```

### 安全测试
```bash
# 漏洞扫描
zap-cli quick-scan --self-contained --start-options "-config api.disablekey=true" http://localhost:3000
```

详细的测试说明请参考 [测试文档](docs/testing.md)。


## 功能特性

### 1. 用户系统
- 用户注册/登录
- 个人信息管理
- 收货地址管理
- 订单历史查询
- 收藏夹管理

### 2. 商品系统
- 商品分类展示
- 商品搜索
- 商品详情
- 商品评价
- 库存管理

### 3. 购物系统
- 购物车管理
- 商品收藏
- 订单管理
- 支付系统
- 物流跟踪

### 4. 促销系统
- 优惠券管理
- 促销活动
- 折扣规则
- 限时特价

## 技术架构

### 前端技术栈
- 框架：Next.js 13+ (App Router)
- UI组件库：Ant Design
- 状态管理：React Context
- 数据请求：Axios
- 样式方案：Tailwind CSS

### 后端技术栈
- 开发语言：Go
- 数据库：MySQL (GORM)
- 缓存：Redis
- 微服务框架：自定义微服务架构

### 系统架构
- 用户服务 (User Service)
- 商品服务 (Product Service)
- 订单服务 (Order Service)
- 购物车服务 (Cart Service)
- 支付回调服务 (PayCallback Service)
- 网关服务 (Gateway Service)

