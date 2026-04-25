import { useState, useEffect } from 'react'
import {urilist, httpServer} from '../utils/index'
import Loading from '../public/Loading'
import Notification from '../public/Notification'

const White = () => {
  const [block_domain, setBlock_domain] = useState([])
  const [direct_app, setDirect_app] = useState([])
  const [direct_domain, setDirect_domain] = useState([])
  const [direct_ip, setDirect_ip] = useState([])
  const [proxy_domain, setProxy_domain] = useState([])
  const [loading, setLoading] = useState(true)
  const [dis, setDis] = useState(false)
  const [notification, setNotification] = useState({})
  useEffect(() => {
    async function getData() {
      const d = await httpServer(urilist.getallrules)
      if (d.status === 0) {
        setBlock_domain(d.data.block_domain)
        setDirect_app(d.data.direct_app)
        setDirect_domain(d.data.direct_domain)
        setDirect_ip(d.data.direct_ip)
        setProxy_domain(d.data.proxy_domain)
        setLoading(false)
      } else {
        setBlock_domain("")
        setDirect_app("")
        setDirect_domain("")
        setDirect_ip("")
        setProxy_domain("")
        setLoading(false)
        OpenNotification(d.message, 'danger')
      }
    }
    getData()
  }, [])
  const bindProxy = (e) => {
    switch (e.target.name) {
      case 'Block_domain':
        setBlock_domain(e.target.value)
        break
      case 'Direct_app':
        setDirect_app(e.target.value)
        break
      case 'Direct_domain':
        setDirect_domain(e.target.value)
        break
      case 'Direct_ip':
        setDirect_ip(e.target.value)
        break
      case 'Proxy_domain':
        setProxy_domain(e.target.value)
        break
      default:
        break
    }
  }
  const Modify = async(t) => {
    setDis(true)
    let data = ""
    switch (t) {
      case 'block_domain':
        data = block_domain
        break
      case 'direct_app':
        data = direct_app
        break
      case 'direct_domain':
        data = direct_domain
        break
      case 'direct_ip':
        data = direct_ip
        break
      case 'proxy_domain':
        data = proxy_domain
        break
      default:
        break
    }
    const params = {
      type: t,
      data: data,
    }
    const d = await httpServer(urilist.updateRule, params, 'put')
    if (d.status === 0) {
      OpenNotification(d.message, 'success')
      setLoading(true)
      const redata = await httpServer(urilist.getallrules)
      if (redata.status === 0) {
        setBlock_domain(redata.data.block_domain)
        setDirect_app(redata.data.direct_app)
        setDirect_domain(redata.data.direct_domain)
        setDirect_ip(redata.data.direct_ip)
        setProxy_domain(redata.data.proxy_domain)
        setLoading(false)
      } else {
        setBlock_domain("")
        setDirect_app("")
        setDirect_domain("")
        setDirect_ip("")
        setProxy_domain("")
        setLoading(false)
        OpenNotification(d.message, 'danger')
      }
    }else{
      OpenNotification(d.message, 'danger')
    }
    setDis(false)
  }
  const OpenNotification = (message, color) => {
    setNotification({active: true, message: message, color: color})
  }
  const CloseNotification = () => {
    setNotification({active: false, message: '', color: ''})
  }
  return (
    <>
      {
        loading ? <Loading /> : 
        <div className="box">
          <div className="field">
            <h5 className="title is-5">屏蔽的域名</h5>
            <h6 className="subtitle is-6">设置屏蔽的域名，每行一个</h6>
            <div className="control">
              <textarea name="Block_domain" className="textarea" rows={8} value={block_domain} onChange={(e)=>bindProxy(e)} />
            </div>
          </div>
          <div className="field is-grouped">
            <div className="control">
              <button className="button is-link" onClick={()=>Modify("block_domain")} disabled={dis? true : false}>更新屏蔽的域名</button>
            </div>
          </div>
          <hr />
          <div className="field mt-5">
            <h5 className="title is-5">屏蔽的应用</h5>
            <h6 className="subtitle is-6">设置屏蔽的应用，每行一个</h6>
            <div className="control">
              <textarea name="Direct_app" className="textarea" rows={8} value={direct_app} onChange={(e)=>bindProxy(e)} />
            </div>
          </div>
          <div className="field is-grouped">
            <div className="control">
              <button className="button is-link" onClick={()=>Modify("direct_app")} disabled={dis? true : false}>更新屏蔽的应用</button>
            </div>
          </div>
          <hr />
          <div className="field mt-5">
            <h5 className="title is-5">不走代理的域名</h5>
            <h6 className="subtitle is-6">设置绕过代理的域名，每行一个</h6>
            <div className="control">
              <textarea name="Direct_domain" className="textarea" rows={8} value={direct_domain} onChange={(e)=>bindProxy(e)} />
            </div>
          </div>
          <div className="field is-grouped">
            <div className="control">
              <button className="button is-link" onClick={()=>Modify("direct_domain")} disabled={dis? true : false}>更新不走代理的域名</button>
            </div>
          </div>
          <hr />
          <div className="field mt-5">
            <h5 className="title is-5">不走代理的IP</h5>
            <h6 className="subtitle is-6">设置不走代理的IP，每行一个</h6>
            <div className="control">
              <textarea name="Direct_ip" className="textarea" rows={8} value={direct_ip} onChange={(e)=>bindProxy(e)} />
            </div>
          </div>
          <div className="field is-grouped">
            <div className="control">
              <button className="button is-link" onClick={()=>Modify("direct_ip")} disabled={dis? true : false}>更新不走代理的IP</button>
            </div>
          </div>
          <hr />
          <div className="field mt-5">
            <h5 className="title is-5">走代理的域名</h5>
            <h6 className="subtitle is-6">设置必须走代理的域名，每行一个</h6>
            <div className="control">
              <textarea name="Proxy_domain" className="textarea" rows={12} value={proxy_domain} onChange={(e) => {bindProxy(e)}} />
            </div>
          </div>
          <div className="field is-grouped">
            <div className="control">
              <button className="button is-link" onClick={()=>Modify("proxy_domain")} disabled={dis? true : false}>更新走代理的域名</button>
            </div>
          </div>
        </div>
      }
      {notification.active ?<Notification props={notification} close={CloseNotification} /> : <></>}
    </>
  )
}
export default White