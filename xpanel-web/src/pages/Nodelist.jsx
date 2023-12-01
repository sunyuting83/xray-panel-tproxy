import { useState, useEffect } from 'react'
import {urilist, httpServer} from '../utils/index'
import Loading from '../public/Loading'
import EmptyEd from '../public/Empty'
import Notification from '../public/Notification'
import { useWsContext } from '../public/ws';

const NodeList = () => {
  const wsuuid = localStorage.getItem("uuid")
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [notification, setNotification] = useState({})
  const [modal, setModal] = useState({})
  const [current, setCurrent] = useState("")
  const [uuid, setUUID] = useState("")
  const globalWs = useWsContext()
  const [ggstatus, setGgstatus] = useState(true)
  const [ggtime, setGgtime] = useState(0)
  const [bdstatus, setBdstatus] = useState(true)
  const [bdtime, setBdtime] = useState(0)
  const [ytstatus, setYtstatus] = useState(true)
  const [yttime, setYttime] = useState(0)
  const [ghstatus, setGhstatus] = useState(true)
  const [ghtime, setGhtime] = useState(0)
  const [tcping, setTcping] = useState(false)
  useEffect(() => {
    const wsuuid = localStorage.getItem("uuid")
    setUUID(wsuuid)
    async function getData() {
      const d = await httpServer(urilist.nodelist)
      setData(d.date)
      const node = localStorage.getItem('current')
      setCurrent(node)
      setLoading(false)
    }
    getData()
    globalWs.onmessage = (data) => {
      if (typeof JSON.parse(data.data) === "object") {
        const jsonData = JSON.parse(data.data)
        if (jsonData.type === "testspeed") {
          if (jsonData.data.indexOf('||||') !== -1) {
            const splitData = jsonData.data.split("||||")
            const speedNumber = parseFloat(splitData[1])
            switch (splitData[0]) {
              case 'https://www.google.com.hk':
                if (speedNumber > 0) {
                  setGgtime(speedNumber)
                }
                setGgstatus(false)
                break;
              case 'https://baidu.com':
                if (speedNumber > 0) {
                  setBdtime(speedNumber)
                }
                setBdstatus(false)
                break;
              case 'https://www.youtube.com/img/desktop/yt_1200.png':
                if (speedNumber > 0) {
                  setYttime(speedNumber)
                }
                setYtstatus(false)
                break;
              case 'https://github.com/webgl-globe/data/data.json':
                if (speedNumber > 0) {
                  setGhtime(speedNumber)
                }
                setGhstatus(false)
                break;
              
              default:
                break;
            }
          }
        }
      }
    }
  }, [globalWs])
  const OpenNotification = (message, color) => {
    setNotification({active: true, message: message, color: color})
  }
  const OpenModal = () => {
    setModal({active: true})
    setTimeout(()=>{
      const sendData = JSON.stringify({'type': 'testspeed', 'uuid': wsuuid, 'data': 'active'})
      if (globalWs.readyState === 1) globalWs.send(sendData)
    }, 1500)
  }
  const CloseModal = () => {
    setGgtime(0)
    setGgstatus(true)
    setBdtime(0)
    setBdstatus(true)
    setYttime(0)
    setYtstatus(true)
    setGhtime(0)
    setGhstatus(true)
    setModal({active: false})
  }
  const CloseNotification = () => {
    setNotification({active: false, message: '', color: ''})
  }

  const TCPing = () => {
    const sendData = JSON.stringify({'type': 'tcping', 'uuid': uuid, 'data': 'active'})
    if (globalWs.readyState === 1) globalWs.send(sendData)
    setTcping(true)

    globalWs.onmessage = (d) => {
      if (typeof JSON.parse(d.data) === "object") {
        const jsonData = JSON.parse(d.data)
        if (jsonData.type === "tcping") {
          if (jsonData.data.indexOf('||||') !== -1) {
            const splitData = jsonData.data.split("||||")
            const i = parseInt(splitData[0])
            const speedNumber = parseFloat(splitData[1])
            // console.log(data, i, speedNumber)
            setData(data.map((e, index) =>{
              if (i === index) {
                e.ping = speedNumber
              }
              return e
            }))
            // console.log(i, data.length)
            if (i === data.length - 1) setTcping(false)
          }
        }
      }
    }
  }

  const DeleteNode = async(i) =>{
    const param = {
      node: String(i)
    }
    const d = await httpServer(urilist.deletenode, param, 'DELETE')
    if (d.status === 0) {
      const newData = data.filter((e, index) => {
        return index !== i
      })
      setData(newData)
    }
  }
  const SetNode = async(i, title) => {
    const param = {
      node: String(i)
    }
    const d = await httpServer(urilist.setnode, param, 'put')
    if (d.status === 0) {
      localStorage.setItem('current', title)
      setCurrent(title)
      OpenModal()
      // OpenNotification(d.message, 'success')
    }else {
      OpenNotification(d.message, 'danger')
    }
  }
  const ShowSpeed = (speed) => {
    return (
      <>
        {speed !== "" ?
          <span>
            {speed > 0 ? speed + 'ms' : '失败'}
          </span>
          :
          <span>待测试</span>
        }
      </>
    )
  }

  const StatusTime = (timestamp) => {
    if (timestamp > 0) {
      return `用时: ${timestamp}秒`
    }else{
      return "测试失败"
    }
  }
  const OpenModle = () => {
    return (
      <>{modal.active ?
        <div id="modal" className="modal is-active">
          <div className="modal-background"></div>
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title is-size-5">设置成功 正在测试连通性</p>
              <button className="delete" aria-label="close" onClick={CloseModal}></button>
            </header>
            <section className="modal-card-body">
              <div className="columns is-mobile is-flex is-flex-wrap-wrap">
                <div className="column is-half">
                  <p className="is-size-7">百度: {bdstatus ? "loading..." : StatusTime(bdtime)}</p>
                </div>
                <div className="column is-half">
                  <p className="is-size-7">Google:{ggstatus ? "loading..." : StatusTime(ggtime)}</p>
                </div>
                <div className="column is-half">
                  <p className="is-size-7">Github:{ghstatus ? "loading..." : StatusTime(ghtime)}</p>
                </div>
                <div className="column is-half">
                  <p className="is-size-7">Youtube:{ytstatus ? "loading..." : StatusTime(yttime)}</p>
                </div>
              </div>
            </section>
            <footer className="modal-card-foot">
              <button className="button is-success" onClick={CloseModal}>关闭</button>
            </footer>
          </div>
        </div>
        :<></>}
      </>
    )
  }
  return (
    <>
      {
        loading ? <Loading /> : 
        <>
          {data.length === 0 ?
          <div>
            <EmptyEd></EmptyEd>
          </div>
          :
          <>
          <div className="columns flex-wrap is-justify-content-space-between mt-1">
            <div className="field ml-3">
              <div className="buttons is-horizontal are-small has-addons">
                <button className="button is-warning" onClick={TCPing} disabled={tcping}>
                  Ping测速
                </button>
              </div>
            </div>
          </div>
          <h6 className="subtitle is-6">Ping测速只能说明节点的端口是通的,无法确认连通性。测试节点是否可用,请使用连通性测试。</h6>
          <div className="table-container">
            <table className="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
              <thead className="is-size-7">
                <tr>
                  <td>序号</td>
                  <td>类型</td>
                  <td>节点名称</td>
                  <td>测速</td>
                  <td>操作</td>
                </tr>
              </thead>
              <tbody className="is-size-7">
                {data.length > 0 &&
                  data.map((item, index) => {
                    return(
                      <tr key={index}>
                        <td>{index + 1}</td>
                        <td>{item.types}</td>
                        <td>{item.title}</td>
                        <td>{ShowSpeed(item.ping)}</td>
                        <td>
                          <div className="buttons">
                            <button className="button is-success is-small" disabled={current === item.title ? true : false } onClick={() =>SetNode(index, item.title)}>{current === item.title ? '当前节点' : '使用节点'}</button>
                            <button className="button is-info is-small" onClick={() =>DeleteNode(index)}>删除节点</button>
                          </div>
                        </td>
                      </tr>
                    )
                })}
              </tbody>
            </table>
          </div>
          </>
          }
        </>
      }
      {notification.active ?<Notification props={notification} close={CloseNotification} /> : <></>}
      {<OpenModle />}
    </>
  )
}
export default NodeList