import { useState, useEffect } from 'react'
import {urilist, httpServer} from '../utils/index'
import Loading from '../public/Loading'
import Notification from '../public/Notification'

const SetDns = () => {
  const [dns, setDns] = useState([])
  const [loading, setLoading] = useState(true)
  const [dis, setDis] = useState(false)
  const [notification, setNotification] = useState({})
  useEffect(() => {
    async function getData() {
      setLoading(true)
      const d = await httpServer(urilist.getdns)
      setDns(d.dns)
      setLoading(false)
    }
    getData()
  }, [])
  const bindDns = (e, i) => {
    setDns(dns.map((el, index) =>{
      if(index === i) { 
        el = e.target.value
      }
      return el
    }))
  }

  const AddDns = () => {
    setDns([...dns, ""])
  }
  const bindDnsAddress = (e, domain, port = false) => {
    setDns(dns.map(el =>{
      if(typeof el !== "string") {
        if (el.domains[0] === domain) {
          if (port) {
            el.port = e.target.value
          }else{
            el.address = e.target.value
          }
        }
      }
      return el
    }))
  }
  const DelDns = (i) => {
    setDns(dns.filter((_, index) =>  index !== i))
  }

  const Modify = async() => {
    setLoading(true)
    setDis(true)
    const dnsJson = JSON.stringify(dns)
    const params = {
      dns: dnsJson
    }
    const data = await httpServer(urilist.setdns, params, 'put')
    if (data.status === 0) {
      const d = await httpServer(urilist.getdns)
      setDns(d.dns)
      OpenNotification(data.message, 'success')
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
            <h5 className="title is-5">DNS设置</h5>
            <h6 className="subtitle is-6">如果不了解什么是DNS,请不要修改本页</h6>
            {dns.map((item, i) => {
              return(
              typeof item === 'string' ?
              <div className="field has-addons" key={i}>
                <div className="control" style={{'width': '100%'}}>
                  <input className="input" type="text" value={item} onChange={(e) => {bindDns(e, i)}} />
                </div>
                <div className="control">
                  {dns.length > 1 ?
                  <button className="button is-info" onClick={() => {DelDns(i)}}>
                    删除本行
                  </button>:<></>}
                </div>
              </div>
              :
              <div className="field has-addons" key={i}>
                <p className="control">
                  <span className="button is-static">
                    {item.domains[0].indexOf("!") !== -1 ? "国外" : "国内"}DNS分流设置
                  </span>
                </p>
                <div className="control" style={{'width': '45%'}}>
                  <input className="input" type="text" value={item.address} onChange={(e) => {bindDnsAddress(e, item.domains[0])}} />
                </div>
                <div className="control" style={{'width': '5%'}}>
                  <input className="input" type="text" value={item.port} onChange={(e) => {bindDnsAddress(e, item.domains[0], true)}}  />
                </div>
                <div className="control" style={{'width': '25%'}}>
                  <input className="input" disabled type="text" defaultValue={item.domains[0]} />
                </div>
                {item.expectIPs && item.expectIPs.length > 0 ?
                <div className="control" style={{'width': '25%'}}>
                  <input className="input" disabled type="text" defaultValue={item.expectIPs[0]} />
                </div>: <div></div>}
              </div>)
            })}
            <div className="field is-grouped">
              <div className="control">
                <div className="buttons are-small">
                  <button className="button is-small is-primary is-light" onClick={AddDns}>增加一行</button>
                  <button className="button is-small is-danger is-light" onClick={Modify}  disabled={dis? true : false}>立即更新</button>
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
export default SetDns