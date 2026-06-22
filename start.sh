#!/bin/bash

# NAI Image Client - 快速启动脚本

set -e

echo "=========================================="
echo "  NAI Image Client - Docker 部署"
echo "=========================================="
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ 错误：未检测到 Docker，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ 错误：未检测到 Docker Compose，请先安装"
    exit 1
fi

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "⚠️  未找到 .env 文件，正在创建..."
    cp .env.example .env
    echo "✅ 已创建 .env 文件，请编辑配置后重新运行"
    echo ""
    echo "需要配置的项："
    echo "  - UPSTREAM_BASE_URL: 公益站 API 地址"
    echo "  - UPSTREAM_API_KEY: API 密钥"
    echo ""
    exit 0
fi

# 创建数据目录
mkdir -p ./backend/data
chmod 777 ./backend/data

echo "📦 构建并启动服务..."
docker-compose up -d --build

echo ""
echo "⏳ 等待服务启动..."
sleep 5

# 检查服务状态
if docker-compose ps | grep -q "Up"; then
    echo ""
    echo "=========================================="
    echo "✅ 服务启动成功！"
    echo "=========================================="
    echo ""
    echo "访问地址："
    echo "  🌐 前端界面: http://localhost:8080"
    echo "  🔧 后端 API: http://localhost:8787"
    echo "  ❤️  健康检查: http://localhost:8787/api/health"
    echo ""
    echo "常用命令："
    echo "  查看日志: docker-compose logs -f"
    echo "  停止服务: docker-compose down"
    echo "  重启服务: docker-compose restart"
    echo ""
else
    echo ""
    echo "❌ 服务启动失败，请查看日志："
    echo "   docker-compose logs"
    exit 1
fi
