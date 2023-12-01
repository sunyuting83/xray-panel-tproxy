import { useState, useEffect } from 'react'
import {urilist, httpServer} from '../utils/index'
import Loading from '../public/Loading'
import Notification from '../public/Notification'

const White = () => {
  const [proxy, setProxy] = useState("")
  const [direct, setDirect] = useState("")
  const [loading, setLoading] = useState(true)
  const [dis, setDis] = useState(false)
  const [notification, setNotification] = useState({})
  useEffect(() => {
    async function getData() {
      const d = await httpServer(urilist.getdomains)
      setProxy(d.proxy)
      setDirect(d.direct)
      setLoading(false)
    }
    getData()
  }, [])
  const bindProxy = (e) => {
    setProxy(e.target.value)
  }
  const bindDirect = (e) => {
    setDirect(e.target.value)
  }
  const Modify = async() => {
    setDis(true)
    const params = {
      proxy: proxy,
      direct: direct,
    }
    const d = await httpServer(urilist.setdomains, params, 'put')
    if (d.status === 0) {
      OpenNotification(d.message, 'success')
      setLoading(true)
      const data = await httpServer(urilist.getdomains)
      setProxy(data.proxy)
      setDirect(data.direct)
      setLoading(false)
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
            <h5 className="title is-5">走代理的域名</h5>
            <h6 className="subtitle is-6">设置必须走代理的域名，每行一个</h6>
            <div className="control">
              <textarea name="Proxy" className="textarea" rows={8} value={proxy} onChange={(e)=>bindProxy(e)} />
            </div>
          </div>
          <div className="field mt-5">
            <h5 className="title is-5">不走代理的域名</h5>
            <h6 className="subtitle is-6">设置绕过代理的域名，每行一个</h6>
            <div className="control">
              <textarea name="Direct" className="textarea" rows={12} value={direct} onChange={(e) => {bindDirect(e)}} />
            </div>
          </div>
          <div className="field is-grouped">
            <div className="control">
              <button className="button is-link" onClick={Modify} disabled={dis? true : false}>修改</button>
            </div>
          </div>
        </div>
      }
      {notification.active ?<Notification props={notification} close={CloseNotification} /> : <></>}
    </>
  )
}
export default White