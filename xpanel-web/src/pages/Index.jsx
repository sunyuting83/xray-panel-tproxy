import { useState, useEffect } from 'react'
import {urilist, httpServer} from '../utils/index'
import Loading from '../public/Loading'
import { useWsContext } from '../public/ws'

const Index = () => {
  const [status, setStatus] = useState(false)
  const [loading, setLoading] = useState(true)
  const [ggstatus, setGgstatus] = useState(true)
  const [ggtime, setGgtime] = useState(0)
  const [bdstatus, setBdstatus] = useState(true)
  const [bdtime, setBdtime] = useState(0)
  const [ytstatus, setYtstatus] = useState(true)
  const [yttime, setYttime] = useState(0)
  const [ghstatus, setGhstatus] = useState(true)
  const [ghtime, setGhtime] = useState(0)
  const [current, setCurrent] = useState("未设定")
  const globalWs = useWsContext()
  useEffect(() => {
    const wsuuid = localStorage.getItem("uuid")
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
    const sendData = JSON.stringify({'type': 'testspeed', 'uuid': wsuuid, 'data': 'active'})
    if (globalWs.readyState === 1) globalWs.send(sendData)
    async function getData() {
      const d = await httpServer(urilist.getstatus)
      if (d.status === 0) {
        setStatus(true)
        setCurrent(d.current)
        localStorage.setItem('current', d.current)
      }else{
        setStatus(false)
      }
      setLoading(false)
    }
    getData()
  }, [globalWs])
  const StatusTime = (timestamp) => {
    if (timestamp > 0) {
      return `用时: ${timestamp}秒`
    }else{
      return "测试失败"
    }
  }
  return (
    <>
      {
        loading ? <Loading /> : 
        <div className="tile is-ancestor mt-3">
          <div className="tile is-vertical is-8">
            <div className="tile">
              <div className="tile is-parent is-vertical">
                <article className={"tile is-child notification" + (status ? " is-primary" : " is-danger")}>
                  <p className="title">状态：</p>
                  <p className="subtitle">{status? "运行中" : "未运行"}</p>
                  <p>当前节点：{current}</p>
                </article>
                <article className="tile is-child notification is-warning">
                  <p className="title">Google</p>
                  <p className="subtitle">{ggstatus ? "loading..." : StatusTime(ggtime)}</p>
                </article>
              </div>
              <div className="tile is-parent">
                <article className="tile is-child notification is-info">
                  <p className="title">百度</p>
                  <p className="subtitle">{bdstatus ? "loading..." : StatusTime(bdtime)}</p>
                  <figure className="image is-4by3">
                  </figure>
                </article>
              </div>
            </div>
            <div className="tile is-parent">
              <article className="tile is-child notification is-danger">
                <p className="title">Youtube</p>
                <p className="subtitle">{ytstatus ? "loading..." : StatusTime(yttime)}</p>
                <div className="content">
                </div>
              </article>
            </div>
          </div>
          <div className="tile is-parent">
            <article className="tile is-child notification is-success">
              <div className="content">
                <p className="title">Github</p>
                <p className="subtitle">{ghstatus ? "loading..." : StatusTime(ghtime)}</p>
                <div className="content">
                </div>
              </div>
            </article>
          </div>
        </div>
      }
    </>
  )
}
export default Index