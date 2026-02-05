import { useState, useCallback } from 'react';

export const useSpeedTestHandler = () => {
  const [ggstatus, setGgstatus] = useState(true);
  const [ggtime, setGgtime] = useState(0);
  const [bdstatus, setBdstatus] = useState(true);
  const [bdtime, setBdtime] = useState(0);
  const [ytstatus, setYtstatus] = useState(true);
  const [yttime, setYttime] = useState(0);
  const [ghstatus, setGhstatus] = useState(true);
  const [ghtime, setGhtime] = useState(0);

  // 重置状态的函数（用于开启新测试前）
  const resetSpeedStatus = useCallback(() => {
    setGgstatus(true); setGgtime(0);
    setBdstatus(true); setBdtime(0);
    setYtstatus(true); setYttime(0);
    setGhstatus(true); setGhtime(0);
  }, []);

  const handleSpeedMessage = useCallback((event) => {
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
      console.error("测速数据解析失败", e);
    }
  }, []);

  return {
    speedState: { ggstatus, ggtime, bdstatus, bdtime, ytstatus, yttime, ghstatus, ghtime },
    handleSpeedMessage,
    resetSpeedStatus
  };
};

export default useSpeedTestHandler;