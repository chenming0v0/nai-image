@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ==========================================
echo   NAI Image Client - Docker 部署
echo ==========================================
echo.

REM 检查 Docker
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误：未检测到 Docker，请先安装 Docker Desktop
    pause
    exit /b 1
)

REM 检查 Docker Compose
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误：未检测到 Docker Compose，请先安装
    pause
    exit /b 1
)

REM 检查 .env 文件
if not exist .env (
    echo ⚠️  未找到 .env 文件，正在创建...
    copy .env.example .env >nul
    echo ✅ 已创建 .env 文件，请编辑配置后重新运行
    echo.
    echo 需要配置的项：
    echo   - UPSTREAM_BASE_URL: 公益站 API 地址
    echo   - UPSTREAM_API_KEY: API 密钥
    echo.
    pause
    exit /b 0
)

REM 创建数据目录
if not exist backend\data mkdir backend\data

echo 📦 构建并启动服务...
docker-compose up -d --build

echo.
echo ⏳ 等待服务启动...
timeout /t 5 /nobreak >nul

REM 检查服务状态
docker-compose ps | findstr "Up" >nul
if errorlevel 1 (
    echo.
    echo ❌ 服务启动失败，请查看日志：
    echo    docker-compose logs
    pause
    exit /b 1
)

echo.
echo ==========================================
echo ✅ 服务启动成功！
echo ==========================================
echo.
echo 访问地址：
echo   🌐 前端界面: http://localhost:8080
echo   🔧 后端 API: http://localhost:8787
echo   ❤️  健康检查: http://localhost:8787/api/health
echo.
echo 常用命令：
echo   查看日志: docker-compose logs -f
echo   停止服务: docker-compose down
echo   重启服务: docker-compose restart
echo.
pause
