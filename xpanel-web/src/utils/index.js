import { httpServer } from './http'
const Rooturi = window.location.origin + '/'
// const Rooturi = "http://localhost:13005/"
// console.log(Rooturi)
const urilist = {
  'nodelist': `${Rooturi}api/nodelist`,
  'updata': `${Rooturi}api/updata`,
  'setnode': `${Rooturi}api/setnode`,
  'deletenode': `${Rooturi}api/deletenode`,
  'getdomains': `${Rooturi}api/GetDomains`,
  'setdomains': `${Rooturi}api/SetDomains`,
  'getsubscribes':`${Rooturi}api/GetSubscribes`,
  'setSubscribes': `${Rooturi}api/SetSubscribes`,
  'setIgnore': `${Rooturi}api/SetIgnore`,
  'getstatus': `${Rooturi}api/GetStatus`,
  'getdns': `${Rooturi}api/GetDns`,
  'setdns': `${Rooturi}api/SetDns`,
  'testproxy': `${Rooturi}api/TestProxy`,
  'checkVersion': `${Rooturi}api/CheckVersion`,
  'getSocks': `${Rooturi}api/GetLocalSocks`,
  'setLocalSocks': `${Rooturi}api/SetLocalSocks`,
}
export {
  httpServer,
  urilist
}