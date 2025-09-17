# 部署模型 (Deployment)

<cite>
**本文档引用文件**  
- [deployment.go](file://backend/internal/model/deployment.go)
- [user.go](file://backend/internal/model/user.go)
- [server.go](file://backend/internal/model/server.go)
- [init.sql](file://scripts/init.sql)
- [model.go](file://backend/internal/model/model.go)
- [types.go](file://backend/internal/api/types.go)
</cite>

## 目录
1. [引言](#引言)
2. [Deployment结构体字段详解](#deployment结构体字段详解)
3. [GORM关联关系实现机制](#gorm关联关系实现机制)
4. [数据库表设计与状态流转](#数据库表设计与状态流转)
5. [GORM操作代码示例](#gorm操作代码示例)
6. [大数据量下的分页与归档策略](#大数据量下的分页与归档策略)
7. [结论](#结论)

## 引言
本文档详细描述了自动化运维平台中的部署模型（Deployment）的设计与实现。该模型用于管理应用部署任务，记录部署过程中的关键信息，支持完整的部署生命周期管理。文档涵盖结构体定义、数据库映射、关联关系、状态设计及实际操作示例，旨在为开发者提供全面的技术参考。

## Deployment结构体字段详解

`Deployment` 结构体定义了部署任务的核心属性，每个字段均有明确的业务含义和数据约束：

- **ID**: 主键，唯一标识一个部署任务，自增整数。
- **Name**: 部署名称，最大长度100字符，不可为空，用于用户识别。
- **ServerID**: 关联的目标服务器ID，建立与`Server`模型的外键关系，带索引以优化查询性能。
- **Server**: 服务器对象，通过`foreignKey:ServerID`与`Server`模型关联，JSON序列化时可选输出。
- **Repository**: 代码仓库地址，最大长度200字符，用于指定部署源。
- **Branch**: 代码分支，默认值为`main`，最大长度50字符。
- **Path**: 部署目标路径，最大长度200字符，指定服务器上的部署目录。
- **Script**: 部署脚本，使用`text`类型存储，可容纳复杂部署逻辑。
- **Status**: 部署状态，整数类型，默认值为0，表示“待部署”。状态值映射为：0=待部署，1=部署中，2=部署成功，3=部署失败。
- **CreatedBy**: 触发部署的用户ID，建立与`User`模型的外键关系，带索引以便按用户查询。
- **User**: 用户对象，通过`foreignKey:CreatedBy`与`User`模型关联，JSON序列化时可选输出。
- **CreatedAt**: 创建时间，自动记录部署任务创建时间。
- **UpdatedAt**: 更新时间，自动记录最后一次更新时间。
- **DeletedAt**: 软删除时间戳，由GORM管理，带索引支持逻辑删除。
- **Logs**: 部署日志列表，通过`foreignKey:DeploymentID`与`DeploymentLog`模型关联，形成一对多关系。

**Section sources**
- [deployment.go](file://backend/internal/model/deployment.go#L1-L36)

## GORM关联关系实现机制

### Belongs To 关系实现
`Deployment` 模型通过 `Belongs To` 关系与 `User` 和 `Server` 模型关联，具体实现如下：

- **与User的关联**：`Deployment` 结构体中的 `User` 字段通过 `gorm:"foreignKey:CreatedBy"` 指定外键为 `CreatedBy`。这意味着每个部署任务都属于一个创建它的用户。当查询部署任务时，可通过 `Preload("User")` 加载关联的用户信息。
- **与Server的关联**：`Deployment` 结构体中的 `Server` 字段通过 `gorm:"foreignKey:ServerID"` 指定外键为 `ServerID`。这表示每个部署任务都针对一个特定的服务器。同样，可通过 `Preload("Server")` 加载服务器详情。

### 外键约束
GORM在迁移时会自动创建外键约束，确保数据完整性：
- `ServerID` 字段的外键约束确保部署任务只能关联到存在的服务器。
- `CreatedBy` 字段的外键约束确保部署任务只能由存在的用户触发。
- 这些约束在数据库层面防止了孤立记录的产生，提升了数据一致性。

**Section sources**
- [deployment.go](file://backend/internal/model/deployment.go#L1-L36)
- [user.go](file://backend/internal/model/user.go#L1-L28)
- [server.go](file://backend/internal/model/server.go#L1-L32)

## 数据库表设计与状态流转

### deployments表结构
根据 `init.sql` 脚本和GORM模型定义，`deployments` 表包含以下字段：
- `id`: 主键
- `name`: 部署名称
- `server_id`: 服务器ID（索引）
- `repository`: 仓库地址
- `branch`: 分支名
- `path`: 部署路径
- `script`: 部署脚本
- `status`: 状态码
- `created_by`: 用户ID（索引）
- `created_at`, `updated_at`, `deleted_at`: 时间戳

### 状态流转设计
状态设计采用整数枚举，共四种状态：
- **pending (0)**: 初始状态，任务已创建但未开始执行。
- **running (1)**: 任务正在执行中，脚本正在服务器上运行。
- **success (2)**: 任务成功完成，所有步骤执行无误。
- **failed (3)**: 任务执行失败，可能由于脚本错误、网络问题等。

**设计考量**：
- 使用整数而非字符串提高存储和查询效率。
- 状态值连续且有序，便于状态转换验证。
- 默认状态为`pending`，确保新任务处于待处理状态。

### 索引优化
- `server_id` 和 `created_by` 字段均建立了索引，支持按服务器或用户快速查询部署历史。
- `deleted_at` 字段索引支持软删除后的高效查询。
- 这些索引显著提升了在大数据量下的查询性能，特别是在分页查询时。

**Section sources**
- [init.sql](file://scripts/init.sql#L0-L15)
- [deployment.go](file://backend/internal/model/deployment.go#L1-L36)

## GORM操作代码示例

### 创建部署任务
```go
// 创建部署任务
deployment := &model.Deployment{
    Name:       "前端部署",
    ServerID:   1,
    Repository: "https://github.com/example/frontend.git",
    Branch:     "main",
    Path:       "/var/www/html",
    Script:     "npm install && npm run build && cp -r dist/* /var/www/html/",
    CreatedBy:  1, // 用户ID
}

// 保存到数据库
if err := db.Create(deployment).Error; err != nil {
    log.Printf("创建部署失败: %v", err)
    return
}
```

### 查询部署历史记录
```go
// 查询特定用户的部署历史（分页）
var deployments []model.Deployment
var total int64

// 获取总数
db.Model(&model.Deployment{}).Where("created_by = ?", userID).Count(&total)

// 分页查询，预加载用户和服务器信息
err := db.Preload("User").Preload("Server").
    Where("created_by = ?", userID).
    Offset((page-1)*pageSize).Limit(pageSize).
    Order("created_at DESC").
    Find(&deployments).Error
```

**Section sources**
- [deployment.go](file://backend/internal/model/deployment.go#L1-L36)
- [types.go](file://backend/internal/api/types.go#L71-L114)

## 大数据量下的分页与归档策略

### 分页策略
- 使用 `PageRequest` 结构体统一处理分页参数，包含 `Page` 和 `PageSize` 字段。
- 查询时通过 `Offset` 和 `Limit` 实现物理分页，避免内存溢出。
- 结合 `Order("created_at DESC")` 确保按时间倒序返回最新记录。
- 返回 `PageResponse` 包含列表、总数、当前页和页大小，便于前端分页控件使用。

### 归档策略
- **时间归档**：定期将超过一定时间（如6个月）的部署记录移动到历史表 `deployments_archive`。
- **状态归档**：将状态为 `success` 且超过3个月的记录归档，保留 `failed` 记录更长时间以便故障分析。
- **索引维护**：归档后重建主表索引，保持查询性能。
- **软删除标记**：对于需要保留但不常查询的记录，可使用 `deleted_at` 标记而非物理删除。

这些策略确保系统在长期运行后仍能保持良好的查询性能和存储效率。

**Section sources**
- [types.go](file://backend/internal/api/types.go#L109-L114)
- [deployment.go](file://backend/internal/model/deployment.go#L1-L36)

## 结论
`Deployment` 模型是自动化运维平台的核心组件之一，通过精心设计的字段、关联关系和状态机，实现了完整的部署任务管理。结合GORM的ORM能力，提供了高效的数据操作接口。合理的索引和分页策略确保了系统在大数据量下的稳定性和性能。该模型为后续的部署自动化、监控和审计功能奠定了坚实的基础。