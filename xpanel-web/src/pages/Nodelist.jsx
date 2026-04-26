import { useState, useEffect, useCallback, useRef } from 'react'
import { urilist, httpServer } from '../utils/index'
import Loading from '../public/Loading'
import EmptyEd from '../public/Empty'
import Notification from '../public/Notification'
import { useWsContext } from '../public/ws'
import { useSpeedTestHandler } from '../hooks/useSpeedTestHandler'

const NodeList = () => {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [notification, setNotification] = useState({})
  const [modal, setModal] = useState({ active: false })
  const [current, setCurrent] = useState("")
  const [tcpingActive, setTcpingActive] = useState(false)

  // 1. 引入共享的测速 Hook
  const { speedState, handleSpeedMessage, resetSpeedStatus } = useSpeedTestHandler()
  const { ws: globalWs, sendMessage, status, addListener, removeListener } = useWsContext()

  // 使用 Ref 处理闭包和定时器清理
  const dataRef = useRef()
  dataRef.current = data
  
  // 用于管理 OpenModal 中测速请求的定时器
  const speedTestTimerRef = useRef(null)

  // 2. 定义处理函数
  const onMessageReceived = useCallback((event) => {
    handleSpeedMessage(event)

    try {
      const jsonData = JSON.parse(event.data)
      if (jsonData.type === "tcping" && jsonData.data.includes('||||')) {
        const [indexStr, speedStr] = jsonData.data.split("||||")
        const i = parseInt(indexStr)
        const speedNumber = parseFloat(speedStr)

        setData(prevData => prevData.map((item, index) => {
          if (i === index) return { ...item, ping: speedNumber }
          return item
        }))

        if (i === dataRef.current.length - 1) {
          setTcpingActive(false)
        }
      }
    } catch (e) {
      console.error("解析消息失败", e)
    }
  }, [handleSpeedMessage])

  // 3. 注册 WebSocket 监听
  useEffect(() => {
    if (!globalWs) return
    
    const bridge = (e) => onMessageReceived(e)
    addListener(bridge)

    const fetchData = async () => {
      try {
        const d = await httpServer(urilist.nodelist)
        setData(d.date || [])
        setCurrent(localStorage.getItem('current'))
      } catch (err) {
        console.error(err)
      } finally {
        setLoading(false)
      }
    }
    fetchData()

    return () => {
      removeListener(bridge)
      // ✅ 卸载组件时清理定时器，防止内存泄漏
      if (speedTestTimerRef.current) clearTimeout(speedTestTimerRef.current)
    }
  }, [globalWs, addListener, removeListener, onMessageReceived])

  // --- 业务操作函数 ---

  const OpenNotification = (message, color) => {
    setNotification({ active: true, message, color })
  }

  // ✅ 修改后的 OpenModal 函数：加入 2 秒延迟
  const OpenModal = () => {
    // 1. 先重置状态并打开弹窗界面
    resetSpeedStatus() 
    setModal({ active: true })

    // 2. 清理掉之前可能还没执行的定时器（防止快速切换节点导致的冲突）
    if (speedTestTimerRef.current) clearTimeout(speedTestTimerRef.current)

    // 3. 开启延迟 2 秒的测速请求
    if (status === 'OPEN') {
      speedTestTimerRef.current = setTimeout(() => {
        sendMessage({
          type: 'testspeed',
          uuid: localStorage.getItem("uuid"),
          data: 'active'
        })
        console.log("已延迟 2s 发送连通性测试请求")
      }, 2000)
    }
  }

  const CloseModal = () => {
    // 关闭弹窗时也建议清理定时器，如果用户秒关弹窗，就没必要测速了
    if (speedTestTimerRef.current) clearTimeout(speedTestTimerRef.current)
    setModal({ active: false })
  }

  const TCPing = () => {
    if (status !== 'OPEN') return
    setTcpingActive(true)
    sendMessage({
      type: 'tcping',
      uuid: localStorage.getItem("uuid"),
      data: 'active'
    })
  }

  const DeleteNode = async (uid) => {
    const d = await httpServer(urilist.deletenode, { node: uid }, 'DELETE')
    if (d.status === 0) {
      setData(data.filter((item) => item.uid !== uid))
    }
  }

  const SetNode = async (uid, title) => {
    const d = await httpServer(urilist.setnode, { node: uid }, 'put')
    if (d.status === 0) {
      localStorage.setItem('current', uid)
      localStorage.setItem('current_title', title)
      setCurrent(uid)
      OpenModal() // 这里会触发延迟测速
    } else {
      OpenNotification(d.message, 'danger')
    }
  }

  // --- UI 辅助函数 ---

  const StatusTimeText = (time, isLoading) => {
    if (isLoading) return "检测中..."
    return time > 0 ? `用时: ${time}秒` : "测试失败"
  }

  const RenderModal = () => {
    if (!modal.active) return null
    return (
      <div className="modal is-active">
        <div className="modal-background" onClick={CloseModal}></div>
        <div className="modal-card">
          <header className="modal-card-head">
            <p className="modal-card-title is-size-5">设置成功 正在测试连通性</p>
            <button className="delete" onClick={CloseModal}></button>
          </header>
          <section className="modal-card-body">
            <div className="columns is-mobile is-multiline">
              <div className="column is-half">
                <p className="is-size-7">百度: {StatusTimeText(speedState.bdtime, speedState.bdstatus)}</p>
              </div>
              <div className="column is-half">
                <p className="is-size-7">Google: {StatusTimeText(speedState.ggtime, speedState.ggstatus)}</p>
              </div>
              <div className="column is-half">
                <p className="is-size-7">Github: {StatusTimeText(speedState.ghtime, speedState.ghstatus)}</p>
              </div>
              <div className="column is-half">
                <p className="is-size-7">Youtube: {StatusTimeText(speedState.yttime, speedState.ytstatus)}</p>
              </div>
            </div>
          </section>
          <footer className="modal-card-foot">
            <button className="button is-success" onClick={CloseModal}>关闭</button>
          </footer>
        </div>
      </div>
    )
  }

  return (
    <>
      {loading ? <Loading /> : (
        <>
          {data.length === 0 ? <EmptyEd /> : (
            <>
              <div className="columns is-vcentered mt-1 ml-1 mr-1">
                <div className="column">
                  <button 
                    className={`button is-warning is-small ${tcpingActive ? 'is-loading' : ''}`} 
                    onClick={TCPing} 
                    disabled={tcpingActive || status !== 'OPEN'}
                  >
                    Ping测速
                  </button>
                </div>
              </div>
              <h6 className="subtitle is-6 ml-3">Ping测速只能说明节点的端口是通的。测试可用性请使用“使用节点”后的连通性测试。</h6>
              
              <div className="table-container">
                <table className="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
                  <thead className="is-size-7">
                    <tr>
                      <th>序号</th><th>类型</th><th>节点名称</th><th>测速</th><th>操作</th>
                    </tr>
                  </thead>
                  <tbody className="is-size-7">
                    {data.map((item) => (
                      <tr key={item.uid}>
                        <td>{item.index}</td>
                        <td>{item.type}</td>
                        <td>{item.title}</td>
                        <td>
                          {item.ping !== undefined ? (
                            <span className={item.ping > 0 ? "has-text-success" : "has-text-danger"}>
                              {item.ping > 0 ? `${item.ping}ms` : '失败'}
                            </span>
                          ) : <span>待测试</span>}
                        </td>
                        <td>
                          <div className="buttons">
                            <button 
                              className="button is-success is-small" 
                              disabled={current === item.uid} 
                              onClick={() => SetNode(item.uid, item.title)}
                            >
                              {current === item.uid ? '当前节点' : '使用节点'}
                            </button>
                            <button className="button is-info is-small" onClick={() => DeleteNode(item.uid)}>删除</button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </>
          )}
        </>
      )}
      {notification.active && (
        <Notification props={notification} close={() => setNotification({ active: false })} />
      )}
      <RenderModal />
    </>
  )
}

export default NodeList