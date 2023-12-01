import Header from './pages/Header'
import { Route, Routes } from 'react-router-dom'
import { HistoryRouter, history } from './utils/history'
import { WsProvider } from './public/ws'

import { lazy, Suspense } from 'react'
import Loading from './public/Loading'

const Index = lazy(() => import('./pages/Index'))
const NodeList = lazy(() => import('./pages/Nodelist'))
const White = lazy(() => import('./pages/White'))
const Subscribe = lazy(() => import('./pages/Subscribe'))
const SetDns = lazy(() => import('./pages/SetDns'))
const CheckVersion = lazy(() => import('./pages/CheckVersion'))
const SetSocks = lazy(() => import('./pages/SetSocks'))
function App() {
  return (
    <div className="App">
      <WsProvider>
        <HistoryRouter history={history}>
          <Header />
          <Suspense
            fallback={
              <div
                style={{
                  textAlign: 'center',
                  marginTop: 200
                }}
              >
                <Loading />
              </div>
            }
          >
          <Routes>
            <Route path='/' element={<Index />} />
            <Route path='/nodelist' element={<NodeList />} />
            <Route path='/white' element={<White />} />
            <Route path='/subscribe' element={<Subscribe />} />
            <Route path='/setdns' element={<SetDns />} />
            <Route path='/update' element={<CheckVersion />} />
            <Route path='/setsocks' element={<SetSocks />} />
          </Routes>
          </Suspense>
        </HistoryRouter>
      </WsProvider>
    </div>
  );
}

export default App;
