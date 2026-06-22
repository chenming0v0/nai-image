# 部署测试清单

在正式部署前，请按照此清单检查配置。

## 前置检查

### 1. 环境检查

```bash
# 检查 Docker
docker --version
# 期望输出：Docker version 20.10.x 或更高

# 检查 Docker Compose
docker-compose --version
# 期望输出：Docker Compose version 2.x.x 或更高

# 检查 Git（可选）
git --version
```

### 2. 文件检查

确认以下文件存在：

- [ ] `docker-compose.yml` - Docker Compose 配置
- [ ] `backend/Dockerfile` - 后端 Dockerfile
- [ ] `deploy/Dockerfile` - 前端 Dockerfile
- [ ] `deploy/nginx.conf` - Nginx 配置
- [ ] `.env.example` - 环境变量模板

### 3. 配置检查

```bash
# 复制环境变量配置
cp .env.example .env

# 编辑 .env 文件
nano .env  # 或使用其他编辑器
```

必须配置的项：
- [ ] `UPSTREAM_BASE_URL` - 公益站 API 地址
- [ ] `UPSTREAM_API_KEY` - API 密钥

可选配置：
- [ ] `BACKEND_PORT` - 后端端口（默认 8787）
- [ ] `FRONTEND_PORT` - 前端端口（默认 8080）
- [ ] `DEFAULT_MODEL` - 默认模型

## 构建测试

### 1. 测试后端构建

```bash
cd backend
docker build -t nai-backend-test .
```

期望：构建成功，无错误信息

### 2. 测试前端构建

```bash
docker build -f deploy/Dockerfile -t nai-frontend-test .
```

期望：构建成功，无错误信息

### 3. 清理测试镜像

```bash
docker rmi nai-backend-test nai-frontend-test
```

## 启动测试

### 1. 启动服务

```bash
# 使用快速启动脚本
bash start.sh  # Linux/Mac
# 或
start.bat      # Windows

# 或直接使用 docker-compose
docker-compose up -d --build
```

### 2. 检查服务状态

```bash
# 查看容器状态
docker-compose ps

# 期望输出：
# NAME                  STATUS
# nai-image-backend     Up
# nai-image-frontend    Up
```

### 3. 查看日志

```bash
# 查看所有日志
docker-compose logs

# 实时查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend
```

### 4. 测试健康检查

```bash
# 测试后端健康检查
curl http://localhost:8787/api/health

# 期望输出：
# {"status":"ok","database":"connected"}

# 测试前端访问
curl -I http://localhost:8080

# 期望输出：
# HTTP/1.1 200 OK
```

### 5. 测试模型列表接口

```bash
# 通过后端获取模型列表
curl http://localhost:8787/api/models

# 期望输出：
# {"object":"list","data":[{"id":"nai-diffusion-4-5-full",...}]}
```

## 功能测试

### 1. 访问前端

浏览器打开：`http://localhost:8080`

- [ ] 页面正常加载
- [ ] 主题切换正常
- [ ] 设置页面可以打开

### 2. 配置测试

在设置页面：
- [ ] 可以看到默认的 API 配置
- [ ] 可以获取模型列表
- [ ] 可以修改配置并保存

### 3. 生成测试

尝试生成一张图片：
- [ ] 输入提示词
- [ ] 选择参数
- [ ] 点击生成
- [ ] 图片生成成功
- [ ] 可以在历史记录中看到

### 4. 历史记录测试

- [ ] 画廊页面正常显示
- [ ] 可以查看图片详情
- [ ] 可以下载图片
- [ ] 可以删除图片

## 性能测试

### 1. 资源使用

```bash
# 查看容器资源使用
docker stats

# 期望：
# CPU 使用率 < 50%（闲置时）
# 内存使用 < 1GB（后端）+ 100MB（前端）
```

### 2. 响应时间

```bash
# 测试 API 响应时间
time curl http://localhost:8787/api/health

# 期望：< 100ms
```

## 数据持久化测试

### 1. 生成测试数据

在前端生成 2-3 张图片

### 2. 重启服务

```bash
docker-compose restart
```

### 3. 检查数据

- [ ] 重启后历史记录仍然存在
- [ ] 数据库文件正常：`ls -lh backend/data/nai.db`

## 故障排查

### 后端无法启动

```bash
# 查看后端详细日志
docker-compose logs backend

# 常见问题：
# - 端口被占用：修改 .env 中的 BACKEND_PORT
# - 数据库权限：chmod 777 backend/data
# - 上游配置错误：检查 UPSTREAM_BASE_URL 和 UPSTREAM_API_KEY
```

### 前端无法访问

```bash
# 查看前端详细日志
docker-compose logs frontend

# 常见问题：
# - 端口被占用：修改 .env 中的 FRONTEND_PORT
# - 后端未启动：先启动后端
# - Nginx 配置错误：检查 deploy/nginx.conf
```

### 无法连接上游 API

```bash
# 在容器内测试网络
docker exec -it nai-image-backend wget -O- http://localhost:8787/api/models

# 检查上游配置
docker exec -it nai-image-backend env | grep UPSTREAM
```

## 清理

### 停止服务

```bash
docker-compose down
```

### 清理数据

```bash
# 仅清理容器和网络
docker-compose down

# 清理容器、网络和卷
docker-compose down -v

# 清理数据库（谨慎！）
rm -rf backend/data/
```

### 清理镜像

```bash
# 查看镜像
docker images | grep nai-image

# 删除镜像
docker rmi nai-image-backend nai-image-frontend
```

## 完成

如果所有测试项都通过，说明部署成功！

遇到问题请查看：
- 详细日志：`docker-compose logs -f`
- Docker 文档：[DOCKER.md](DOCKER.md)
- 后端文档：[backend/README.md](backend/README.md)
- API 文档：[API接入文档-20260527.md](API接入文档-20260527.md)
