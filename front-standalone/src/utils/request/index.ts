import axios, { AxiosError, AxiosInstance, AxiosPromise, AxiosRequestConfig, AxiosResponse } from 'axios';

export interface Result<T>{
    data: T;
    code: number;
    msg: string;
}

class HttpService{
    private readonly http!: AxiosInstance;

    constructor() {
        const proto = process.env.VUE_APP_API_PROTO;
        const host = process.env.VUE_APP_API_HOST;
        const port = process.env.VUE_APP_API_PORT?.length>0?":".concat(process.env.VUE_APP_API_PORT):"";
        const route = process.env.VUE_APP_API_ROUTE;
        this.http = axios.create({
            baseURL: proto+"://"+host+port+route,
            timeout: 60000,
        });

        this.addInterceptors(this.http);
    }

    get<T>(url: string, config?: AxiosRequestConfig) {
        return this.handleErrorWrapper<T>(this.http.get(url, config));
    }

    post<T>(url: string, param: unknown, config?: AxiosRequestConfig) {
        return this.handleErrorWrapper<T>(this.http.post(url, param, config));
    }

    postDownload<T>(url: string, param: unknown) {
        return this.handleErrorWrapper<T>(this.http.post(url, param, {responseType: 'arraybuffer'}));
    }

    put<T>(url: string, param: unknown, config?: AxiosRequestConfig) {
        return this.handleErrorWrapper<T>(this.http.put(url, param, config));
    }

    delete<T>(url: string, param: unknown, config?: AxiosRequestConfig) {
        return this.handleErrorWrapper<T>(this.http.delete(url, {data: param, ...config}));
    }

    private addInterceptors(http: AxiosInstance) {
        http.interceptors.request.use((config) => {
            return config;
        });

        http.interceptors.response.use(
            (response: AxiosResponse) => {
                return response;
            },
            (error) => {
                return Promise.reject(error);
            },
        );
    }

    private async handleErrorWrapper<T>(p: AxiosPromise): Promise<AxiosResponse<any> | {}> {
        return p
            .then((response) => {
                return response.data;
            })
            .catch((error: AxiosError) => {
                return {
                    ...(error.response?.data as object)
                };
            });
    }
}

export const http = new HttpService();