import { useEffect } from 'react'
import {NavLink} from 'react-router-dom'
import { useWsContext } from '../public/ws'
const Header = () => {
  const globalWs = useWsContext()
  useEffect(() => {
    let uuid = ""
    globalWs.onmessage = (data) => {
      if (typeof JSON.parse(data.data) === "object") {
        const jsonData = JSON.parse(data.data)
        uuid = jsonData.uuid
        localStorage.setItem("uuid", jsonData.uuid)
      }
    }
    setInterval(() => {
      if (uuid !== "") {
        const sendData = JSON.stringify({'type': 'active', 'uuid': uuid, 'data': 'active'})
        if (globalWs.readyState === 1) globalWs.send(sendData)
      }
    }, 1000 * 60 * 5)
  },[globalWs])
  return (
    <div className="navbar">
      <div className="navbar-tabs">
        <NavLink to="/" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>状态</NavLink>
        <NavLink to="/nodelist" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>节点管理</NavLink>
        <NavLink to="/subscribe" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>订阅设置</NavLink>
        <NavLink to="/white" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>域名分流设置</NavLink>
        <NavLink to="/setdns" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>DNS设置</NavLink>
        <NavLink to="/update" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>更新</NavLink>
        {/* <NavLink to="/setsocks" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>设置本地Socks</NavLink> */}
      </div>
    </div>
  )
}
export default Header