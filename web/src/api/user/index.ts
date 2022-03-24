import { getRequest, getRequestById, postRequest, putRequest } from "/@/utils/requestApi";

// 查询列表
export function getList(params: any) {
  return getRequest('/sys/user/list', params);
}

// 获取详情
export function getDetailById(id: any) {
  return getRequestById(`/sys/user/detail/${ id }`);
}

// 新增
export function save(data: object) {
  return postRequest('/sys/user/save', data);
}

// 修改
export function update(data: object) {
  return putRequest('/sys/user/update', data);
}

// 删除
export function deleteById(data: object) {
  return postRequest('/sys/user/delete', data);
}
