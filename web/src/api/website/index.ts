import { getRequest } from "/@/utils/requestApi";

// 查询列表
export function getList(params: object) {
  return getRequest('/sys/website/list', params);
}
