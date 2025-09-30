# Go Backend Template

<div align="center">

**[English](#english) | [中文](#chinese)**

一个基于领域驱动设计（DDD）的现代化Golang后端开发模版

A modern Golang backend development template based on Domain-Driven Design (DDD)

[![Go Version](https://img.shields.io/badge/Go-1.25.1-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![DDD](https://img.shields.io/badge/architecture-DDD-brightgreen.svg)](https://en.wikipedia.org/wiki/Domain-driven_design)

</div>

---

<a name="english"></a>

## English

### 📖 Overview

This is a production-ready Golang backend template built with **Domain-Driven Design (DDD)** architecture principles. It provides a solid foundation for building scalable, maintainable, and testable web applications.

### ✨ Features

- 🏗️ **DDD Architecture** - Clean separation of concerns with Domain, Application, Infrastructure, and Interface layers
- 💉 **Dig DI** - Type-safe dependency injection using uber-go/dig
- 🔐 **OAuth2 Integration** - Support for GitHub, Google authentication
- 🔑 **JWT Authentication** - Secure token-based authentication
- 💾 **PostgreSQL** - Robust relational database with GORM
- 🚀 **Redis Cache** - High-performance caching layer
- 📦 **Object Storage** - Support for MinIO and Tencent COS
- 🤖 **OpenAI Integration** - Ready-to-use AI capabilities
- 📝 **Swagger Documentation** - Auto-generated API documentation
- ⚡ **High Performance** - Built with Fiber web framework
- 🔄 **Cron Jobs** - Scheduled task support

### 🏛️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Interface Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ HTTP Handler │  │   Middleware  │  │    Router    │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────┐
│                  Application Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  App Service │  │     DTO      │  │   Commands   │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────┐
│                    Domain Layer                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Entity     │  │ Value Object │  │   Repository │  │
│  └──────────────┘  └──────────────┘  │  Interface   │  │
│  ┌──────────────┐  ┌──────────────┐  └──────────────┘  │
│  │   Aggregate  │  │Domain Service│                     │
│  └──────────────┘  └──────────────┘                     │
└─────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────┐
│                Infrastructure Layer                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Repository   │  │ OAuth2 Impl  │  │   Database   │  │
│  │ Implementation│  └──────────────┘  └──────────────┘  │
│  └──────────────┘  ┌──────────────┐  ┌──────────────┐  │
│  │ Dig Container│  │    Cache     │  │   Storage    │  │
│  │   (DI)       │  └──────────────┘  └──────────────┘  │
│  └──────────────┘                                        │
└─────────────────────────────────────────────────────────┘
```

### 📁 Project Structure

```
.
├── cmd/                      # Command line tools
├── internal/
│   ├── application/          # Application Layer
│   │   ├── auth/            # Auth application services
│   │   └── user/            # User application services
│   ├── domain/              # Domain Layer
│   │   ├── auth/            # Auth domain models
│   │   └── user/            # User aggregate
│   ├── infrastructure/       # Infrastructure Layer
│   │   ├── container/       # Dependency injection
│   │   ├── oauth2/          # OAuth2 providers
│   │   └── persistence/     # Repository implementations
│   ├── interfaces/          # Interface Layer
│   │   └── http/            # HTTP handlers
│   ├── middleware/          # HTTP middlewares
│   ├── router/              # Route definitions
│   ├── protocol/            # API protocols & DTOs
│   ├── resource/            # External resources
│   │   ├── cache/           # Redis cache
│   │   ├── database/        # Database models & DAOs
│   │   ├── llm/             # LLM integration
│   │   └── storage/         # Object storage
│   ├── config/              # Configuration
│   ├── logger/              # Logging
│   └── util/                # Utilities
├── docker/                   # Docker files
├── env/                      # Environment configs
└── docs/                     # Swagger documentation
```

### 🚀 Quick Start

#### Prerequisites

- Go 1.24+
- PostgreSQL 12+
- Redis 6+
- (Optional) MinIO / Tencent COS
- (Optional) OpenAI API Key

#### Installation

1. **Clone the repository**
```bash
git clone https://github.com/hcd233/go-backend-tmpl.git
cd go-backend-tmpl
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment**
```bash
cp env/api.env.template env/api.env
# Edit env/api.env with your configuration
```

4. **Initialize database**
```bash
go run main.go database migrate
```

5. **Run the server**
```bash
go run main.go server start --port 8080
```

#### Using Docker Compose

```bash
# Start all services
docker-compose -f docker/docker-compose.yml up -d

# Without middleware services
docker-compose -f docker/docker-compose-without-middlewares.yml up -d
```

### 📚 API Documentation

Once the server is running, visit:
- Swagger UI: `http://localhost:8080/swagger/`

### 🔑 Environment Variables

Key environment variables (see `env/api.env.template` for complete list):

```bash
# Server
READ_TIMEOUT=10
WRITE_TIMEOUT=10

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DATABASE=backend

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_ACCESS_TOKEN_SECRET=your-access-secret
JWT_ACCESS_TOKEN_EXPIRED=1h
JWT_REFRESH_TOKEN_SECRET=your-refresh-secret
JWT_REFRESH_TOKEN_EXPIRED=168h

# OAuth2
OAUTH2_GITHUB_CLIENT_ID=your-github-client-id
OAUTH2_GITHUB_CLIENT_SECRET=your-github-client-secret
OAUTH2_GOOGLE_CLIENT_ID=your-google-client-id
OAUTH2_GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### 🏗️ DDD Concepts

#### Domain Layer (`internal/domain`)
- **Entities**: Business objects with unique identity
- **Value Objects**: Immutable objects defined by their attributes
- **Aggregates**: Cluster of domain objects treated as a single unit
- **Domain Services**: Operations that don't naturally fit within entities
- **Repository Interfaces**: Abstract data access

#### Application Layer (`internal/application`)
- **Application Services**: Orchestrate use cases and business workflows
- **DTOs**: Data Transfer Objects for communication between layers
- **Commands/Queries**: CQRS-inspired message objects

#### Infrastructure Layer (`internal/infrastructure`)
- **Repository Implementations**: Concrete data access implementations
- **External Service Adapters**: OAuth2, storage, etc.
- **Dependency Injection Container**: Wire up dependencies

#### Interface Layer (`internal/interfaces`)
- **HTTP Handlers**: Handle HTTP requests/responses
- **Protocol Conversion**: Convert between HTTP and application DTOs

### 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/domain/user
```

### 📦 Build & Deploy

```bash
# Build binary
go build -o bin/server main.go

# Run binary
./bin/server server start

# Build Docker image
docker build -f docker/dockerfile -t go-backend-tmpl:latest .
```

### 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

<a name="chinese"></a>

## 中文

### 📖 概述

这是一个基于**领域驱动设计（DDD）**架构原则构建的生产就绪的Golang后端模版。它为构建可扩展、可维护和可测试的Web应用程序提供了坚实的基础。

### ✨ 特性

- 🏗️ **DDD架构** - 通过领域层、应用层、基础设施层和接口层清晰分离关注点
- 💉 **Dig依赖注入** - 使用uber-go/dig实现类型安全的依赖注入
- 🔐 **OAuth2集成** - 支持GitHub、Google认证
- 🔑 **JWT认证** - 安全的基于令牌的认证
- 💾 **PostgreSQL** - 使用GORM的健壮关系型数据库
- 🚀 **Redis缓存** - 高性能缓存层
- 📦 **对象存储** - 支持MinIO和腾讯云COS
- 🤖 **OpenAI集成** - 开箱即用的AI能力
- 📝 **Swagger文档** - 自动生成的API文档
- ⚡ **高性能** - 使用Fiber Web框架构建
- 🔄 **定时任务** - 支持计划任务

### 🏛️ 架构

```
┌────────────────────────────────────────────────────────┐
│                      接口层                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  HTTP处理器  │  │   中间件     │  │    路由器    │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────────────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│                      应用层                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  应用服务    │  │     DTO      │  │    命令      │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────────────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│                      领域层                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │    实体      │  │   值对象     │  │   仓储接口   │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│  ┌──────────────┐  ┌──────────────┐                    │
│  │    聚合      │  │  领域服务    │                    │
│  └──────────────┘  └──────────────┘                    │
└────────────────────────────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│                    基础设施层                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  仓储实现    │  │ OAuth2实现   │  │   数据库     │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Dig依赖容器  │  │    缓存      │  │   存储       │  │
│  │   (DI)       │  └──────────────┘  └──────────────┘  │
│  └──────────────┘                                       │
└────────────────────────────────────────────────────────┘
```

### 📁 项目结构

```
.
├── cmd/                      # 命令行工具
├── internal/
│   ├── application/          # 应用层
│   │   ├── auth/            # 认证应用服务
│   │   └── user/            # 用户应用服务
│   ├── domain/              # 领域层
│   │   ├── auth/            # 认证领域模型
│   │   └── user/            # 用户聚合
│   ├── infrastructure/       # 基础设施层
│   │   ├── container/       # 依赖注入
│   │   ├── oauth2/          # OAuth2提供商
│   │   └── persistence/     # 仓储实现
│   ├── interfaces/          # 接口层
│   │   └── http/            # HTTP处理器
│   ├── middleware/          # HTTP中间件
│   ├── router/              # 路由定义
│   ├── protocol/            # API协议和DTO
│   ├── resource/            # 外部资源
│   │   ├── cache/           # Redis缓存
│   │   ├── database/        # 数据库模型和DAO
│   │   ├── llm/             # LLM集成
│   │   └── storage/         # 对象存储
│   ├── config/              # 配置
│   ├── logger/              # 日志
│   └── util/                # 工具
├── docker/                   # Docker文件
├── env/                      # 环境配置
└── docs/                     # Swagger文档
```

### 🚀 快速开始

#### 前置要求

- Go 1.24+
- PostgreSQL 12+
- Redis 6+
- （可选）MinIO / 腾讯云COS
- （可选）OpenAI API密钥

#### 安装

1. **克隆仓库**
```bash
git clone https://github.com/hcd233/go-backend-tmpl.git
cd go-backend-tmpl
```

2. **安装依赖**
```bash
go mod download
```

3. **配置环境**
```bash
cp env/api.env.template env/api.env
# 编辑env/api.env配置你的参数
```

4. **初始化数据库**
```bash
go run main.go database migrate
```

5. **启动服务器**
```bash
go run main.go server start --port 8080
```

#### 使用Docker Compose

```bash
# 启动所有服务
docker-compose -f docker/docker-compose.yml up -d

# 不包含中间件服务
docker-compose -f docker/docker-compose-without-middlewares.yml up -d
```

### 📚 API文档

服务器运行后，访问：
- Swagger UI: `http://localhost:8080/swagger/`

### 🔑 环境变量

关键环境变量（完整列表见`env/api.env.template`）：

```bash
# 服务器
READ_TIMEOUT=10
WRITE_TIMEOUT=10

# 数据库
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DATABASE=backend

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_ACCESS_TOKEN_SECRET=你的访问密钥
JWT_ACCESS_TOKEN_EXPIRED=1h
JWT_REFRESH_TOKEN_SECRET=你的刷新密钥
JWT_REFRESH_TOKEN_EXPIRED=168h

# OAuth2
OAUTH2_GITHUB_CLIENT_ID=你的github客户端ID
OAUTH2_GITHUB_CLIENT_SECRET=你的github客户端密钥
OAUTH2_GOOGLE_CLIENT_ID=你的google客户端ID
OAUTH2_GOOGLE_CLIENT_SECRET=你的google客户端密钥
```

### 🏗️ DDD概念

#### 领域层 (`internal/domain`)
- **实体（Entities）**: 具有唯一标识的业务对象
- **值对象（Value Objects）**: 由属性定义的不可变对象
- **聚合（Aggregates）**: 作为单一单元处理的领域对象集群
- **领域服务（Domain Services）**: 不自然适合实体的操作
- **仓储接口（Repository Interfaces）**: 抽象数据访问

#### 应用层 (`internal/application`)
- **应用服务（Application Services）**: 编排用例和业务工作流
- **DTO**: 层间通信的数据传输对象
- **命令/查询（Commands/Queries）**: CQRS风格的消息对象

#### 基础设施层 (`internal/infrastructure`)
- **仓储实现（Repository Implementations）**: 具体的数据访问实现
- **外部服务适配器（External Service Adapters）**: OAuth2、存储等
- **依赖注入容器（Dependency Injection Container）**: 连接依赖关系

#### 接口层 (`internal/interfaces`)
- **HTTP处理器（HTTP Handlers）**: 处理HTTP请求/响应
- **协议转换（Protocol Conversion）**: 在HTTP和应用DTO之间转换

### 🧪 测试

```bash
# 运行所有测试
go test ./...

# 带覆盖率运行
go test -cover ./...

# 运行特定包
go test ./internal/domain/user
```

### 📦 构建和部署

```bash
# 构建二进制文件
go build -o bin/server main.go

# 运行二进制文件
./bin/server server start

# 构建Docker镜像
docker build -f docker/dockerfile -t go-backend-tmpl:latest .
```

### 🤝 贡献

欢迎贡献！请随时提交Pull Request。

### 📄 许可证

本项目采用Apache License 2.0许可证 - 详见[LICENSE](LICENSE)文件。

---

<div align="center">

Made with ❤️ by [hcd233](https://github.com/hcd233)

</div>
