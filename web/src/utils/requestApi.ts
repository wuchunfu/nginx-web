import request from '/@/utils/request'

export const getRequest = (url: string, params: object) => {
  return request({
    method: 'get',
    url: url,
    params: params
  })
}

export const getRequestById = (url: string) => {
  return request({
    method: 'get',
    url: url
  })
}

export const postRequest = (url: string, data: object) => {
  return request({
    method: 'post',
    url: url,
    data: data
  })
}

export const postRequestById = (url: string) => {
  return request({
    method: 'post',
    url: url
  })
}

export const putRequest = (url: string, data: object) => {
  return request({
    method: 'put',
    url: url,
    data: data
  })
}

export const putRequestById = (url: string) => {
  return request({
    method: 'put',
    url: url,
  })
}

export const deleteRequestById = (url: string) => {
  return request({
    method: 'delete',
    url: url
  })
}

export const uploadFileRequest = (url: string, data: object) => {
  return request({
    method: 'post',
    url: url,
    data: data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
