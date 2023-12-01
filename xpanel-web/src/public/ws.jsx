import { createContext, useContext } from 'react'

const WsContext = createContext()

export const WsProvider = ({ children }) => {
  // const wsUri = `ws://localhost:13005/ws`
  const a = window.location.origin
  const b = a.split("://")[1]
  const wsUri = `ws://${b}/ws`
  const globalWs = new WebSocket(wsUri)

  return (
    <WsContext.Provider value={globalWs}>
      {children}
    </WsContext.Provider>
  )
}

export const useWsContext = () => useContext(WsContext)