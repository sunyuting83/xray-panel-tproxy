import { useState, useEffect } from 'react'
import {urilist, httpServer} from '../utils/index'
import Loading from '../public/Loading'
import Notification from '../public/Notification'
import { useWsContext } from '../public/ws'

const CheckVersion = () => {
  const globalWs = useWsContext()
  const wsuuid = localStorage.getItem("uuid")
  const [loading, setLoading] = useState(true)
  const [geo, setGeo] = useState(false)
  const [dis, setDis] = useState(false)
  const [geoVersion, setGeoVersion] = useState("")
  const [geoDown, setGeoDown] = useState(false)
  const [geoDownErr, setGeoDownErr] = useState("")
  const [geoprogress, setGeoprogress] = useState("0")
  const [geositeDown, setGeositeDown] = useState(false)
  const [geositeDownErr, setGeositeDownErr] = useState(false)
  const [geositeprogress, setGeositeprogress] = useState("0")
  const [geoIP, setGeoIP] = useState("")
  const [geoSite, setGeoSite] = useState("")
  const [core, setCore] = useState(false)
  const [coreUri, setCoreUri] = useState("")
  const [coreVersion, setCoreVersion] = useState("")
  const [coreName, setCoreName] = useState("")
  const [coreDown, setCoreDown] = useState(false)
  const [coreDownErr, setCoreDownErr] = useState("")
  const [coreDownProgress, setCoreDownProgress] = useState("0")
  const [notification, setNotification] = useState({})
  useEffect(() => {
    async function getData() {
      const d = await httpServer(urilist.checkVersion)
      // console.log(d.CoreVersion)
      setCore(d.CoreVersion)
      if (d.CoreVersion) {
        setCoreUri(d.CoreUri)
        const coreurl = d.CoreUri.split('/')
        const urilen = coreurl.length
        const core_name = coreurl[urilen - 1]
        // console.log(core_name)
        setCoreName(core_name)
        setCoreVersion(d.CurrentCoreVersion)
      }
      setGeo(d.GeoVersion)
      if (d.GeoVersion) {
        setGeoIP(d.GeoIP)
        setGeoSite(d.GeoSite)
        setGeoVersion(d.CurrentGeoVersion)
      }
      setLoading(false)
    }
    getData()
    globalWs.onmessage = (data) => {
      if (typeof JSON.parse(data.data) === "object") {
        const jsonData = JSON.parse(data.data)
        if (jsonData.type === "download") {
          if (jsonData.data.indexOf('----') !== -1) {
            const splitData = jsonData.data.split("----")
            const fileName = splitData[0]
            // console.log(fileName)
            // console.log(coreName)
            switch (fileName) {
              case 'geoip.dat':
                if (!geoDown) setGeoDown(true)
                setGeoprogress(splitData[1])
                break;
              case 'geosite.dat':
                if (!geositeDown) setGeositeDown(true)
                setGeositeprogress(splitData[1])
                break;
              default:
                if (!coreDown) setCoreDown(true)
                setCoreDownProgress(splitData[1])
                break;
            }
          }
        }
        if (jsonData.type === "error") {
          if (jsonData.data.indexOf('----') !== -1) {
            const splitData = jsonData.data.split("----")
            const fileName = splitData[0]
            switch (fileName) {
              case 'geoip.dat':
                if (!geoDown) setGeoDown(true)
                setGeoDownErr(splitData[1])
                break;
              case 'geosite.dat':
                if (!geositeDown) setGeositeDown(true)
                setGeositeDownErr(splitData[1])
                break;
              default:
                if (!coreDown) setCoreDown(true)
                setCoreDownErr(splitData[1])
                break;
            }
          }
        }
      }
    }
  }, [globalWs])
  const CloseNotification = () => {
    setNotification({active: false, message: '', color: ''})
  }
  const UpCore = () => {
    setDis(true)
    setCoreDown(true)
    // console.log(coreName)
    const coredata = `${coreName}||||${coreUri}||||${coreVersion}`
    const sendData = JSON.stringify({'type': 'download', 'uuid': wsuuid, 'data': coredata})
    if (globalWs.readyState === 1) {
      globalWs.send(sendData)
    }
  }
  const UpGeo = () => {
    setDis(true)
    setGeoDown(true)
    setGeositeDown(true)
    const geodata = `geoip.dat||||${geoIP}||||${geoVersion}`
    const sendData = JSON.stringify({'type': 'download', 'uuid': wsuuid, 'data': geodata})
    const geositedata = `geosite.dat||||${geoSite}||||${geoVersion}`
    const geositeSendData = JSON.stringify({'type': 'download', 'uuid': wsuuid, 'data': geositedata})
    if (globalWs.readyState === 1) {
      globalWs.send(sendData)
      setTimeout(() => {
        
      }, 1000)
      globalWs.send(geositeSendData)
    }
  }
  return (
    <>
      {
        loading ? <Loading /> : 
        <div className="box">
          <div className="field">
            <h5 className="title is-5">核心更新</h5>
            <h6 className="subtitle is-6">{core ? '需要更新' : '已经是最新版本'}</h6>
            <div className="control">
              {core?
              <button className="button is-link small" onClick={UpCore} disabled={dis? true : false}>更新核心</button>
              :
              <></>
              }
              {coreDown?
              <div className='columns mt-1'>
                <div className='column is-1'>Core</div>
                <div className='column is-10'>
                  {coreDownErr.length > 0?
                  <h6 className="subtitle is-6">{coreDownErr}</h6>
                  :
                  <progress className="progress is-small is-info" value={coreDownProgress} max="100">{coreDownProgress}%</progress>
                  }
                </div>
              </div>
              :
              <></>}
              
            </div>
          </div>
          <hr />
          <div className="field mt-5">
            <h5 className="title is-5">GeoIP库更新</h5>
            <h6 className="subtitle is-6">{geo ? '需要更新' : '已经是最新版本'}</h6>
            <div className="control">
              {geo?
              <button className="button is-warning is-small" onClick={UpGeo} disabled={dis? true : false}>更新IP库</button>
              :
              <></>}
              {geoDown?
                <div className='columns mt-1'>
                <div className='column is-1'>GeoIP</div>
                <div className='column is-10'>
                  {geoDownErr.length > 0?
                  <h6 className="subtitle is-6">{geoDownErr}</h6>
                  :
                  <progress className="progress is-small is-warning" value={geoprogress} max="100">{geoprogress}%</progress>
                  }
                </div>
                </div>
              :
              <></>}
              {geositeDown?
                <div className='columns mt-1'>
                  <div className='column is-1'>GeoSite</div>
                  <div className='column is-10'>
                    {geositeDownErr.length > 0?
                    <h6 className="subtitle is-6">{geoDownErr}</h6>
                    :
                    <progress className="progress is-small is-warning" value={geositeprogress} max="100">{geositeprogress}%</progress>
                    }
                  </div>
                </div>
              :
              <></>}
            </div>
          </div>
        </div>
      }
      {notification.active ?<Notification props={notification} close={CloseNotification} /> : <></>}
    </>
  )
}
export default CheckVersion