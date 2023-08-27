# kit-study
- 微服务 go-kit study 后续会根据时间需求 用go-kratos 进行重构

# v1- kit base service & fx demo
- 建立 kit 基础服务
- fx 依赖注入demo

# v2- middleware & request uuid & log & project-layout
- 建立中间件 
- 请求uuid 
- log集成zap方便扩展 
- project-layout目录结构初步设定

# v3- jwt & error & log & auth middleware
- 建立 jwt 认证
- error 自定义错误码
- log 封装zap
- auth中间件

# v4- rate-limit & grpc fx,grpc client,server & makefile protoc & log errorhandler 
- 建立限流 rate-limit
- grpc服务建立 transport, endpoint层rpc改造
- grpc_server grpc_client 建立
- makefile scripts 雏形建造 初步做了protoc 命令建立 后续添加其他make命令
- fx 兼容grpc_server
- log 日志库改造 errorhandler ZapLogger 全局使用

# v5- etcd 服务注册与发现