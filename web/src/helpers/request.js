import axios from "axios";

const request = axios.create({
  timeout: 10 * 1000
});

request.interceptors.response.use(null, err => {
  const { response } = err;
  if (err.code === "ECONNABORTED") {
    err.message = "请求超时，请稍候再试";
  } else if (response) {
    if (response.data && response.data.message) {
      err.message = response.data.message;
      err.code = response.data.code;
    } else {
      err.message = `unknown error[${response.statusCode || -1}]`;
    }
  }
  return Promise.reject(err);
});

export default request;
