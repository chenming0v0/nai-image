# 开发环境快速启动指南

## 一键启动开发环境

### 前置要求

- Node.js 20+
- Go 1.23+
- npm 或 pnpm

### 第一次运行

```bash
# 1. 安装依赖（包含 concurrently）
npm install

# 2. 一键启动前端 + 后端
npm start
```

就这么简单！🎉

### 可用的 npm 命令

```bash
# 🚀 一键启动前后端（推荐）
npm start
# 或
npm run dev:all

# 🎨 只启动前端
npm run dev:frontend
# 或
npm run dev

# 🔧 只启动后端
npm run dev:backend

# 📦 构建前端
npm run build

# 🔨 构建后端
npm run build:backend

# 📦 构建前端 + 后端
npm run build:all

# 👀 预览生产构建
npm run preview

# 🧪 运行测试
npm test
npm run test:watch
```

## 启动后访问

运行 `npm start` 后：

- **前端**: http://localhost:5111
- **后端 API**: http://localhost:8787
- **后端健康检查**: http://localhost:8787/api/health

## 输出说明

启动后你会看到带颜色的输出：

```
[backend] | 2026/06/21 18:30:00 Server running on :8787
[frontend] | 
[frontend] |   VITE v6.3.2  ready in 1234 ms
[frontend] | 
[frontend] |   ➜  Local:   http://localhost:5173/
```

- **蓝色** = 后端日志
- **紫色** = 前端日志

## 停止服务

按 `Ctrl+C` 即可同时停止前后端服务。

## 环境变量配置（可选）

### 后端配置

如果需要配置上游 API，可以：

**方式一：环境变量**

```bash
# Windows PowerShell
$env:UPSTREAM_BASE_URL="https://api.example.com/v1"; npm start

# Linux/Mac
UPSTREAM_BASE_URL=https://api.example.com/v1 npm start
```

**方式二：创建 backend/.env 文件**

```env
UPSTREAM_BASE_URL=https://api.example.com/v1
UPSTREAM_API_KEY=your-key
DEFAULT_MODEL=nai-diffusion-4-5-full
```

### 前端配置

如果需要配置默认 API 地址，创建 `.env.local`：

```env
VITE_API_BASE_URL=http://localhost:8787
```

## 常见问题

### Q: 端口被占用怎么办？

**前端端口冲突**（5173）：

修改 `vite.config.ts`：
```ts
export default defineConfig({
  server: {
    port: 5174  // 改成其他端口
  }
})
```

**后端端口冲突**（8787）：

设置环境变量：
```bash
# Windows
$env:PORT=8788; npm start

# Linux/Mac
PORT=8788 npm start
```

### Q: Go 命令找不到？

确认已安装 Go：
```bash
go version
```

如果未安装，访问 https://go.dev/dl/ 下载安装。

### Q: 后端启动失败？

检查 Go 依赖：
```bash
cd backend
go mod download
go run ./cmd/server
```

### Q: 前端无法连接后端？

1. 确认后端已启动（检查 http://localhost:8787/api/health）
2. 检查浏览器控制台是否有 CORS 错误
3. 确认前端配置的 API 地址正确

### Q: 想单独调试某个服务？

```bash
# 只启动前端（适合调试前端）
npm run dev:frontend

# 只启动后端（适合调试后端）
npm run dev:backend
```

## 开发技巧

### 1. 使用热重载

- **前端**：修改代码后自动刷新（Vite HMR）
- **后端**：需要手动重启（或使用 air/nodemon）

### 2. 安装 Go 热重载工具（可选）

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 在 backend 目录运行
cd backend
air
```

然后修改 `package.json`：
```json
"dev:backend": "cd backend && air"
```

### 3. 查看实时日志

使用 `npm start` 时，前后端日志会同时显示。

如果日志太多，可以分开启动：
```bash
# 终端 1
npm run dev:backend

# 终端 2
npm run dev:frontend
```

### 4. API 调试

推荐使用：
- **Postman** / **Insomnia** - 测试后端 API
- **Browser DevTools** - 查看网络请求

健康检查端点：
```bash
curl http://localhost:8787/api/health
```

## 生产构建

### 构建所有

```bash
npm run build:all
```

构建产物：
- 前端：`dist/` 目录
- 后端：`dist/nai-backend` 或 `dist/nai-backend.exe`

### 运行生产构建

```bash
# 启动后端
./dist/nai-backend  # Linux/Mac
dist\nai-backend.exe  # Windows

# 预览前端（使用 Vite preview）
npm run preview
```

## 推荐的开发工作流

```bash
# 1. 启动开发环境
npm start

# 2. 打开编辑器
code .  # VS Code
# 或其他编辑器

# 3. 开始开发
# - 修改前端代码：src/ 目录
# - 修改后端代码：backend/ 目录

# 4. 测试
# - 前端：http://localhost:5173
# - 后端：http://localhost:8787

# 5. 提交代码
git add .
git commit -m "feat: xxx"
git push
```

## 团队协作建议

### 环境变量管理

不要提交 `.env` 文件到 Git！

```bash
# .gitignore 已包含
.env
.env.local
backend/.env
```

团队成员各自创建自己的 `.env` 文件。

### 数据库

开发时使用本地 SQLite：
- 位置：`backend/data/nai.db`
- 已在 `.gitignore` 中排除
- 每个开发者有独立的数据库

### 推荐的 IDE 设置

**VS Code 插件**：
- Go (官方)
- ESLint
- Prettier
- Tailwind CSS IntelliSense
- Error Lens

**WebStorm / GoLand**：
- 内置支持，开箱即用

## 更多帮助

- 前端问题：查看 [README.md](README.md)
- 后端问题：查看 [backend/README.md](backend/README.md)
- Docker 部署：查看 [DOCKER.md](DOCKER.md)
- API 文档：查看 [API接入文档-20260527.md](API接入文档-20260527.md)

---

**快速启动命令记忆口诀**：
```bash
npm install  # 第一次
npm start    # 每次开发
```

就这么简单！🚀
