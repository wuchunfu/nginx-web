import {
  deleteRequestById,
  getRequest,
  getRequestById,
  postRequest,
  putRequest,
  putRequestById
} from "/@/utils/requestApi";

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

// 启用
export function enable(fileName: string) {
  return putRequestById(`/sys/website/enable/${ fileName }`);
}

// 禁用
export function disable(fileName: string) {
  return putRequestById(`/sys/website/disable/${ fileName }`);
}

// 删除
export function deleteById(fileName: string) {
  return deleteRequestById(`/sys/website/delete/${ fileName }`);
}
