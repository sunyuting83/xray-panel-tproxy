import {NavLink} from 'react-router-dom'
const Header = () => {
  return (
    <div className="navbar">
      <div className="navbar-tabs">
        <NavLink to="/" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>状态</NavLink>
        <NavLink to="/nodelist" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>节点管理</NavLink>
        <NavLink to="/subscribe" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>订阅设置</NavLink>
        <NavLink to="/white" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>分流设置</NavLink>
        <NavLink to="/setdns" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>DNS设置</NavLink>
        <NavLink to="/update" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>更新</NavLink>
        {/* <NavLink to="/setsocks" className={({ isActive }) => "navbar-item is-tab" + (isActive ? " is-active" : "")}>设置本地Socks</NavLink> */}
      </div>
    </div>
  )
}
export default Header