import { useEffect, useCallback, useState } from 'react'
import { useWsContext } from '../public/ws'
import { useSpeedTestHandler } from '../hooks/useSpeedTestHandler'
import Loading from '../public/Loading'
import { urilist, httpServer } from '../utils/index';

const Index = () => {
  const [custatus, setStatus] = useState(false)
  const [current, setCurrent] = useState("未设定")
  // 从自定义 Hook 中获取状态和处理函数
  const { speedState, handleSpeedMessage } = useSpeedTestHandler()
  // 从 Context 中获取 WebSocket 实例和状态
  const { ws: globalWs, sendMessage, status, addListener, removeListener } = useWsContext()

  // 定义测速辅助函数
  const runSpeedTest = useCallback(() => {
    if (status === 'OPEN') {
      const wsuuid = localStorage.getItem("uuid")
      sendMessage({
        type: 'testspeed',
        uuid: wsuuid,
        data: 'active'
      })
    }
  }, [status, sendMessage])

  useEffect(() => {
    if (!globalWs) return

     // 获取 HTTP 状态
    async function getData() {
      try {
        const d = await httpServer(urilist.getstatus);
        if (d?.status === 0) {
          setStatus(true);
          setCurrent(d.current_title);
          localStorage.setItem('current', d.current);
        } else {
          setStatus(false);
        }
      } catch (err) {
        console.error("HTTP获取失败", err);
      }
    }
    getData();

    // 绑定监听：收到消息时交给 Hook 处理
    const bridge = (e) => handleSpeedMessage(e)
    addListener(bridge)

    // 如果 WS 已经连接，加载页面时直接运行一次测速
    if (status === 'OPEN') {
      runSpeedTest()
    }

    return () => {
      // 卸载时移除监听
      removeListener(bridge)
    }
  }, [globalWs, status, addListener, removeListener, handleSpeedMessage, runSpeedTest])

  // 辅助渲染：格式化显示时间
  const renderTime = (time, loadingStatus) => {
    if (loadingStatus) return <span className="tag is-light">检测中...</span>
    return time > 0 ? (
      <span className="tag is-success is-light">{time}s</span>
    ) : (
      <span className="tag is-danger is-light">超时/失败</span>
    )
  }

  return (
    <div className="container mt-5">
      <article className={"tile is-child notification" + (custatus ? " is-primary" : " is-danger")}>
        <p className="title">状态：</p>
        <p className="subtitle">{custatus ? "运行中" : "未运行"}</p>
        <p>当前节点：{current}</p>
      </article>
      <div className="columns is-multiline">
        {/* 测速卡片展示 */}
        <div className="column is-12">
          <div className="box">
            <div className="level">
              <div className="level-left">
                <h3 className="title is-4">网络连通性状态</h3>
              </div>
              <div className="level-right">
                <button 
                  className={`button is-info is-small ${status !== 'OPEN' ? 'is-loading' : ''}`} 
                  onClick={runSpeedTest}
                  disabled={status !== 'OPEN'}
                >
                  重新测试
                </button>
              </div>
            </div>

            <div className="columns is-mobile is-multiline">
              <div className="column is-half-mobile is-one-quarter-desktop">
                <div className="notification is-white has-text-centered border-light">
                  <p className="heading">百度 (Baidu)</p>
                  <p className="title is-5">{renderTime(speedState.bdtime, speedState.bdstatus)}</p>
                </div>
              </div>
              
              <div className="column is-half-mobile is-one-quarter-desktop">
                <div className="notification is-white has-text-centered border-light">
                  <p className="heading">谷歌 (Google)</p>
                  <p className="title is-5">{renderTime(speedState.ggtime, speedState.ggstatus)}</p>
                </div>
              </div>

              <div className="column is-half-mobile is-one-quarter-desktop">
                <div className="notification is-white has-text-centered border-light">
                  <p className="heading">GitHub</p>
                  <p className="title is-5">{renderTime(speedState.ghtime, speedState.ghstatus)}</p>
                </div>
              </div>

              <div className="column is-half-mobile is-one-quarter-desktop">
                <div className="notification is-white has-text-centered border-light">
                  <p className="heading">YouTube</p>
                  <p className="title is-5">{renderTime(speedState.yttime, speedState.ytstatus)}</p>
                </div>
              </div>
            </div>
            
            <p className="is-size-7 has-text-grey mt-3">
              * 测速基于当前已激活的节点。如果测试全部失败，请检查节点是否可用。
            </p>
          </div>
        </div>
      </div>
      
      {/* 如果 WS 连接还没准备好，可以显示一个小的加载提示 */}
      {status !== 'OPEN' && <Loading />}
    </div>
  )
}

export default Index