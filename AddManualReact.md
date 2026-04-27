# 🚀 Agent Instruction: Manual Node Form Implementation

## 🎯 任务目标
在 `xpanel-web/src/pages/` 下创建 `AddManualNode.jsx`。实现一个多协议切换的手动添加节点页面，要求 UI 响应迅速、字段对齐 Go 后端 `CodeList` 结构、并使用既有的 `httpServer` 进行通信。

---

## 后端表单需要对齐的结构体
```
type CodeList struct {
	// --- 基础识别 ---
	Index   int    `json:"index"` // 节点索引
	UID     string `json:"uid"`
	Type    string `json:"type"`    // vless, vmess, shadowsocks, trojan
	Title   string `json:"title"`   // 节点名称
	Address string `json:"address"` // 服务器 IP/域名
	Port    int    `json:"port"`    // 端口

	// --- 认证信息 ---
	ID       string `json:"id"`       // vless/vmess 的 UUID
	Password string `json:"password"` // ss/trojan 的密码
	Security string `json:"security"` // vless 专用, 一般为 "none"

	// --- 传输层 (StreamSettings) ---
	Network string `json:"network"` // tcp, ws, grpc, h2, mkcp, quic
	Path    string `json:"path"`    // ws/grpc/h2 的路径
	Host    string `json:"host"`    // ws/h2 的 host 头部

	// --- 安全层 (TLS/Reality) ---
	StreamSecurity string `json:"stream_security"` // "", "tls", "reality"
	Sni            string `json:"sni"`             // 服务器域名 (SNI)
	Fingerprint    string `json:"fingerprint"`     // 指纹: chrome, edge, safari, firefox
	AllowInsecure  bool   `json:"allow_insecure"`  // 跳过证书检查

	// Reality 特有
	PublicKey string `json:"public_key"`
	ShortId   string `json:"short_id"`
	SpiderX   string `json:"spider_x"`

	// --- 协议特性 ---
	Flow    string `json:"flow"`   // vless 专用: xtls-rprx-vision
	Method  string `json:"method"` // ss 专用: 加密方法
	AlterID int    `json:"aid"`    // vmess 专用 (虽然现在通常为 0)
}
```

---

## 🛠️ 技术约束
1.  **文件路径**：仅限 `xpanel-web` 目录，禁止越权访问。
2.  **样式框架**：使用 **Bulma CSS**。
3.  **接口库**：使用 `xpanel-web/utils/http.js` 中的 `httpServer`。
4.  **配置库**：使用 `xpanel-web/utils/index.js` 中的 `urilist`。
5.  **协议类型**：支持 `ss`, `trojan`, `vless`, `vmess`, `socks5`。

---

## 📄 页面逻辑与 UI 规范

### 1. 状态管理 (useState)
定义一个 `type` 状态（默认 `ss`）和 `formData` 对象。
```javascript
const [activeType, setActiveType] = useState('ss');
const [formData, setFormData] = useState({
  type: 'ss',
  title: '',
  address: '',
  port: 443,
  // ... 各协议共有及特有字段初始值
});
```

### 2. 协议切换逻辑
在页面顶部渲染一组按钮。
- 当前选中的协议按钮需增加 `disabled` 属性并应用 `is-link` 类名。
- 点击时清空/重置 `formData` 并更新 `activeType`。

### 3. 表单字段映射 (Form Fields)

| 协议 | 必填字段 (基于 CodeList 结构) |
| :--- | :--- |
| **Common** | `title`, `address`, `port` |
| **ss** | `method`, `password` (隐去所有 streamSettings) |
| **trojan** | `password`, `sni` |
| **vless** | `id`, `flow`, `network`, `stream_security`, `sni`, `public_key`, `short_id` |
| **vmess** | `id`, `aid`, `network`, `path`, `host` |
| **socks5** | `id` (as username), `password` |

---

## 📜 代码生成指南 (用于 AddManualNode.jsx)

### A. 提交函数实现
```javascript
import { urilist, httpServer } from '../utils/index';

const handleAddNode = async () => {
  try {
    // 处理 shadowsocks 的特殊逻辑：提交前确保 streamSettings 相关字段为空字符串
    let payload = { ...formData, type: activeType };
    if (activeType === 'ss' || activeType === 'socks5') {
      payload.network = "";
      payload.stream_security = "";
    }

    const response = await httpServer(urilist.addManualNode, payload, "POST");
    if (response.status === 0) {
      // 提示成功逻辑
    }
  } catch (err) {
    console.error("提交失败", err);
  }
};
```

### B. UI 结构参考 (Bulma 风格)
```jsx
<div className="container p-5">
  {/* 协议切换按钮组 */}
  <div className="buttons has-addons is-centered">
    {['ss', 'trojan', 'vless', 'vmess', 'socks5'].map((t) => (
      <button 
        key={t}
        className={`button ${activeType === t ? 'is-link' : ''}`}
        disabled={activeType === t}
        onClick={() => setActiveType(t)}
      >
        {t.toUpperCase()}
      </button>
    ))}
  </div>

  <div className="box">
    {/* 通用字段 */}
    <div className="field">
      <label className="label">节点名称</label>
      <div className="control">
        <input className="input" type="text" name="title" value={formData.title} onChange={handleInputChange} />
      </div>
    </div>
    
    {/* 动态字段渲染逻辑... */}
    {activeType === 'vless' && (
      <div className="field">
        <label className="label">UUID (ID)</label>
        <div className="control">
          <input className="input" type="text" name="id" value={formData.id} onChange={handleInputChange} />
        </div>
      </div>
    )}

    {/* 提交按钮 */}
    <div className="control">
      <button className="button is-primary is-fullwidth" onClick={handleAddNode}>保存节点</button>
    </div>
  </div>
</div>
```

---

## 🚩 特别提醒 Agent
- **onChange 处理**：编写一个通用的 `handleInputChange` 函数来处理 `e.target.name` 和 `e.target.value`。
- **字段对齐**：JSON Payload 必须包含 `stream_security` (下划线) 而非 `streamSecurity`。
- **数据类型**：`port` 和 `aid` 必须确保是 `Number` 类型后再发送。