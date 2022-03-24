import { postRequest, postRequestById } from "/@/utils/requestApi";

/**
 * 用户登录
 * @param data 要传的参数值
 * @returns 返回接口数据
 */
export function signIn(data: any) {
  const parse = JSON.parse(JSON.stringify(data));
  const params = {
    username: parse.username,
    password: parse.password
  }
  return postRequest('/sys/login', params);
}

/**
 * 用户退出登录
 * @returns 返回接口数据
 */
export function signOut() {
  return postRequestById('/sys/logout');
}
