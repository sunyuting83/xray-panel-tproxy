import { createContext, useContext, useEffect, useRef, useState } from 'react';

const WsContext = createContext();

// 在组件外定义实例，确保全局唯一
let globalInstance = null;

export const WsProvider = ({ children }) => {
  const [status, setStatus] = useState('CONNECTING');
  const reconnectCount = useRef(0);
  const wsRef = useRef(null);

  const connect = () => {
    // 防止重复连接
    if (globalInstance && (globalInstance.readyState === WebSocket.OPEN || globalInstance.readyState === WebSocket.CONNECTING)) {
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUri = `${protocol}//${window.location.host}/ws`;
    const socket = new WebSocket(wsUri);

    socket.onopen = () => {
      console.log("✅ WS Connected");
      setStatus('OPEN');
      reconnectCount.current = 0;
    };

    socket.onmessage = (event) => {
      try {
        const jsonData = JSON.parse(event.data);
        if (jsonData.uuid) localStorage.setItem("uuid", jsonData.uuid);
      } catch (e) {}
    };

    socket.onclose = () => {
      setStatus('CLOSED');
      console.log("❌ WS Closed. Reconnecting...");
      globalInstance = null; // 清除实例允许重连
      const delay = Math.min(1000 * 2 ** reconnectCount.current, 30000);
      setTimeout(() => {
        reconnectCount.current++;
        connect();
      }, delay);
    };

    globalInstance = socket;
    wsRef.current = socket;
  };

  useEffect(() => {
    connect();
    const heartbeat = setInterval(() => {
      if (globalInstance?.readyState === WebSocket.OPEN) {
        const uuid = localStorage.getItem("uuid");
        globalInstance.send(JSON.stringify({ type: 'active', uuid, data: 'active' }));
      }
    }, 1000 * 60 * 5);

    return () => clearInterval(heartbeat);
  }, []);

  // 包装稳定的发送函数
  const sendMessage = (data) => {
    if (globalInstance?.readyState === WebSocket.OPEN) {
      globalInstance.send(typeof data === 'string' ? data : JSON.stringify(data));
    }
  };

  const value = {
    // 导出实例和方法
    ws: globalInstance,
    status,
    sendMessage,
    // 专家建议：直接通过 globalInstance 操作，避免 context 刷新导致的引用丢失
    addListener: (cb) => globalInstance?.addEventListener('message', cb),
    removeListener: (cb) => globalInstance?.removeEventListener('message', cb),
  };

  return <WsContext.Provider value={value}>{children}</WsContext.Provider>;
};

export const useWsContext = () => useContext(WsContext);