import {http} from "/@/utils/request";
import {getApiBase} from "/@/utils/x";

type Result = {
  code?: number;
  message?: string;
  data?: object | any;
};

/** 获取当前配置信息 */
export const fetchConfig = () => {
  return http.get<Result>("/config");
};
