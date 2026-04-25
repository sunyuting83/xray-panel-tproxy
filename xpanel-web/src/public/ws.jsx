import { createContext, useContext, useEffect, useRef, useState, useCallback } from 'react';

const WsContext = createContext();

// 在组件外定义实例，确保全局唯一
let globalInstance = null;

export const WsProvider = ({ children }) => {
  const [status, setStatus] = useState('CONNECTING');
  const reconnectCount = useRef(0);
  const wsRef = useRef(null);

  // 1. 使用 useCallback 包裹 connect，确保其引用地址稳定
  const connect = useCallback(() => {
    // 防止重复连接
    if (globalInstance && (globalInstance.readyState === WebSocket.OPEN || globalInstance.readyState === WebSocket.CONNECTING)) {
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUri = `${protocol}//${window.location.host}/ws`;
    // const wsUri = `${protocol}//localhost:13005/ws`;
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
        connect(); // 这里会安全地引用 useCallback 包裹后的函数
      }, delay);
    };

    globalInstance = socket;
    wsRef.current = socket;
  }, []); // 依赖项为空，确保 connect 函数永远不会变

  useEffect(() => {
    connect();
    
    const heartbeat = setInterval(() => {
      if (globalInstance?.readyState === WebSocket.OPEN) {
        const uuid = localStorage.getItem("uuid");
        globalInstance.send(JSON.stringify({ type: 'active', uuid, data: 'active' }));
      }
    }, 1000 * 60 * 5);

    return () => clearInterval(heartbeat);
  }, [connect]); // ✅ 现在 connect 已经稳定且加入了依赖数组，警告消失

  // 包装稳定的发送函数
  const sendMessage = (data) => {
    if (globalInstance?.readyState === WebSocket.OPEN) {
      globalInstance.send(typeof data === 'string' ? data : JSON.stringify(data));
    }
  };

  const value = {
    ws: globalInstance,
    status,
    sendMessage,
    addListener: (cb) => globalInstance?.addEventListener('message', cb),
    removeListener: (cb) => globalInstance?.removeEventListener('message', cb),
  };

  return <WsContext.Provider value={value}>{children}</WsContext.Provider>;
};

export const useWsContext = () => useContext(WsContext);