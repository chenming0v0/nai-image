# nai-image 后端

Go 服务：为前端提供结构化 NAI 绘图 API，并代理上游 NewAPI（OpenAI 兼容 `chat/completions`）。

## 功能特性

- **RESTful API**：完整的任务管理、历史查询接口
- **SQLite 存储**：任务历史、图片元信息持久化
- **NewAPI 代理**：适配公益站的 OpenAI 兼容接口格式
- **结构化参数**：简化前端调用，自动转换为 NAI 格式

## 运行

### 本地开发

```bash
cd backend
go run ./cmd/server
```

默认监听 `http://127.0.0.1:8787`，数据库 `./data/nai.db`。

### 生产构建

```bash
cd backend
go build -o nai-backend ./cmd/server
./nai-backend
```

### Docker 运行

```bash
docker build -t nai-backend .
docker run -p 8787:8787 \
  -v $(pwd)/data:/app/data \
  -e UPSTREAM_BASE_URL=https://api.example.com/v1 \
  -e UPSTREAM_API_KEY=your-key \
  nai-backend
```

## 环境变量

| 变量 | 默认 | 说明 |
|------|------|------|
| `PORT` | `8787` | HTTP 端口 |
| `DB_PATH` | `./data/nai.db` | SQLite 路径 |
| `UPSTREAM_BASE_URL` | 空 | 上游 NewAPI 地址（也可通过 API 设置） |
| `UPSTREAM_API_KEY` | 空 | 上游密钥（也可通过 API 设置） |
| `DEFAULT_MODEL` | `nai-diffusion-4-5-full` | 默认模型 |
| `REQUEST_TIMEOUT_SECONDS` | `180` | 上游请求超时 |
| `MAX_IMAGE_BYTES` | `536870912` | 请求体大小上限（512MB） |

**注意**：API 配置优先级：数据库存储 > 环境变量

## HTTP API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查 |
| GET/PUT | `/api/settings` | 上游配置（密钥仅返回是否已设置） |
| GET | `/api/models` | 代理上游 `/v1/models` |
| POST | `/api/generate` | 结构化绘图（见 `internal/models.DrawRequest`） |
| GET | `/api/tasks` | 历史列表（支持分页、筛选） |
| GET | `/api/tasks/:id` | 单条任务详情 |
| DELETE | `/api/tasks/:id` | 删除单条任务 |
| DELETE | `/api/tasks` | 清空历史 |
| GET | `/api/images/:id` | 输出图二进制 |
| GET | `/api/images/:id/meta` | 图片元信息 |

错误响应：`{"error":"..."}`。

## API 示例

### 获取模型列表

```bash
curl http://localhost:8787/api/models
```

### 生成图片

```bash
curl -X POST http://localhost:8787/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "1girl, solo, masterpiece, best quality",
    "negative_prompt": "lowres, bad anatomy",
    "size": [832, 1216],
    "steps": 23,
    "model": "nai-diffusion-4-5-full"
  }'
```

### 查看历史任务

```bash
curl http://localhost:8787/api/tasks?limit=10&offset=0
```

## 数据存储

### 数据库结构

- **tasks** 表：任务历史
  - id, model, prompt, negative_prompt, parameters, status, created_at, updated_at
- **images** 表：生成的图片
  - id, task_id, image_data, format, size, seed, created_at

### 数据库位置

默认：`./data/nai.db`

可通过环境变量 `DB_PATH` 修改。

## 开发

### 目录结构

```
backend/
├── cmd/
│   └── server/          # 主程序入口
├── internal/
│   ├── api/             # HTTP 处理器
│   ├── database/        # 数据库操作
│   ├── models/          # 数据模型
│   └── proxy/           # 上游代理逻辑
├── go.mod
└── go.sum
```

### 添加新接口

1. 在 `internal/models/` 定义请求/响应结构
2. 在 `internal/api/` 添加处理器
3. 在 `cmd/server/main.go` 注册路由

## 绘图参数规则

详细的绘图参数说明见仓库根目录 [API接入文档-20260527.md](../API接入文档-20260527.md)。

支持的主要参数：
- `prompt`, `negative_prompt`：提示词（必须英文）
- `size`：尺寸数组 `[width, height]`
- `steps`：迭代步数（最大 28）
- `scale`：引导强度
- `sampler`：采样器
- `seed`：随机种子
- `characters`：多角色定义
- `i2i`：图生图
- `inpaint`：局部重绘
- `controlnet`：Vibe Transfer
- `character_references`：角色参考

## 测试

```bash
cd backend
go test ./...
```

## 许可证

MIT License
