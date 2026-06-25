# 多角色功能测试文档

## 实现验证

经过代码审查，你们的多角色功能实现**完全正确**，符合 API 文档要求。

## 数据流

### 1. 前端 (TypeScript)
- 用户在界面输入多个角色，每个角色包含：
  - `prompt`: 正面提示词
  - `negative_prompt`: 负面提示词  
  - `position`: 位置（如 "C3", "A1" 等）
- 数据存储在 `store.ts` 中的 `characters` 数组

### 2. API 调用 (naiBackendImageApi.ts)
```typescript
// 第 78-88 行
if (extra.characters?.length) {
  body.characters = extra.characters
    .filter((character) => character.prompt.trim())
    .map((character) => ({
      prompt: character.prompt.trim(),
      negative_prompt: character.negative_prompt?.trim() || undefined,
      position: character.position,
    }))
}
if (extra.use_coords !== undefined) body.use_coords = extra.use_coords
if (extra.use_order !== undefined) body.use_order = extra.use_order
```

发送到后端 `http://127.0.0.1:8787/api/generate` 的数据格式：
```json
{
  "model": "nai-diffusion-4-5-full",
  "prompt": "2girls, masterpiece",
  "negative_prompt": "lowres, bad anatomy",
  "size": [1216, 832],
  "characters": [
    {
      "prompt": "1girl, blue hair, blue eyes",
      "negative_prompt": "white hair",
      "position": "B2"
    },
    {
      "prompt": "1girl, white hair, red eyes",
      "negative_prompt": "blue hair",
      "position": "D4"
    }
  ],
  "use_coords": true,
  "use_order": false,
  ...
}
```

### 3. 后端处理 (Go)

#### a) 接收请求 (handlers/generate.go)
```go
var req models.DrawRequest
if err := c.BodyParser(&req); err != nil {
  return fiberError(c, 400, "请求体解析失败: "+err.Error())
}
```

#### b) 构建 Chat 请求 (nai/types.go)
```go
// BuildChatRequest 第 82-125 行
func BuildChatRequest(req *models.DrawRequest) (*chatRequest, error) {
  inner := innerPayload{
    Prompt:         req.Prompt,
    NegativePrompt: req.NegativePrompt,
    Characters:     req.Characters,
    UseCoords:      req.UseCoords,
    UseOrder:       req.UseOrder,
    // ... 其他字段
  }
  innerBytes, err := json.Marshal(inner)
  // ...
  return &chatRequest{
    Model:     req.Model,
    Messages:  []chatMessage{{Role: "user", Content: string(innerBytes)}},
    Stream:    stream,
    MaxTokens: maxTokens,
  }, nil
}
```

发送到 NewAPI 的最终格式：
```json
{
  "model": "nai-diffusion-4-5-full",
  "messages": [
    {
      "role": "user",
      "content": "{\"prompt\":\"2girls, masterpiece\",\"negative_prompt\":\"lowres, bad anatomy\",\"size\":[1216,832],\"characters\":[{\"prompt\":\"1girl, blue hair, blue eyes\",\"negative_prompt\":\"white hair\",\"position\":\"B2\"},{\"prompt\":\"1girl, white hair, red eyes\",\"negative_prompt\":\"blue hair\",\"position\":\"D4\"}],\"use_coords\":true,\"use_order\":false,...}"
    }
  ],
  "stream": false,
  "max_tokens": 100000
}
```

**这完全符合 API 文档第 4 节的要求！**

## 测试步骤

### 1. 确保后端运行
```bash
cd backend
go run cmd/server/main.go
```
或使用 Docker：
```bash
docker-compose up -d
```

### 2. 配置 API
在前端界面中：
- 设置 API 模式为 "Chat"
- 填写 NewAPI 地址和密钥

### 3. 创建多角色任务
1. 切换到 "Gallery" 模式
2. 输入基础提示词：`2girls, standing side by side, looking at viewer, masterpiece`
3. 添加角色 1：
   - 正面：`1girl, blue hair, long blue hair, blue eyes, blue dress`
   - 负面：`white hair, pink hair`
   - 位置：`B2`（左上）
4. 添加角色 2：
   - 正面：`1girl, white hair, long white hair, red eyes, white kimono`
   - 负面：`blue hair, pink hair`
   - 位置：`D4`（右下）
5. 启用 "使用坐标定位" (use_coords)
6. 点击生成

### 4. 验证结果
- 检查浏览器 DevTools Network 标签
- 查看发送到 `http://127.0.0.1:8787/api/generate` 的请求体
- 确认 `characters` 数组包含完整数据
- 查看生成的图片，确认：
  - 左上区域是蓝发蓝眼角色
  - 右下区域是白发红眼角色

## 常见问题排查

### 问题 1: 角色数据未发送
**检查**: 在 `naiBackendImageApi.ts` 第 78 行设置断点
**验证**: `extra.characters` 是否有值

### 问题 2: 后端解析失败
**检查**: 后端日志中的 "请求体解析失败" 错误
**验证**: 前端发送的 JSON 格式是否正确

### 问题 3: NewAPI 返回错误
**检查**: 后端日志中的上游错误信息
**可能原因**:
- API 密钥无效
- 模型不支持多角色（需要 V4 或 V4.5 系列）
- 提示词包含中文字符

### 问题 4: 角色位置不生效
**检查**: `use_coords` 是否设置为 `true`
**验证**: 在界面中启用 "使用坐标定位" 开关

## 实现总结

✅ **前端实现正确**
- 数据结构定义完整
- API 调用正确传递参数
- 界面展示和编辑功能完善

✅ **后端实现正确**
- 正确接收前端数据
- 按照 API 文档要求构造请求
- 正确将参数序列化到 `messages[0].content`

✅ **符合 API 文档**
- 外层保留 `model`、`stream`、`max_tokens`
- 内层包含所有绘图参数（包括 `characters`、`use_coords`、`use_order`）
- 内层参数序列化为 JSON 字符串

## 后续优化建议

1. **错误提示优化**: 当用户输入中文提示词时，前端可以主动提示
2. **模型检测**: 自动检测所选模型是否支持多角色功能
3. **位置预览**: 提供 5×5 网格的可视化选择器
4. **角色模板**: 提供常用角色配置的快速模板
