# Router 拆分 Agent

## 任务目标
将 `Router/router.go` 中第 457 行到文件末尾的路由处理函数，拆分成独立的 `.go` 文件，每个文件包含一个函数，放置在 `Router` 目录下。

## 工作目录限制
**严禁访问 `Router` 以外的任何目录。** 所有操作必须在 `Router/` 目录内进行。

---

## 执行步骤

### Step 1: 备份原文件
```
操作: 复制 Router/router.go → Router/router.go.bak
```
**目的**: 保留回滚依据，严禁删除此备份。

---

### Step 2: 读取并分析目标区域
```
操作: 读取 Router/router.go 第 457 行到 EOF
```

**分析内容:**
1. 找出所有 `func ` 开头的函数定义
2. 记录每个函数的：
   - 函数名（如 `UpData`）
   - 函数签名（如 `func UpData(c *gin.Context)`）
   - 函数体起止行号
   - 该函数使用的 import 包

---

### Step 3: 提取函数列表
**判断标准:**
- 函数签名格式为 `func XXX(c *gin.Context)` 或 `func XXX(c *gin.Context) { ... }`
- 这些函数在路由注册中被引用

**示例提取结果:**
```
函数: UpData
  签名: func UpData(c *gin.Context)
  位置: router.go 第 460-520 行
  使用 import: gin, json, fmt
```

---

### Step 4: 创建独立文件

**命名规则:**
```
函数名 → 文件名（snake_case）
UpData      → up_data.go
GetUserList → get_user_list.go
Index       → index.go
```

**文件位置:** `Router/` 目录下

**文件模板:**
```go
package router

import (
    "github.com/gin-gonic/gin"
    // 仅该函数实际使用的其他 import
)

// UpData 处理 /updata 路由
// [原函数的完整代码，严禁修改业务逻辑]
func UpData(c *gin.Context) {
    // ... 原代码完整复制，一字不改 ...
}
```

**Import 处理原则:**
1. 查看原 `router.go` 的 import 块
2. 分析每个函数实际使用了哪些包
3. 只在新文件中保留该函数实际需要的 import
4. 如果函数没有使用某个 import，不要在新文件中引入
5. 每个函数文件必须包含 `"github.com/gin-gonic/gin"`（因为 handler 签名需要）

---

### Step 5: 修改 router.go

**保留内容（第 431-476 行）:**
```go
package router

import (
    "github.com/gin-gonic/gin"
    // 其他需要的 import
)

// 路由注册代码 - 必须保留
func InitRouter() *gin.Engine {
    router := gin.New()
    // 初始化 gws Upgrader
    upgrader := gws.NewUpgrader(websocket.Manager, &gws.ServerOption{})
    router.Use(utils.CORSMiddleware())           // ← 保留引用

    api := router.Group("/api")
    {
        api.GET("/updata", UpData)   // ← 保留引用
        // ... 其他路由注册
    }

    return router
}

// 第 457 行之前的其他代码（变量、常量、类型等）全部保留
```

**删除内容（第 457 行到 EOF）:**
- 删除所有函数定义体
- **不要删除**路由注册行中的函数名引用

**修改后 router.go 的末尾:**
第 456 行之后应该是文件结束，或者只保留类型定义/变量定义（如果有的话）。

**注意:** 由于所有文件都是 `package Router`，同包函数之间可以直接调用，无需 import。

---

### Step 6: 清理 import
检查 `router.go` 的 import 块：
- 如果某个 import 只在被拆分的函数中使用，且 `router.go` 中不再使用 → 从 `router.go` 的 import 中删除
- 如果 import 仍然被 `router.go` 中的代码使用 → 保留

---

### Step 7: 验证

**检查清单:**
- [ ] 所有新文件都在 `Router/` 目录下
- [ ] 没有访问过 `Router/` 以外的目录
- [ ] 每个函数文件都使用 `package Router`
- [ ] 函数体代码与原备份完全一致
- [ ] `router.go` 中的路由注册行全部保留
- [ ] `router.go.bak` 仍然存在

---

## 关键约束

| 约束 | 说明 |
|------|------|
| **目录限制** | 只能读写 `Router/` 目录下的文件 |
| **不修改逻辑** | 函数体内的代码必须原样复制，一字不改 |
| **包名一致** | 所有新文件使用 `package router`，与原文件同包 |
| **Import 精准** | 每个新文件只导入该函数实际使用的包 |
| **路由保留** | `router.go` 中的路由注册行必须保留 |
| **备份保留** | `router.go.bak` 必须保留 |

---

## 特殊情况处理

### 函数间相互调用
如果函数 A 调用了同文件中的函数 B：
- 两个函数分别拆分到 `a.go` 和 `b.go`
- 同包（`package Router`），无需额外 import，直接调用即可

### 共享变量/常量
如果函数使用了第 457 行之前定义的变量或常量：
- 这些变量/常量保留在 `router.go` 中
- 新文件中的函数仍然可以访问（同包）

### 类型定义在第 457 行之后
如果第 457 行之后有 `type XXX struct` 定义：
- 优先保留在 `router.go` 中
- 如果类型只被某个拆分后的函数使用，可移到该函数文件中

---

## 输出报告

完成后输出以下报告：

```
=== Router 拆分报告 ===

备份文件:
  Router/router.go.bak ✓

拆分的函数:
  ┌─────────────┬──────────────────────┐
  │ 函数名      │ 文件                 │
  ├─────────────┼──────────────────────┤
  │ UpData      │ Router/up_data.go    │
  │ Index       │ Router/index.go      │
  │ GetUserList │ Router/get_user_list │
  │ ...         │ ...                  │
  └─────────────┴──────────────────────┘

原 router.go 修改:
  - 删除行范围: 457 → EOF
  - 保留: package, import, 路由注册, 变量/常量/类型定义
  - 清理未使用的 import: [是/否]

目录访问确认:
  - 仅访问 Router/ 目录 ✓
  - 未访问其他目录 ✓

编译状态: [通过/失败]
```