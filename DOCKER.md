# Docker 部署指南

本文档详细说明如何使用 Docker 和 Docker Compose 部署 NAI Image Client。

## 快速开始

### 1. 准备工作

确保已安装：
- Docker 20.10+
- Docker Compose 2.0+

检查版本：
```bash
docker --version
docker-compose --version
```

### 2. 克隆项目

```bash
git clone https://github.com/chenming0v0/nai-image.git
cd nai-image
```

### 3. 配置环境变量

复制环境变量模板：
```bash
cp .env.example .env
```

编辑 `.env` 文件，填写必要的配置：

```env
# 端口配置
BACKEND_PORT=8787
FRONTEND_PORT=8080

# 上游接口配置（必填）
UPSTREAM_BASE_URL=https://你的公益站地址/v1
UPSTREAM_API_KEY=你的API密钥

# 模型配置
DEFAULT_MODEL=nai-diffusion-4-5-full
```

### 4. 启动服务

```bash
# 构建并启动（首次运行或代码更新后）
docker-compose up -d --build

# 或者直接启动（使用已构建的镜像）
docker-compose up -d
```

### 5. 访问应用

打开浏览器访问：
- 前端界面：`http://localhost:8080`
- 后端 API：`http://localhost:8787`
- 健康检查：`http://localhost:8787/api/health`

## 服务架构

```
┌─────────────────┐
│   浏览器        │
└────────┬────────┘
         │ :8080
         ▼
┌─────────────────┐
│   Nginx         │ 前端服务（静态文件）
│   (Frontend)    │
└────────┬────────┘
         │ 内部网络
         ▼
┌─────────────────┐
│   Go/Fiber      │ 后端服务（API）
│   (Backend)     │ :8787
└────────┬────────┘
         │
         ├─► SQLite (本地存储)
         │
         └─► 上游 NewAPI (公益站)
```

## 常用命令

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看前端日志
docker-compose logs -f frontend

# 查看后端日志
docker-compose logs -f backend
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启单个服务
docker-compose restart backend
docker-compose restart frontend
```

### 停止服务

```bash
# 停止并删除容器
docker-compose down

# 停止并删除容器、网络、卷
docker-compose down -v
```

### 更新服务

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up -d --build
```

### 查看服务状态

```bash
docker-compose ps
```

## 数据持久化

### 数据存储位置

- **SQLite 数据库**：`./backend/data/nai.db`
- **任务历史**：存储在 SQLite 数据库中
- **图片元信息**：存储在 SQLite 数据库中

### 备份数据

```bash
# 备份数据库
cp ./backend/data/nai.db ./backup/nai-$(date +%Y%m%d).db

# 或使用 tar 打包
tar -czf nai-data-$(date +%Y%m%d).tar.gz ./backend/data/
```

### 恢复数据

```bash
# 停止服务
docker-compose down

# 恢复数据库
cp ./backup/nai-20260621.db ./backend/data/nai.db

# 启动服务
docker-compose up -d
```

## 环境变量说明

### 端口配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `FRONTEND_PORT` | `8080` | 前端服务端口 |
| `BACKEND_PORT` | `8787` | 后端 API 端口 |

### 上游接口配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `UPSTREAM_BASE_URL` | 空 | NewAPI 地址（必填） |
| `UPSTREAM_API_KEY` | 空 | API 密钥（必填） |
| `DEFAULT_MODEL` | `nai-diffusion-4-5-full` | 默认模型 |

### 高级配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `REQUEST_TIMEOUT_SECONDS` | `180` | 请求超时时间（秒） |
| `MAX_IMAGE_BYTES` | `536870912` | 最大图片大小（512MB） |

## 故障排查

### 1. 端口被占用

**错误信息**：
```
Error starting userland proxy: listen tcp4 0.0.0.0:8080: bind: address already in use
```

**解决方案**：
- 修改 `.env` 文件中的端口号
- 或停止占用端口的程序

### 2. 后端无法连接上游

**错误信息**：
```
upstream connection failed
```

**解决方案**：
- 检查 `UPSTREAM_BASE_URL` 是否正确
- 检查 `UPSTREAM_API_KEY` 是否有效
- 确认网络连接正常

### 3. 数据库文件权限问题

**错误信息**：
```
unable to open database file
```

**解决方案**：
```bash
# 创建数据目录
mkdir -p ./backend/data

# 修改权限
chmod 777 ./backend/data
```

### 4. 查看容器内部日志

```bash
# 进入容器
docker exec -it nai-image-backend sh

# 查看日志文件
ls -la /app/data/
```

## 性能优化

### 1. 限制资源使用

编辑 `docker-compose.yml`，添加资源限制：

```yaml
services:
  backend:
    # ... 其他配置
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### 2. 启用日志轮转

```yaml
services:
  backend:
    # ... 其他配置
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 安全建议

### 1. 使用 HTTPS

建议在生产环境中使用 Nginx 反向代理并配置 SSL：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api/ {
        proxy_pass http://localhost:8787;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. 限制访问

使用防火墙或 Docker 网络隔离：

```bash
# 只允许本地访问
FRONTEND_PORT=127.0.0.1:8080
BACKEND_PORT=127.0.0.1:8787
```

### 3. 定期更新

```bash
# 定期拉取最新代码
git pull

# 重新构建镜像
docker-compose build --no-cache

# 重启服务
docker-compose up -d
```

## 多实例部署

如果需要运行多个独立实例：

```bash
# 复制项目目录
cp -r nai-image nai-image-instance2
cd nai-image-instance2

# 修改 .env 中的端口
FRONTEND_PORT=8081
BACKEND_PORT=8788

# 启动
docker-compose up -d
```

## 监控和维护

### 健康检查

后端服务包含健康检查端点：

```bash
curl http://localhost:8787/api/health
```

响应示例：
```json
{
  "status": "ok",
  "database": "connected"
}
```

### 定期清理

```bash
# 清理未使用的镜像
docker image prune -a

# 清理未使用的容器
docker container prune

# 清理未使用的卷
docker volume prune
```

## 获取帮助

如果遇到问题：

1. 查看日志：`docker-compose logs -f`
2. 检查服务状态：`docker-compose ps`
3. 提交 Issue：https://github.com/chenming0v0/nai-image/issues
