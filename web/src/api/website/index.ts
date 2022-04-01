import { getRequest, getRequestById, postRequest, putRequest } from "/@/utils/requestApi";

// 查询列表
export function getList(params: object) {
  return getRequest('/sys/website/list', params);
}

// 获取详情
export function getDetailById(fileName: string) {
  return getRequestById(`/sys/website/detail/${ fileName }`);
}

// 新增
export function save(data: object) {
  return postRequest('/sys/website/save', data);
}

// 修改
export function update(data: object) {
  return putRequest('/sys/website/update', data);
}
