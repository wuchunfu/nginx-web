import axios, { AxiosRequestHeaders } from 'axios';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Session } from '/@/utils/storage';

axios.defaults.headers.common['Content-Type'] = 'application/json; charset=utf-8'

// 配置新建一个 axios 实例
const service = axios.create({
// `baseURL` will be prepended to `url` unless `url` is absolute.
  // It can be convenient to set `baseURL` for an instance of axios to pass relative URLs
  // to methods of that instance.
  baseURL: import.meta.env.VITE_API_URL as any,
// `timeout` specifies the number of milliseconds before the request times out.
  // If the request takes longer than `timeout`, the request will be aborted.
  timeout: 50000,
  // `transformRequest` allows changes to the request data before it is sent to the server
  // This is only applicable for request methods 'PUT', 'POST', 'PATCH' and 'DELETE'
  // The last function in the array must return a string or an instance of Buffer, ArrayBuffer,
  // FormData or Stream
  // You may modify the headers object.
  transformRequest: [
    (data: any, headers?: AxiosRequestHeaders) => {
      // Do whatever you want to transform the data
      if (headers && headers['Content-Type'] === 'multipart/form-data; charset=UTF-8') {
        return data
      } else {
        if (headers) {
          headers['Content-Type'] = 'application/json'
        }
      }
      return JSON.stringify(data);
    }
  ],
  // `transformResponse` allows changes to the response data to be made before
  // it is passed to then/catch
  transformResponse: [
    (data: any) => {
      // Do whatever you want to transform the data
      return JSON.parse(data);
    }],
  // `headers` 是即将被发送的自定义请求头
  headers: { 'Content-Type': 'application/json' },
  // `responseType` indicates the type of data that the server will respond with
  // options are: 'arraybuffer', 'document', 'json', 'text', 'stream'
  //   browser only: 'blob'
  // default
  responseType: 'json',
  // `responseEncoding` indicates encoding to use for decoding responses (Node.js only)
  // Note: Ignored for `responseType` of 'stream' or client-side requests
  // default
  responseEncoding: 'utf8',
  // `xsrfCookieName` is the name of the cookie to use as a value for xsrf token
  // default
  xsrfCookieName: 'XSRF-TOKEN',
  // `xsrfHeaderName` is the name of the http header that carries the xsrf token value
  // default
  xsrfHeaderName: 'X-XSRF-TOKEN',
  // `withCredentials` indicates whether or not cross-site Access-Control requests
  // should be made using credentials
  // default
  withCredentials: false,
});

// 添加请求拦截器
service.interceptors.request.use((config: any) => {
    // 是否需要设置 token
    const isToken = (config.headers || {}).isToken === false
    // 在发送请求之前做些什么 token
    if (Session.get('token') && !isToken) {
      // 让每个请求携带自定义token 请根据实际情况自行修改
      (<any>config.headers).common['Authorization'] = `${ Session.get('token') }`;
    }
    // get请求映射params参数
    if (config.method === 'get' && config.params) {
      let url = config.url + '?';
      for (const propName of Object.keys(config.params)) {
        const value = config.params[propName];
        const part = encodeURIComponent(propName) + "=";
        if (value !== null && typeof (value) !== "undefined") {
          if (typeof value === 'object') {
            for (const key of Object.keys(value)) {
              if (value[key] !== null && typeof (value[key]) !== 'undefined') {
                let params = propName + '[' + key + ']';
                let subPart = encodeURIComponent(params) + '=';
                url += subPart + encodeURIComponent(value[key]) + '&';
              }
            }
          } else {
            url += part + encodeURIComponent(value) + "&";
          }
        }
      }
      url = url.slice(0, -1);
      config.params = {};
      config.url = url;
    }
    return config;
  }, (error) => {
    console.log(error)
    // 对请求错误做些什么
    return Promise.reject(error);
  }
);

// 添加响应拦截器
service.interceptors.response.use((response) => {
    // 对响应数据做点什么
    const res = response.data;
    if (res.code && res.code !== 200) {
      // token 过期或者账号已在别处登录
      if (res.code === 401 || res.code === 4001) {
        // 清除浏览器全部临时缓存
        Session.clear();
        ElMessageBox.confirm('登录状态已过期，您可以继续留在该页面，或者重新登录', '系统提示', {
          confirmButtonText: '重新登录',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          // 去登录页
          window.location.href = '/';
        }).catch(() => {
        });
        return Promise.reject("认证失败，无法访问系统资源");
      } else if (res.code === 500) {
        ElMessage({
          type: 'error',
          message: '系统未知错误，请反馈给管理员',
          showClose: true,
        })
        return Promise.reject(new Error("系统未知错误，请反馈给管理员"))
      } else if (res.code === 403) {
        ElMessage({
          type: 'error',
          message: '当前操作没有权限',
          showClose: true,
        })
        return Promise.reject(new Error("当前操作没有权限"))
      } else if (res.code === 404) {
        ElMessage({
          type: 'error',
          message: '访问资源不存在',
          showClose: true,
        })
        return Promise.reject(new Error("访问资源不存在"))
      } else {
        if (res.msg !== "") {
          ElMessage({
            type: 'error',
            message: res.msg,
            showClose: true,
          })
        } else {
          ElMessage({
            type: 'error',
            message: '系统未知错误，请反馈给管理员',
            showClose: true,
          })
        }
        return Promise.reject(service.interceptors.response);
      }
    } else {
      return response.data;
    }
  }, (error) => {
    console.log(error)
    // 对响应错误做点什么
    let { message, response } = error;
    if (message.indexOf('timeout') != -1) {
      message = '系统接口请求超时';
    } else if (message == 'Network Error') {
      message = '后端接口连接异常';
    } else if (message.includes("Request failed with status code")) {
      message = "系统接口" + message.substr(message.length - 3) + "异常";
    } else {
      if (response.data) {
        message = response.statusText;
      } else {
        message = '接口路径找不到';
      }
    }
    ElMessage({
      type: 'error',
      message: message,
      showClose: true,
    })
    return Promise.reject(error);
  }
);

// 导出 axios 实例
export default service;
