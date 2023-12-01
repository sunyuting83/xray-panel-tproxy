import { useState, useEffect } from 'react'
import {urilist, httpServer} from '../utils/index'
import Loading from '../public/Loading'
import Notification from '../public/Notification'

const Subscribe = () => {
  const [sub, setSub] = useState([])
  const [loading, setLoading] = useState(true)
  const [proxy, setProxy] = useState(false)
  const [notification, setNotification] = useState({})
  const [dis, setDis] = useState(false)
  const [ignore, setIgnore] = useState("")
  
  useEffect(() => {
    async function getData() {
      const d = await httpServer(urilist.getsubscribes)
      if (d.status === 0) {
        let subArray = []
        if (d.subscribes.indexOf('\n') !== -1) {
          subArray = d.subscribes.split('\n')
        }else{
          subArray = [...subArray, d.subscribes]
        }
        setSub(subArray)
        setIgnore(d.ignore)
        setLoading(false)
      }else{
        CloseNotification(d.message, 'danger')
      }
    }
    getData()
  }, [])
  const bindSub = (e, i) => {
    setSub(sub.map((el, index) =>{
      if(index === i) { 
        el = e.target.value
      }
      return el
    }))
  }
  const bindIgnore = (e) => {
    setIgnore(e.target.value)
  }
  const AddSub = () => {
    setSub([...sub, ""])
  }
  const DelSub = (i) => {
    setSub(sub.filter((_, index) =>  index !== i))
  }
  const Modify = async() => {
    setDis(true)
    let subData = ""
    if (sub.length > 1) {
      subData = sub.join("\n")
    }else{
      subData = sub[0]
    }
    const params = {
      subscribes: subData
    }
    const d = await httpServer(urilist.setSubscribes, params, 'put')
    if (d.status === 0) {
      setSub([])
      OpenNotification(d.message, 'success')
      setLoading(true)
      const data = await httpServer(urilist.getsubscribes)
      let subArray = []
      if (data.subscribes.indexOf('\n') !== -1) {
        subArray = data.subscribes.split('\n')
      }else{
        subArray = [...subArray, data.subscribes]
      }
      setSub(subArray)
      setLoading(false)
    }else{
      OpenNotification(d.message, 'danger')
    }
    setDis(false)
  }
  const Updata = async() => {
    setLoading(true)
    setDis(true)
    let params = {
      proxy: "0"
    }
    if (proxy) {
      params.proxy = "1"
    }
    const data = await httpServer(urilist.updata, params, 'get')
    if (data.status === 0) {
      OpenNotification(data.message, 'success')
    }else {
      OpenNotification(data.message, 'danger')
    }
    setLoading(false)
    setDis(false)
  }
  const saveIgnore = async() => {
    setDis(true)
    setLoading(true)
    const params = {
      ignore: ignore
    }
    const d = await httpServer(urilist.setIgnore, params, 'put')
    if (d.status === 0) {
      OpenNotification(d.message, 'success')
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
          <h5 className="title is-5">订阅地址</h5>
          <h6 className="subtitle is-6">一行一条订阅地址，如果有修改，请先保存后再更新。</h6>
          {sub.map((item, i) => {
          return(
          <div className="field has-addons" key={i}>
            <div className="control" style={{'width': '100%'}}>
              <input className="input" type="text" value={item} onChange={(e) => {bindSub(e, i)}} />
            </div>
            <div className="control">
              {sub.length > 1 ?
              <button className="button is-info" onClick={() => {DelSub(i)}}>
                删除本行
              </button>:<></>}
            </div>
          </div>)
          })}
          <div className="field is-grouped">
            <div className="control">
              <div className="buttons are-small">
                <button className="button is-small is-primary is-light" onClick={AddSub}>增加一行</button>
                <button className="button is-small is-danger is-light" onClick={Updata}>立即更新</button>
                <label className="checkbox">
                  <input type="checkbox" value={proxy} onChange={() => {setProxy(!proxy)}} />
                  使用代理更新
                </label>
              </div>
            </div>
          </div>
          <div className="field is-grouped">
            <div className="control">
                <button className="button is-link" onClick={Modify}  disabled={dis? true : false}>保存</button>
            </div>
          </div>
          <h5 className="title is-5 mt-6">过滤关键字</h5>
          <h6 className="subtitle is-6">包含以下文字的节点将被过滤，使用|分割多个词,输入框失去焦点自动保存</h6>
          <div className="field has-addons">
            <div className="control" style={{'width': '100%'}}>
              <input className="input" type="text" value={ignore} onChange={(e) => {bindIgnore(e)}} onBlur={saveIgnore} />
            </div>
          </div>
        </div>
      }
      {notification.active ?<Notification props={notification} close={CloseNotification} /> : <></>}
    </>
  )
}
export default Subscribe