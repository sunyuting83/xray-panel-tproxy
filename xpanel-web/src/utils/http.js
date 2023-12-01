const httpServer = (url, params = {}, method = 'GET',) => {
    method = method.toUpperCase()
    if (method === 'GET') {
      let dataStr = ''
      Object.keys(params).forEach(key => {
        dataStr += key + '=' + params[key] + '&'
      })
  
      if (dataStr !== '') {
        dataStr = dataStr.substr(0, dataStr.lastIndexOf('&'))
        url = url + '?' + dataStr
      }
    }
    let requestConfig = {
      method: method,
    }
    if (method === 'POST' || method === 'PUT' || method === 'DELETE') {
      if (!Object.hasOwnProperty.call(params, "files")) {
        requestConfig.headers = {
          Accept: '*/*',
          'Content-Type': 'application/json;charset=UTF-8',
        }
        Object.defineProperty(requestConfig, 'body', {
          value: JSON.stringify(params)
        })
      }else {
        requestConfig.headers = {
          Accept: '*/*',
          'Content-Type': 'multipart/form-data',
        }
        requestConfig.processData = false
        const formData = new FormData()
        for (const key in params) {
          formData.append(key,params[key])
        }
        // Object.defineProperty(requestConfig, 'body', formData)
        requestConfig.body = formData
      }
    }
    return new Promise((resolve) => {
      fetch(url, requestConfig)
        .then(res => {
          if(res.ok) {
            return res.json()
          }else {
            resolve({
              status: 1,
              message: "访问出错"
            })
          }
        })
        .then(json => {
          resolve(json)
        })
        .catch((err) => {
          resolve({
            status: 1,
            message: err.message
          })
        })
    })
}

export { httpServer }