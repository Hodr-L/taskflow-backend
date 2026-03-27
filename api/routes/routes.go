package routes

// 路由注册包
// 将路由按功能模块拆分到不同文件，便于管理和协作

// 导入所有路由注册函数，这里不直接导出，而是通过RegisterAllRoutes统一调用
// 避免外部直接依赖各个RegisterXXXRoutes函数，保持封装性
