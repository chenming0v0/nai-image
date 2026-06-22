# 更新日志

本项目基于 [CookSleep/gpt_image_playground](https://github.com/CookSleep/gpt_image_playground) 修改，适配 NovelAI 公益站接口。

## [1.0.0] - 2026-06-21

### 新增
- ✨ Go 后端服务（Fiber + SQLite）
- ✨ Docker Compose 一键部署支持
- ✨ 完整的 NAI V4/V4.5 接口适配
- ✨ 多角色生成支持（5×5 位置网格）
- ✨ Vibe Transfer 功能
- ✨ Character Reference 支持
- ✨ 图生图（i2i）和局部重绘（inpaint）
- ✨ 后端 RESTful API
- ✨ SQLite 数据持久化
- ✨ 任务历史管理

### 修改
- 🔧 前端适配后端 API 调用
- 🔧 Chat 模式走 NAI 后端
- 🔧 设置页支持从后端拉取模型列表
- 🔧 更新项目文档和说明

### 原项目特性保留
- ✅ React 19 + TypeScript + Vite
- ✅ Tailwind CSS 3 + Framer Motion
- ✅ 画廊瀑布流展示
- ✅ 收藏夹管理
- ✅ IndexedDB 本地存储
- ✅ 批量操作
- ✅ 亮色/暗色主题

## 原项目来源

- **原作者**: [CookSleep](https://github.com/CookSleep)
- **原项目**: [gpt_image_playground](https://github.com/CookSleep/gpt_image_playground)
- **许可证**: MIT License

## 主要变更

### 架构变化
- 从纯前端应用改为前后端分离架构
- 新增 Go 后端处理 API 代理和数据存储
- 支持 Docker 容器化部署

### 接口适配
- 从 OpenAI GPT-Image API 改为 NovelAI API
- 适配 NewAPI 兼容格式（通过 chat/completions 调用）
- 支持完整的 NAI 参数（多角色、Vibe Transfer 等）

### 数据存储
- 新增后端 SQLite 数据库
- 前端保留 IndexedDB 本地存储
- 双重存储机制提供更好的数据管理

## 未来计划

- [ ] OSS 云端数据同步
- [ ] 提示词历史记录
- [ ] 更多 NAI 模型参数支持
- [ ] 批量生成队列
- [ ] 图片编辑功能增强
- [ ] API 密钥管理优化

## 致谢

感谢原作者 [CookSleep](https://github.com/CookSleep) 提供优秀的前端界面基础。
