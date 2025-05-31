import axios, {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
} from "axios";

const ENV_BASE_URL =
  import.meta.env.MODE != "production"
    ? import.meta?.env?.VITE_LOCAL_SERVER || "http://localhost:8080"
    : import.meta?.env?.VITE_PROD_SERVER;

const BASE_URL = `${ENV_BASE_URL}${import.meta.env.VITE_SERVER_API_SUFFIX}/${
  import.meta.env.VITE_SERVER_API_VERSION
}`;

const DEFAULT_HEADERS = {
  "Content-Type": "application/json",
};

export interface APIResponseData<T> {
  status_code: number;
  success: boolean;
  error_message?: string | null;
  data?: T;
}

export interface CustomAxiosError<T = any> extends AxiosError {
  response?: AxiosError["response"] & { data: APIResponseData<T> };
}

// Custom error interface
export interface ApiError {
  status: number;
  message: string;
  data?: APIResponseData<null>;
}

type ApiMethod = "get" | "post" | "put" | "delete" | "patch";

// Auth event system for notifying context about token invalidation
class AuthEventManager {
  private static instance: AuthEventManager;
  private listeners: (() => void)[] = [];

  static getInstance(): AuthEventManager {
    if (!AuthEventManager.instance) {
      AuthEventManager.instance = new AuthEventManager();
    }
    return AuthEventManager.instance;
  }

  addListener(callback: () => void): void {
    this.listeners.push(callback);
  }

  removeListener(callback: () => void): void {
    this.listeners = this.listeners.filter(listener => listener !== callback);
  }

  notifyTokenInvalidation(): void {
    this.listeners.forEach(callback => callback());
  }
}

export const authEventManager = AuthEventManager.getInstance();

class ApiHandler {
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = axios.create({
      baseURL: BASE_URL,
      headers: DEFAULT_HEADERS,
      timeout: 15000, // 15 seconds timeout
      // withCredentials: true,
    });

    this.setupInterceptors();
  }

  private setupInterceptors(): void {
    // Request interceptor to add auth token
    this.axiosInstance.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem("access_token");

        config.headers.Authorization = `Bearer ${token}`;

        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    this.axiosInstance.interceptors.response.use(
      (response) => response,
      (error) => {
        return Promise.reject(this.handleError(error));
      }
    );
  }

  private handleError(error: CustomAxiosError): ApiError {
    if (error.response) {
      console.log("error response", error.response);
      // Server responded with an error status
      const status = error.response.status;

      switch (status) {
        case 401:
          // Handle unauthorized - clear tokens and notify auth context
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
          localStorage.removeItem("user");
          // Notify AuthContext about token invalidation
          authEventManager.notifyTokenInvalidation();
          break;
        case 404:
          // Handle not found
          console.error("Resource not found");
          break;
        case 429:
          // Handle rate limiting
          console.error("Rate limit exceeded");
          break;
        case 500:
          // Handle server error
          console.error("Internal server error");
          break;
      }

      return {
        status,
        message: error.response.data?.error_message || "An error occurred",
        data: error.response.data,
      };
    }
    if (error.request) {
      // Network error - request made but no response received
      return {
        status: 0,
        message: navigator.onLine
          ? "Server is not responding"
          : "No internet connection",
      };
    }

    // Something happened in setting up the request
    return {
      status: 0,
      message: error.message || "An unexpected error occurred",
    };
  }

  private async request<T>(
    method: ApiMethod,
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<APIResponseData<T>> {
    try {
      const response: AxiosResponse<APIResponseData<T>> =
        await this.axiosInstance.request({
          method,
          url,
          data,
          ...config,
        });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        throw this.handleError(error);
      }
      throw error;
    }
  }

  public async get<T>(
    url: string,
    config?: AxiosRequestConfig
  ): Promise<APIResponseData<T>> {
    return this.request<T>("get", url, undefined, config);
  }

  public async post<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<APIResponseData<T>> {
    return this.request<T>("post", url, data, config);
  }

  public async put<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<APIResponseData<T>> {
    return this.request<T>("put", url, data, config);
  }

  public async delete<T>(
    url: string,
    config?: AxiosRequestConfig
  ): Promise<APIResponseData<T>> {
    return this.request<T>("delete", url, undefined, config);
  }

  public async patch<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<APIResponseData<T>> {
    return this.request<T>("patch", url, data, config);
  }
}

const apiHandler = new ApiHandler();
export default apiHandler;
