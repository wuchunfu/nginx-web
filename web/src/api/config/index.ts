import { getRequest, putRequest } from "/@/utils/requestApi";

// 查询列表
export function getList(params: object) {
  return getRequest('/sys/config/list', params);
}

// 获取文件夹下面的文件以及文件夹
export function getListFolder(params: object) {
  return getRequest('/sys/config/changeFolder', params);
}

// 获取详情
export function getDetail(params: object) {
  return getRequest('/sys/config/detail', params);
}

// 修改
export function update(data: object) {
  return putRequest('/sys/config/update', data);
}
