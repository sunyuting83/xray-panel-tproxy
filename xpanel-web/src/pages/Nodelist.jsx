import { useState, useEffect, useCallback, useRef } from 'react'
import { urilist, httpServer } from '../utils/index'
import Loading from '../public/Loading'
import EmptyEd from '../public/Empty'
import Notification from '../public/Notification'
import { useWsContext } from '../public/ws'

const NodeList = () => {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [notification, setNotification] = useState({})
  const [modal, setModal] = useState({})
  const [current, setCurrent] = useState("")
  const [uuid, setUUID] = useState("")
  
  const [ggstatus, setGgstatus] = useState(true)
  const [ggtime, setGgtime] = useState(0)
  const [bdstatus, setBdstatus] = useState(true)
  const [bdtime, setBdtime] = useState(0)
  const [ytstatus, setYtstatus] = useState(true)
  const [yttime, setYttime] = useState(0)
  const [ghstatus, setGhstatus] = useState(true)
  const [ghtime, setGhtime] = useState(0)
  const [tcping, setTcping] = useState(false)

  const { ws: globalWs, sendMessage, status, addListener, removeListener } = useWsContext()

  // 使用 Ref 解决闭包问题，确保 handleMessage 永远能拿到最新的 state
  const stateRef = useRef();
  stateRef.current = { data, uuid, tcping };

  const handleMessage = useCallback((event) => {
    try {
      const jsonData = JSON.parse(event.data);
      
      // 逻辑A: 处理连通性测速 (testspeed)
      if (jsonData.type === "testspeed" && jsonData.data.includes('||||')) {
        const [url, speedStr] = jsonData.data.split("||||");
        const speedNumber = parseFloat(speedStr);
        switch (url) {
          case 'https://www.google.com.hk':
            if (speedNumber > 0) setGgtime(speedNumber);
            setGgstatus(false);
            break;
          case 'https://baidu.com':
            if (speedNumber > 0) setBdtime(speedNumber);
            setBdstatus(false);
            break;
          case 'https://www.youtube.com/img/desktop/yt_1200.png':
            if (speedNumber > 0) setYttime(speedNumber);
            setYtstatus(false);
            break;
          case 'https://github.com/webgl-globe/data/data.json':
            if (speedNumber > 0) setGhtime(speedNumber);
            setGhstatus(false);
            break;
          default:
            break;
        }
      }

      // 逻辑B: 处理 Ping 测速 (tcping)
      if (jsonData.type === "tcping" && jsonData.data.includes('||||')) {
        const [indexStr, speedStr] = jsonData.data.split("||||");
        const i = parseInt(indexStr);
        const speedNumber = parseFloat(speedStr);

        setData(prevData => prevData.map((e, index) => {
          if (i === index) return { ...e, ping: speedNumber };
          return e;
        }));

        if (i === stateRef.current.data.length - 1) {
          setTcping(false);
        }
      }
    } catch (e) {
      console.error("解析消息失败", e);
    }
  }, []);

  useEffect(() => {
    if (!globalWs) return;

    // 绑定监听
    const bridge = (e) => handleMessage(e);
    addListener(bridge);

    const wsuuid = localStorage.getItem("uuid");
    setUUID(wsuuid);

    async function getData() {
      try {
        const d = await httpServer(urilist.nodelist);
        setData(d.date || []);
        const node = localStorage.getItem('current');
        setCurrent(node);
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    }
    getData();

    return () => removeListener(bridge);
  }, [globalWs, addListener, removeListener, handleMessage]);

  const OpenNotification = (message, color) => {
    setNotification({ active: true, message: message, color: color })
  }

  const OpenModal = () => {
    setModal({ active: true })
    if (status === 'OPEN') {
      sendMessage({
        type: 'testspeed',
        uuid: uuid || localStorage.getItem("uuid"),
        data: 'active'
      })
    }
  }

  const CloseModal = () => {
    setGgtime(0); setGgstatus(true);
    setBdtime(0); setBdstatus(true);
    setYttime(0); setYtstatus(true);
    setGhtime(0); setGhstatus(true);
    setModal({ active: false })
  }

  const CloseNotification = () => {
    setNotification({ active: false, message: '', color: '' })
  }

  const TCPing = () => {
    if (status !== 'OPEN') return;
    setTcping(true);
    sendMessage({
      type: 'tcping',
      uuid: uuid || localStorage.getItem("uuid"),
      data: 'active'
    });
  }

  const DeleteNode = async (i) => {
    const d = await httpServer(urilist.deletenode, { node: String(i) }, 'DELETE')
    if (d.status === 0) {
      setData(data.filter((_, index) => index !== i))
    }
  }

  const SetNode = async (i, title) => {
    const d = await httpServer(urilist.setnode, { node: String(i) }, 'put')
    if (d.status === 0) {
      localStorage.setItem('current', title)
      setCurrent(title)
      OpenModal()
    } else {
      OpenNotification(d.message, 'danger')
    }
  }

  const ShowSpeed = (speed) => (
    <span>{speed === undefined || speed === "" ? "待测试" : (speed > 0 ? `${speed}ms` : "失败")}</span>
  );

  const StatusTime = (timestamp) => (timestamp > 0 ? `用时: ${timestamp}秒` : "测试失败");

  const RenderModal = () => {
    if (!modal.active) return null;
    return (
      <div className="modal is-active">
        <div className="modal-background"></div>
        <div className="modal-card">
          <header className="modal-card-head">
            <p className="modal-card-title is-size-5">设置成功 正在测试连通性</p>
            <button className="delete" onClick={CloseModal}></button>
          </header>
          <section className="modal-card-body">
            <div className="columns is-mobile is-flex-wrap-wrap">
              <div className="column is-half"><p className="is-size-7">百度: {bdstatus ? "loading..." : StatusTime(bdtime)}</p></div>
              <div className="column is-half"><p className="is-size-7">Google: {ggstatus ? "loading..." : StatusTime(ggtime)}</p></div>
              <div className="column is-half"><p className="is-size-7">Github: {ghstatus ? "loading..." : StatusTime(ghtime)}</p></div>
              <div className="column is-half"><p className="is-size-7">Youtube: {ytstatus ? "loading..." : StatusTime(yttime)}</p></div>
            </div>
          </section>
          <footer className="modal-card-foot">
            <button className="button is-success" onClick={CloseModal}>关闭</button>
          </footer>
        </div>
      </div>
    );
  }

  return (
    <>
      {loading ? <Loading /> : (
        <>
          {data.length === 0 ? <EmptyEd /> : (
            <>
              <div className="columns flex-wrap is-justify-content-space-between mt-1">
                <div className="field ml-3">
                  <div className="buttons are-small">
                    <button className="button is-warning" onClick={TCPing} disabled={tcping || status !== 'OPEN'}>
                      Ping测速
                    </button>
                  </div>
                </div>
              </div>
              <h6 className="subtitle is-6">Ping测速只能说明节点的端口是通的,无法确认连通性。</h6>
              <div className="table-container">
                <table className="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
                  <thead className="is-size-7">
                    <tr><td>序号</td><td>类型</td><td>节点名称</td><td>测速</td><td>操作</td></tr>
                  </thead>
                  <tbody className="is-size-7">
                    {data.map((item, index) => (
                      <tr key={index}>
                        <td>{index + 1}</td>
                        <td>{item.types}</td>
                        <td>{item.title}</td>
                        <td>{ShowSpeed(item.ping)}</td>
                        <td>
                          <div className="buttons">
                            <button className="button is-success is-small" disabled={current === item.title} onClick={() => SetNode(index, item.title)}>
                              {current === item.title ? '当前节点' : '使用节点'}
                            </button>
                            <button className="button is-info is-small" onClick={() => DeleteNode(index)}>删除节点</button>
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
      {notification.active && <Notification props={notification} close={CloseNotification} />}
      <RenderModal />
    </>
  )
}
export default NodeList