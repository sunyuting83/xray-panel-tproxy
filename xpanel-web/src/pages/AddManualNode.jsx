import { useState } from 'react'
import { urilist, httpServer } from '../utils/index'
import Notification from '../public/Notification'

const AddManualNode = () => {
  const [activeType, setActiveType] = useState('ss')
  const [dis, setDis] = useState(false)
  const [notification, setNotification] = useState({})

  const initialForm = {
    type: 'ss',
    title: '',
    address: '',
    port: 443,
    id: '',
    password: '',
    security: 'none',
    network: 'tcp',
    path: '',
    host: '',
    stream_security: '',
    sni: '',
    fingerprint: 'chrome',
    allow_insecure: false,
    public_key: '',
    short_id: '',
    spider_x: '/',
    flow: '',
    method: 'aes-256-gcm',
    aid: 0,
  }

  const [formData, setFormData] = useState({ ...initialForm })

  // --- 选项常量 ---
  const protocolList = ['ss', 'trojan', 'vless', 'vmess', 'socks5']
  const ssMethods = ['aes-256-gcm', 'aes-128-gcm', 'chacha20-poly1305', 'none']
  const networkTypes = ['tcp', 'ws', 'grpc', 'h2']
  const streamSecurities = ['', 'tls', 'reality']
  const fingerprints = ['chrome', 'edge', 'safari', 'firefox']
  const flows = ['', 'xtls-rprx-vision']

  const handleTypeChange = (t) => {
    setActiveType(t)
    // 切换协议时，注入该协议必要的默认值
    const newForm = { ...initialForm, type: t }
    if (t === 'vless') newForm.security = 'none'
    if (t === 'trojan') newForm.stream_security = 'tls'
    setFormData(newForm)
  }

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target
    let val = type === 'checkbox' ? checked : value
    if (name === 'port' || name === 'aid') val = Number(value)
    setFormData((prev) => ({ ...prev, [name]: val }))
  }

  const handleAddNode = async () => {
    if (!formData.title.trim() || !formData.address.trim() || !formData.port) {
      OpenNotification('基础信息（名称、地址、端口）不能为空', 'warning')
      return
    }

    setDis(true)
    try {
      // 深度清洗数据：只保留当前协议和传输层需要的字段
      let p = { ...formData, type: activeType }
      
      // 1. 基础清理
      if (activeType === 'ss') {
        p = { type: 'ss', title: p.title, address: p.address, port: p.port, method: p.method, password: p.password }
      } else if (activeType === 'socks5') {
        p = { type: 'socks5', title: p.title, address: p.address, port: p.port, id: p.id, password: p.password }
      } else {
        // VLESS / VMess / Trojan 逻辑清理
        if (activeType !== 'vmess') p.aid = 0
        if (activeType !== 'vless') { p.flow = ''; p.security = '' }
        
        // 传输层清理：如果不选 TLS/Reality，清空相关字段
        if (!p.stream_security) {
          p.sni = ''; p.fingerprint = ''; p.public_key = ''; p.short_id = ''; p.spider_x = ''
        }
        // Network 清理：非 WS/gRPC/H2 清空 path/host
        if (['tcp', 'mkcp'].includes(p.network)) { p.path = ''; p.host = '' }
      }

      const response = await httpServer(urilist.addManualNode, p, 'POST')
      if (response.status === 0) {
        OpenNotification('节点添加成功', 'success')
        setFormData({ ...initialForm, type: activeType })
      } else {
        OpenNotification(response.message || '添加失败', 'danger')
      }
    } catch (err) {
      OpenNotification('请求失败，请检查后端服务', 'danger')
    }
    setDis(false)
  }

  const OpenNotification = (message, color) => {
    setNotification({ active: true, message, color })
  }

  return (
    <div className="container p-5">
      <div className="buttons has-addons is-centered mb-5">
        {protocolList.map((t) => (
          <button
            key={t}
            className={`button is-uppercase ${activeType === t ? 'is-link is-selected' : ''}`}
            disabled={activeType === t || dis}
            onClick={() => handleTypeChange(t)}
          >
            {t}
          </button>
        ))}
      </div>

      <div className="box anim-fade-in">
        <h5 className="title is-5 mb-4">基础配置</h5>
        <div className="columns is-multiline">
          <div className="column is-6">
            <div className="field">
              <label className="label">节点名称</label>
              <input className="input" name="title" value={formData.title} onChange={handleInputChange} placeholder="My Server" />
            </div>
          </div>
          <div className="column is-4">
            <div className="field">
              <label className="label">地址</label>
              <input className="input" name="address" value={formData.address} onChange={handleInputChange} placeholder="example.com" />
            </div>
          </div>
          <div className="column is-2">
            <div className="field">
              <label className="label">端口</label>
              <input className="input" type="number" name="port" value={formData.port} onChange={handleInputChange} />
            </div>
          </div>
        </div>

        <hr />
        <h5 className="title is-5 mb-4">协议与传输配置</h5>

        {/* --- Shadowsocks --- */}
        {activeType === 'ss' && (
          <div className="columns">
            <div className="column is-6">
              <label className="label">加密方法</label>
              <div className="select is-fullwidth">
                <select name="method" value={formData.method} onChange={handleInputChange}>
                  {ssMethods.map(m => <option key={m} value={m}>{m}</option>)}
                </select>
              </div>
            </div>
            <div className="column is-6">
              <label className="label">密码</label>
              <input className="input" name="password" value={formData.password} onChange={handleInputChange} />
            </div>
          </div>
        )}

        {/* --- VLESS / VMess 共同部分 --- */}
        {(activeType === 'vless' || activeType === 'vmess') && (
          <>
            <div className="field">
              <label className="label">用户 ID (UUID)</label>
              <input className="input" name="id" value={formData.id} onChange={handleInputChange} placeholder="UUID" />
            </div>
            <div className="columns">
              <div className="column is-4">
                <label className="label">传输协议 (Network)</label>
                <div className="select is-fullwidth">
                  <select name="network" value={formData.network} onChange={handleInputChange}>
                    {networkTypes.map(n => <option key={n} value={n}>{n}</option>)}
                  </select>
                </div>
              </div>
              <div className="column is-4">
                <label className="label">安全传输 (StreamSecurity)</label>
                <div className="select is-fullwidth">
                  <select name="stream_security" value={formData.stream_security} onChange={handleInputChange}>
                    {streamSecurities.map(s => <option key={s} value={s}>{s || 'none'}</option>)}
                  </select>
                </div>
              </div>
              {activeType === 'vless' && (
                <div className="column is-4">
                  <label className="label">流控 (Flow)</label>
                  <div className="select is-fullwidth">
                    <select name="flow" value={formData.flow} onChange={handleInputChange}>
                      {flows.map(f => <option key={f} value={f}>{f || 'none'}</option>)}
                    </select>
                  </div>
                </div>
              )}
              {activeType === 'vmess' && (
                <div className="column is-4">
                  <label className="label">AlterID</label>
                  <input className="input" type="number" name="aid" value={formData.aid} onChange={handleInputChange} />
                </div>
              )}
            </div>

            {/* 高级传输层：WS/gRPC 路径 */}
            {['ws', 'grpc', 'h2'].includes(formData.network) && (
              <div className="columns">
                <div className="column is-6">
                  <label className="label">Path</label>
                  <input className="input" name="path" value={formData.path} onChange={handleInputChange} placeholder="/v2ray" />
                </div>
                <div className="column is-6">
                  <label className="label">Host</label>
                  <input className="input" name="host" value={formData.host} onChange={handleInputChange} placeholder="bing.com" />
                </div>
              </div>
            )}
          </>
        )}

        {/* --- Reality / TLS 共有字段 --- */}
        {['tls', 'reality'].includes(formData.stream_security) && (
          <div className="columns is-multiline">
            <div className="column is-6">
              <label className="label">SNI</label>
              <input className="input" name="sni" value={formData.sni} onChange={handleInputChange} />
            </div>
            <div className="column is-6">
              <label className="label">指纹 (Fingerprint)</label>
              <div className="select is-fullwidth">
                <select name="fingerprint" value={formData.fingerprint} onChange={handleInputChange}>
                  {fingerprints.map(f => <option key={f} value={f}>{f}</option>)}
                </select>
              </div>
            </div>
            {formData.stream_security === 'reality' && (
              <>
                <div className="column is-6">
                  <label className="label">Public Key</label>
                  <input className="input" name="public_key" value={formData.public_key} onChange={handleInputChange} />
                </div>
                <div className="column is-3">
                  <label className="label">Short ID</label>
                  <input className="input" name="short_id" value={formData.short_id} onChange={handleInputChange} />
                </div>
                <div className="column is-3">
                  <label className="label">SpiderX</label>
                  <input className="input" name="spider_x" value={formData.spider_x} onChange={handleInputChange} />
                </div>
              </>
            )}
          </div>
        )}

        {/* --- Trojan / Socks5 --- */}
        {activeType === 'trojan' && (
          <div className="field">
            <label className="label">Password</label>
            <input className="input" name="password" value={formData.password} onChange={handleInputChange} />
          </div>
        )}
        {activeType === 'socks5' && (
          <div className="columns">
            <div className="column is-6"><label className="label">用户名</label><input className="input" name="id" value={formData.id} onChange={handleInputChange} /></div>
            <div className="column is-6"><label className="label">密码</label><input className="input" name="password" value={formData.password} onChange={handleInputChange} /></div>
          </div>
        )}

        <div className="field mt-6">
          <button className={`button is-primary is-fullwidth is-medium ${dis ? 'is-loading' : ''}`} onClick={handleAddNode}>
            🚀 保存并添加到手动节点
          </button>
        </div>
      </div>
      
      {notification.active && (
        <Notification props={notification} close={() => setNotification({ active: false })} />
      )}
    </div>
  )
}

export default AddManualNode