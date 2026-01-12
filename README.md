# Demo App

一个简单的 Go Web 应用，用于测试 CI/CD 平台的完整流程。

## 功能

- 健康检查接口 `/health`
- 版本信息接口 `/version`
- 简单的 API 接口 `/api/hello`

## 本地运行

```bash
go run main.go
```

服务将在 `http://localhost:8000` 启动。

## Docker 构建

```bash
docker build -t demo-app .
docker run -p 8000:8000 demo-app
```

## API 接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/version` | GET | 版本信息 |
| `/api/hello` | GET | Hello World |
| `/api/hello?name=xxx` | GET | 个性化问候 |

## CI/CD 流程

1. 推送代码到 GitHub
2. 在 CI/CD 平台中创建项目并关联仓库
3. 平台自动分析代码，生成推荐的流水线配置
4. 在平台 UI 中配置流水线（无需修改代码）
5. 配置 Webhook 触发器
6. 代码变更时自动触发构建和部署
