package service

// 抽象:对应 Service 安装中间件 (serivce加一层装饰)

const ContextReqUUid = "req_uuid"

type NewMiddlewareServer func(Service) Service
