import { useState, useEffect } from 'react'
import {urilist, httpServer} from '../utils/index'
import Loading from '../public/Loading'
import Notification from '../public/Notification'

const SetSocks = () => {
  const [socks, setSocks] = useState([])
  const [auth, setAuth] = useState(false)
  const [loading, setLoading] = useState(true)
  const [dis, setDis] = useState(false)
  const [notification, setNotification] = useState({})
  useEffect(() => {
    async function getData() {
      setLoading(true)
      const d = await httpServer(urilist.getSocks)
      if (d.status === 0) {
        setAuth(d.SocksStatus)
        setSocks(d.Auths)
      }
      setLoading(false)
    }
    getData()
  }, [])
  const setSocksStatus = () => {
    setAuth(!auth)
  }
  const bindSocksUser = (e, index) => {
    setSocks(socks.map((el, i) =>{
      if (i === index) {
        el.user = e.target.value
      }
      return el
    }))
  }
  const bindSocksPas = (e, index) => {
    setSocks(socks.map((el, i) =>{
      if (i === index) {
        el.pass = e.target.value
      }
      return el
    }))
  }
  const AddSocks = () => {
    setSocks([...socks, {"user": "", "pass":""}])
  }
  const DelSocks = (i) => {
    setSocks(socks.filter((_, index) =>  index !== i))
  }

  const SaveSetting = async() => {
    setLoading(true)
    setDis(true)
    const socksJson = JSON.stringify(socks)
    let authstatus = "false"
    if (auth) {
      authstatus = "true"
    }
    const params = {
      sockStuts: authstatus,
      auths: socksJson
    }
    const data = await httpServer(urilist.setLocalSocks, params, 'put')
    if (data.status === 0) {
      const d = await httpServer(urilist.getSocks)
      if (d.status === 0) {
        setAuth(d.SocksStatus)
        setSocks(d.Auths)
        OpenNotification(data.message, 'success')
      }
      setLoading(false)
    }else{
      OpenNotification(data.message, 'danger')
      setLoading(false)
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
            <h5 className="title is-5">本地Socks设置</h5>
            <h6 className="subtitle is-6">如果不了解什么是本地Socks,请不要修改本页</h6>

            <div className="field is-grouped">
              <div className="control">
                <div className="buttons are-small">
                  {auth?
                  <button className="button is-small is-danger is-rounded" onClick={setSocksStatus} disabled={dis? true : false}>关闭本地验证</button>
                  :
                  <button className="button is-small is-primary is-rounded" onClick={setSocksStatus}  disabled={dis? true : false}>开启本地验证</button>
                  }
                </div>
              </div>
            </div>
            {socks.length > 0 && socks.map((item, i) => {
              return(
              <div className="field has-addons" key={i}>
                <p className="control">
                  <span className="button is-static">
                    用户名
                  </span>
                </p>
                <div className="control" style={{'width': '100%'}}>
                  <input className="input" type="text" value={item.user} onChange={(e) => {bindSocksUser(e, i)}} />
                </div>
                <p className="control">
                  <span className="button is-static">
                    密码
                  </span>
                </p>
                <div className="control" style={{'width': '100%'}}>
                  <input className="input" type="text" value={item.pass} onChange={(e) => {bindSocksPas(e, i)}} />
                </div>
                <div className="control">
                  {socks.length > 1 ?
                  <button className="button is-info" onClick={() => {DelSocks(i)}}>
                    删除本行
                  </button>:<></>}
                </div>
              </div>
              )}
            )}
            <div className="field is-grouped">
              <div className="control">
                <div className="buttons are-small">
                  <button className="button is-small is-primary is-light" onClick={AddSocks}>增加一行</button>
                  <button className="button is-small is-danger" disabled={dis? true : false} onClick={SaveSetting}>保存设置</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      }
      {notification.active ?<Notification props={notification} close={CloseNotification} /> : <></>}
    </>
  )
}
export default SetSocks