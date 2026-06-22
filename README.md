<div align="center">

# 🎨 NAI Image Client

[![License](https://img.shields.io/badge/license-MIT-10b981?style=flat-square)](https://github.com/chenming0v0/nai-image/blob/main/LICENSE)
[![React](https://img.shields.io/badge/React-19-20232A?style=flat-square&logo=react&logoColor=61DAFB)](https://react.dev/)
[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)

**基于 NovelAI 公益站接口的图像生成客户端**

现代化的 Web UI + Go 后端，支持文生图、图生图、多角色生成、Vibe Transfer、历史记录本地存储。

</div>

---

## ✨ 项目特点

本项目是为 **NovelAI 公益站接口**（NewAPI 兼容格式）定制的个人客户端，具有以下特色：

### 🎯 核心优势
- **前后端分离架构**：React 19 前端 + Go 后端，稳定高效
- **完整的 NAI V4/V4.5 功能支持**：文生图、图生图、多角色、Vibe Transfer、Character Reference
- **本地数据管理**：前端 IndexedDB + 后端 SQLite 双重存储，支持历史记录管理
- **适合个人使用**：轻量部署，支持 Docker 一键启动

### 🎨 前端特性
- **现代化设计**：基于 Tailwind CSS 3，支持亮色/暗色主题
- **流畅动画**：Framer Motion 驱动的交互体验
- **多角色管理**：支持角色位置网格（5×5）、多角色组合生成
- **Vibe Transfer**：风格参考图支持，带 cache_id 复用机制
- **历史画廊**：瀑布流展示、收藏夹管理、批量下载
- **Chat 模式**：对话式生图，支持上下文记忆

### 🔧 后端特性
- **Go + Fiber**：高性能 HTTP 服务
- **SQLite 存储**：任务历史、图片元信息持久化
- **NewAPI 代理**：适配公益站的 OpenAI 兼容接口格式
- **RESTful API**：完整的任务管理、历史查询、图片存储接口

---

## 📸 界面预览

<details>
<summary><b>点击展开截图展示</b></summary>
<br>

<div align="center">
  <b>主界面</b><br>
  <img src="docs/images/example_pc_1.jpg" alt="主界面" />
</div>

<br>

<div align="center">
  <b>多角色生成</b><br>
  <img src="docs/images/example_pc_3.jpg" alt="多角色生成" />
</div>

<br>

<div align="center">
  <b>Chat 对话模式</b><br>
  <img src="docs/images/example_pc_4.jpg" alt="Chat 对话模式" />
</div>

</details>

---

## 🚀 快速开始

### 方式一：Docker Compose 部署（推荐）

最简单的部署方式，前后端一键启动：

```bash
# 克隆仓库
git clone https://github.com/chenming0v0/nai-image.git
cd nai-image

# 启动服务
docker-compose up -d
```

访问 `http://localhost:8080` 即可使用。

**环境变量配置**：

编辑 `docker-compose.yml` 或创建 `.env` 文件：

```env
# 后端配置
BACKEND_PORT=8787
UPSTREAM_BASE_URL=https://你的公益站地址
UPSTREAM_API_KEY=你的API密钥
DEFAULT_MODEL=nai-diffusion-4-5-full

# 前端配置
FRONTEND_PORT=8080
```

### 方式二：本地开发运行

**前置要求**：
- Node.js 20+
- Go 1.23+

**启动后端**：

```bash
cd backend
go run ./cmd/server
```

后端默认监听 `http://127.0.0.1:8787`。

**启动前端**：

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
```

前端开发服务器默认 `http://localhost:5173`。

### 方式三：仅前端部署（Vercel/Cloudflare）

如果只需要前端，可以部署到静态托管平台，通过前端直接调用公益站 API：

**Vercel 部署**：

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fchenming0v0%2Fnai-image)

**Cloudflare Workers 部署**：

```bash
npm run deploy:cf
```

---

## 📖 使用说明

### 1. 配置 API

首次使用需要在设置中配置：

- **API 地址**：公益站提供的 NewAPI 兼容地址（如 `https://example.com/v1`）
- **API Key**：公益站提供的密钥
- **模型选择**：通常使用 `nai-diffusion-4-5-full` 或 `nai-diffusion-4-5-curated`

### 2. 生成图片

**基础文生图**：
1. 输入英文提示词（prompt）
2. 可选：添加负面提示词（negative prompt）
3. 选择尺寸、步数等参数
4. 点击生成

**多角色生成**：
1. 点击"多角色"按钮
2. 为每个角色设置独立的提示词
3. 使用位置网格指定角色位置（A1-E5）
4. 生成多角色组合图片

**图生图（i2i）**：
1. 上传参考图
2. 调整 strength（变化强度）
3. 输入新的提示词
4. 生成

**Vibe Transfer**：
1. 上传 1-4 张风格参考图
2. 设置每张图的 info_extracted 和 strength
3. 生成带有参考风格的图片

### 3. 历史管理

- **画廊**：查看所有生成历史
- **收藏夹**：创建多个收藏夹分类管理
- **批量操作**：支持批量下载、删除
- **数据导出**：一键打包下载 ZIP

---

## 🛠️ API 接口说明

后端提供完整的 RESTful API，详见 [API接入文档-20260527.md](API接入文档-20260527.md)。

### 主要端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查 |
| GET/PUT | `/api/settings` | 上游配置管理 |
| GET | `/api/models` | 获取可用模型列表 |
| POST | `/api/generate` | 生成图片 |
| GET | `/api/tasks` | 获取历史任务 |
| GET | `/api/images/:id` | 获取图片 |

---

## 🐳 Docker 部署详解

### Docker Compose 配置

项目包含完整的 `docker-compose.yml`，支持：

- **前端服务**：Nginx 提供静态文件
- **后端服务**：Go API 服务
- **数据持久化**：SQLite 数据库挂载到宿主机
- **环境变量注入**：灵活配置上游接口

### 端口说明

| 服务 | 容器端口 | 宿主机端口（默认） |
|------|---------|-------------------|
| 前端 | 80 | 8080 |
| 后端 | 8787 | 8787 |

### 数据卷

```yaml
volumes:
  - ./backend/data:/app/data  # SQLite 数据库
```

历史记录和图片元信息存储在 `./backend/data/nai.db`。

---

## 💻 技术栈

### 前端

<div align="center">
  <a href="https://react.dev/"><img src="https://img.shields.io/badge/React_19-20232A?style=for-the-badge&logo=react&logoColor=61DAFB" alt="React 19" /></a>
  <a href="https://www.typescriptlang.org/"><img src="https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript" /></a>
  <a href="https://vite.dev/"><img src="https://img.shields.io/badge/Vite-B73BFE?style=for-the-badge&logo=vite&logoColor=FFD62E" alt="Vite" /></a>
  <a href="https://tailwindcss.com/"><img src="https://img.shields.io/badge/Tailwind_CSS_3-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white" alt="Tailwind CSS 3" /></a>
  <a href="https://www.framer.com/motion/"><img src="https://img.shields.io/badge/Framer_Motion-0055FF?style=for-the-badge&logo=framer&logoColor=white" alt="Framer Motion" /></a>
  <a href="https://zustand.docs.pmnd.rs/"><img src="https://img.shields.io/badge/Zustand-764ABC?style=for-the-badge&logo=react&logoColor=white" alt="Zustand" /></a>
</div>

### 后端

<div align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go_1.23-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" /></a>
  <a href="https://gofiber.io/"><img src="https://img.shields.io/badge/Fiber-00ACD7?style=for-the-badge&logo=go&logoColor=white" alt="Fiber" /></a>
  <a href="https://www.sqlite.org/"><img src="https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite" /></a>
</div>

---

## 📝 开发计划

- [x] 基础文生图功能
- [x] 多角色生成
- [x] 图生图、Vibe Transfer
- [x] Go 后端 API
- [x] Docker 部署支持
- [ ] 提示词历史记录
- [ ] 更多模型参数支持

---

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

**项目来源**：
- 前端界面基于 [CookSleep/gpt_image_playground](https://github.com/CookSleep/gpt_image_playground) 修改
- 后端为原创开发
- 适配 NovelAI 公益站接口

---

## 🙏 致谢

- [CookSleep/gpt_image_playground](https://github.com/CookSleep/gpt_image_playground) - 前端界面参考
- NovelAI 公益站 - API 接口支持
- 所有为开源社区做出贡献的开发者

---

## ⚠️ 免责声明

本项目仅供学习交流使用，请遵守相关服务条款和法律法规。使用本项目产生的任何后果由使用者自行承担。

---

<div align="center">

**如果这个项目对你有帮助，欢迎 Star ⭐**

</div>
