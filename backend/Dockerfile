# ---- Build stage ----
FROM golang:1.23-alpine AS builder

WORKDIR /build

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# ---- Production stage ----
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite-libs

# 从构建阶段复制二进制文件
COPY --from=builder /build/server .

# 创建数据目录
RUN mkdir -p /app/data

# 环境变量
ENV PORT=8787
ENV DB_PATH=/app/data/nai.db
ENV GIN_MODE=release

EXPOSE 8787

CMD ["./server"]
