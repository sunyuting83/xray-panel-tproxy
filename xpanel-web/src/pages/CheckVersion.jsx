import { useState, useEffect, useCallback, useRef } from 'react'
import { useWsContext } from '../public/ws'
import { urilist, httpServer } from '../utils/index'
import Loading from '../public/Loading'
import Notification from '../public/Notification'

const CheckVersion = () => {
  const { ws: globalWs, sendMessage, status, addListener, removeListener } = useWsContext()
  const [loading, setLoading] = useState(true)
  const [dis, setDis] = useState(false)
  const [notification, setNotification] = useState({})

  // Geo 状态
  const [geo, setGeo] = useState(false)
  const [geoVersion, setGeoVersion] = useState("")
  const [geoIP, setGeoIP] = useState("")
  const [geoSite, setGeoSite] = useState("")
  const [geoDown, setGeoDown] = useState(false)
  const [geoDownErr, setGeoDownErr] = useState("")
  const [geoprogress, setGeoprogress] = useState("0")

  // GeoSite 状态
  const [geositeDown, setGeositeDown] = useState(false)
  const [geositeDownErr, setGeositeDownErr] = useState("")
  const [geositeprogress, setGeositeprogress] = useState("0")

  // Core 状态
  const [core, setCore] = useState(false)
  const [coreUri, setCoreUri] = useState("")
  const [coreVersion, setCoreVersion] = useState("")
  const [coreName, setCoreName] = useState("")
  const [coreDown, setCoreDown] = useState(false)
  const [coreDownErr, setCoreDownErr] = useState("")
  const [coreDownProgress, setCoreDownProgress] = useState("0")

  // 定时器引用，用于清理
  const timerRef = useRef(null)
  // 状态引用，确保 handleMessage 里的逻辑能拿到最新的状态值
  const stateRef = useRef({})
  stateRef.current = { geoSite, geoVersion, status }

  // 1. 定义 UpGeoSite (使用最新的 stateRef)
  const UpGeoSite = useCallback(() => {
    console.log("执行 UpGeoSite (延迟后)");
    const wsuuid = localStorage.getItem("uuid");
    setDis(true);
    setGeositeDown(true);
    
    // 从 stateRef 拿数据，防止闭包拿不到最新的 geoSite
    const { geoSite: latestSite, geoVersion: latestVer, status: latestStatus } = stateRef.current;
    
    const geositedata = `geosite.dat||||${latestSite}||||${latestVer}`;
    if (latestStatus === 'OPEN') {
      sendMessage({ type: 'download', 'uuid': wsuuid, data: geositedata });
    }
  }, [sendMessage]);

  // 2. 核心消息处理
  const handleMessage = useCallback((event) => {
    try {
      const jsonData = JSON.parse(event.data);
      if (jsonData.type === "download" && jsonData.data.includes('----')) {
        const [fileName, progress] = jsonData.data.split("----");

        switch (fileName) {
          case 'geoip.dat':
            setGeoDown(true);
            setGeoprogress(progress);
            if (progress === "100") {
              console.log("GeoIP 下载完成，6秒后开始更新 GeoSite...");
              setGeositeDown(true);
              setGeositeDownErr("GeoIP 下载完成，6秒后开始更新 GeoSite...,切勿关闭或切换本页面")
              // 清除之前的定时器，防止重复触发
              if (timerRef.current) clearTimeout(timerRef.current);
              // 设置 1500ms 延迟
              timerRef.current = setTimeout(() => {
                UpGeoSite();
              }, 6000);
            }
            break;
          case 'geosite.dat':
            setGeositeDownErr("")
            setGeositeprogress(progress)
            break;
          default:
            setCoreDown(true);
            setCoreDownProgress(progress)
            break;
        }
      }
      
      // 错误处理逻辑
      if (jsonData.type === "error" && jsonData.data.includes('----')) {
        const [fileName, errMsg] = jsonData.data.split("----");
        switch (fileName) {
          case 'geoip.dat': setGeoDownErr(errMsg); break;
          case 'geosite.dat': setGeositeDownErr(errMsg); break;
          default: setCoreDownErr(errMsg); break;
        }
      }
    } catch (e) {
      console.error("解析下载消息失败", e);
    }
  }, [UpGeoSite]);

  // 3. 监听与初始化
  useEffect(() => {
    if (!globalWs) return;

    const bridge = (e) => handleMessage(e);
    addListener(bridge);

    async function getData() {
      try {
        const d = await httpServer(urilist.checkVersion);
        setCore(d.CoreVersion);
        if (d.CoreVersion) {
          setCoreUri(d.CoreUri);
          const parts = d.CoreUri.split('/');
          setCoreName(parts[parts.length - 1]);
          setCoreVersion(d.CurrentCoreVersion);
        }
        setGeo(d.GeoVersion);
        if (d.GeoVersion) {
          setGeoIP(d.GeoIP);
          setGeoSite(d.GeoSite);
          setGeoVersion(d.CurrentGeoVersion);
        }
      } catch (err) { console.error(err); }
      finally { setLoading(false); }
    }
    getData();

    return () => {
      removeListener(bridge);
      // 组件卸载时，一定要清理定时器，防止报错
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, [globalWs, addListener, removeListener, handleMessage]);

  // 手动触发函数
  const UpCore = () => {
    setDis(true); setCoreDown(true);
    const wsuuid = localStorage.getItem("uuid");
    const coredata = `${coreName}||||${coreUri}||||${coreVersion}`;
    if (status === 'OPEN') sendMessage({ type: 'download', 'uuid': wsuuid, data: coredata });
  };

  const UpGeo = () => {
    setDis(true); setGeoDown(true);
    const wsuuid = localStorage.getItem("uuid");
    const geodata = `geoip.dat||||${geoIP}||||${geoVersion}`;
    if (status === 'OPEN') sendMessage({ type: 'download', 'uuid': wsuuid, data: geodata });
  };

  const CloseNotification = () => setNotification({ active: false, message: '', color: '' });

  return (
    <>
      {loading ? <Loading /> : (
        <div className="box">
          {/* 核心更新 */}
          <div className="field">
            <h5 className="title is-5">核心更新</h5>
            <h6 className="subtitle is-6">{core ? '需要更新' : '已经是最新版本'}</h6>
            <div className="control">
              {core && <button className="button is-link is-small" onClick={UpCore} disabled={dis}>更新核心</button>}
              {coreDown && (
                <div className='columns mt-1'>
                  <div className='column is-1'>Core</div>
                  <div className='column is-10'>
                    {coreDownErr ? <h6 className="help is-danger">{coreDownErr}</h6> : 
                    <progress className="progress is-small is-info" value={coreDownProgress} max="100">{coreDownProgress}%</progress>}
                  </div>
                </div>
              )}
            </div>
          </div>
          <hr />
          {/* Geo 数据更新 */}
          <div className="field mt-5">
            <h5 className="title is-5">Geo数据更新</h5>
            <h6 className="subtitle is-6">{geo ? '需要更新' : '已经是最新版本'}</h6>
            <div className="control">
              {geo && <button className="button is-warning is-small" onClick={UpGeo} disabled={dis}>更新Geo库</button>}
              {geoDown && (
                <div className='columns mt-1'>
                  <div className='column is-1'>GeoIP</div>
                  <div className='column is-10'>
                    {geoDownErr ? <h6 className="help is-danger">{geoDownErr}</h6> : 
                    <progress className="progress is-small is-warning" value={geoprogress} max="100">{geoprogress}%</progress>}
                  </div>
                </div>
              )}
              {geositeDown && (
                <div className='columns mt-1'>
                  <div className='column is-1'>GeoSite</div>
                  <div className='column is-10'>
                    {geositeDownErr ? <h6 className="help is-danger">{geositeDownErr}</h6> : 
                    <progress className="progress is-small is-primary" value={geositeprogress} max="100">{geositeprogress}%</progress>}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
      {notification.active && <Notification props={notification} close={CloseNotification} />}
    </>
  )
}

export default CheckVersion