import { useState, useEffect, useRef } from 'react';
import { urilist, httpServer } from '../utils/index';
import Loading from '../public/Loading';
import { useWsContext } from '../public/ws';

const Index = () => {
  const [status, setStatus] = useState(false);
  const [loading, setLoading] = useState(true);
  const [ggstatus, setGgstatus] = useState(true);
  const [ggtime, setGgtime] = useState(0);
  const [bdstatus, setBdstatus] = useState(true);
  const [bdtime, setBdtime] = useState(0);
  const [ytstatus, setYtstatus] = useState(true);
  const [yttime, setYttime] = useState(0);
  const [ghstatus, setGhstatus] = useState(true);
  const [ghtime, setGhtime] = useState(0);
  const [current, setCurrent] = useState("未设定");

  const { sendMessage, status: wsstatus, addListener, removeListener } = useWsContext();

  // --- 关键重构：使用 Ref 解决闭包陷阱 ---
  const logicRef = useRef();

  // 每次渲染都更新 Ref，确保 handleMessage 逻辑永远是最新的
  logicRef.current = (event) => {
    try {
      const jsonData = JSON.parse(event.data);
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
    } catch (e) {
      console.error("解析失败", e);
    }
  };

  useEffect(() => {
    // 一个稳定的桥接函数，它不依赖组件状态，只调用 Ref
    const bridge = (e) => logicRef.current?.(e);

    // 绑定监听器
    addListener(bridge);

    // 发送初始测速请求
    if (wsstatus === 'OPEN') {
      sendMessage({
        type: 'testspeed',
        uuid: localStorage.getItem("uuid"),
        data: 'active'
      });
    }

    // 获取 HTTP 状态
    async function getData() {
      try {
        const d = await httpServer(urilist.getstatus);
        if (d?.status === 0) {
          setStatus(true);
          setCurrent(d.current);
          localStorage.setItem('current', d.current);
        } else {
          setStatus(false);
        }
      } catch (err) {
        console.error("HTTP获取失败", err);
      } finally {
        setLoading(false);
      }
    }
    getData();

    // 清理：组件销毁时移除监听
    return () => {
      removeListener(bridge);
    };
  }, [wsstatus, sendMessage, addListener, removeListener]); // 移除 globalWs 依赖

  // 处理 WS 连接中的状态
  if (wsstatus !== 'OPEN') {
    return <Loading />;
  }

  const StatusTime = (timestamp) => {
    return timestamp > 0 ? `用时: ${timestamp}秒` : "测试失败";
  };

  return (
    <>
      {loading ? <Loading /> : (
        <div className="tile is-ancestor mt-3">
          {/* ... 你的渲染内容保持不变 ... */}
          <div className="tile is-vertical is-8">
            <div className="tile">
              <div className="tile is-parent is-vertical">
                <article className={"tile is-child notification" + (status ? " is-primary" : " is-danger")}>
                  <p className="title">状态：</p>
                  <p className="subtitle">{status ? "运行中" : "未运行"}</p>
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
                </article>
              </div>
            </div>
            <div className="tile is-parent">
              <article className="tile is-child notification is-danger">
                <p className="title">Youtube</p>
                <p className="subtitle">{ytstatus ? "loading..." : StatusTime(yttime)}</p>
              </article>
            </div>
          </div>
          <div className="tile is-parent">
            <article className="tile is-child notification is-success">
              <div className="content">
                <p className="title">Github</p>
                <p className="subtitle">{ghstatus ? "loading..." : StatusTime(ghtime)}</p>
              </div>
            </article>
          </div>
        </div>
      )}
    </>
  );
};

export default Index;