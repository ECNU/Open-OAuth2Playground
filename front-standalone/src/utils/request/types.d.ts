import {
    Method,
    AxiosError,
    AxiosResponse,
    AxiosRequestConfig
} from "axios";

export type RequestMethods = Extract<Method,
    "get" | "post" | "put" | "delete" | "patch" | "option" | "head">;

export interface TSHttpError extends AxiosError{
    isCancelRequest?: boolean;
}

export interface TSHttpResponse extends AxiosResponse{
    config: TSHttpRequestConfig;
}

export interface TSHttpRequestConfig extends AxiosRequestConfig{
    beforeRequestCallback?: (request: TSHttpRequestConfig) => void;
    beforeResponseCallback?: (response: TSHttpResponse) => void;
}