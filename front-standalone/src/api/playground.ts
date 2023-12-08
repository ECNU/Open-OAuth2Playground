import { http } from "/@/utils/request";

type Result = {
    code: string;
    msg: string;
    data: {
        request?: object | any;
        response?:object | any;
        rawjson?: object | any;
        example?: object | any;
    }
}

/** authorization_code */
/** Step 2 */
/** Get access_token with code */
export const fetchACToken = (data) => {
    return http.post<Result>("/oauth2/authorization_code", data);
};

/** Get new access_token with refresh_token */
export const fetchRefreshToken = (data) => {
    return http.post<Result>("/oauth2/refresh_token", data);
}

/** Step 3*/
/** Get data from API */
export const fetchApiData = (data) => {
    return http.post<Result>("/api", data);
}


/** Client Credentials */
export const fetchACTokenByClient = (data) => {
    return http.post<Result>("/oauth2/client_credentials", data);
};

/** Password */
export const fetchACTokenByPassword = (data) => {
    return http.post<Result>("/oauth2/password", data);
};

/** Device Flow */
/** Step 2 */
/** Get access_token with device_code */
export const fetchACTokenByDevice = (data) => {
    return http.post<Result>("/oauth2/device_flow", data);
};

/** PKCE */
/** Step 2 */
/** Get access_token with PKCE */
export const fetchACTokenByPkce = (data) => {
    return http.post<Result>("/oauth2/pkce", data);
};
