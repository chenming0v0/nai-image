# nai-image 后端

Go 服务：为前端提供结构化 NAI 绘图 API，并代理上游 NewAPI（OpenAI 兼容 `chat/completions`）。

## 运行

```powershell
cd backend
go run ./cmd/server
```

默认监听 `http://127.0.0.1:8787`，数据库 `./data/nai.db`。

## 环境变量

| 变量 | 默认 | 说明 |
|------|------|------|
| `PORT` | `8787` | HTTP 端口 |
| `DB_PATH` | `./data/nai.db` | SQLite 路径 |
| `UPSTREAM_BASE_URL` | 空 | 上游 NewAPI 地址（也可通过 API 设置） |
| `UPSTREAM_API_KEY` | 空 | 上游密钥 |
| `DEFAULT_MODEL` | `nai-diffusion-4-5-full` | 默认模型 |
| `REQUEST_TIMEOUT_SECONDS` | `180` | 上游请求超时 |
| `MAX_IMAGE_BYTES` | `536870912` | 请求体大小上限 |

## HTTP API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查 |
| GET/PUT | `/api/settings` | 上游配置（密钥仅返回是否已设置） |
| GET | `/api/models` | 代理上游 `/v1/models` |
| POST | `/api/generate` | 结构化绘图（见 `internal/models.DrawRequest`） |
| GET | `/api/tasks` | 历史列表 |
| GET/DELETE | `/api/tasks/:id` | 单条任务 |
| DELETE | `/api/tasks` | 清空历史 |
| GET | `/api/images/:id` | 输出图二进制 |
| GET | `/api/images/:id/meta` | 图片元信息 |

错误响应：`{"error":"..."}`。

绘图参数规则见仓库根目录 `API接入文档-20260527.md`。
