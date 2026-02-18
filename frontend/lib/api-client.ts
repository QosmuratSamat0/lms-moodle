import { ApiError } from "@/types/common";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

// Token management
const TOKEN_KEY = "lms_access_token";
const REFRESH_TOKEN_KEY = "lms_refresh_token";
const REMEMBER_ME_KEY = "lms_remember_me";

// Get the appropriate storage based on remember me setting
const getStorage = (): Storage | null => {
  if (typeof window === "undefined") return null;
  const rememberMe = localStorage.getItem(REMEMBER_ME_KEY) === "true";
  return rememberMe ? localStorage : sessionStorage;
};

export const tokenStorage = {
  setRememberMe: (remember: boolean): void => {
    if (typeof window === "undefined") return;
    localStorage.setItem(REMEMBER_ME_KEY, String(remember));
  },

  getRememberMe: (): boolean => {
    if (typeof window === "undefined") return false;
    return localStorage.getItem(REMEMBER_ME_KEY) === "true";
  },

  getAccessToken: (): string | null => {
    if (typeof window === "undefined") return null;
    // Check both storages for backwards compatibility
    return sessionStorage.getItem(TOKEN_KEY) || localStorage.getItem(TOKEN_KEY);
  },

  setAccessToken: (token: string): void => {
    const storage = getStorage();
    if (!storage) return;
    storage.setItem(TOKEN_KEY, token);
  },

  getRefreshToken: (): string | null => {
    if (typeof window === "undefined") return null;
    // Check both storages for backwards compatibility
    return sessionStorage.getItem(REFRESH_TOKEN_KEY) || localStorage.getItem(REFRESH_TOKEN_KEY);
  },

  setRefreshToken: (token: string): void => {
    const storage = getStorage();
    if (!storage) return;
    storage.setItem(REFRESH_TOKEN_KEY, token);
  },

  clearTokens: (): void => {
    if (typeof window === "undefined") return;
    // Clear from both storages
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    sessionStorage.removeItem(TOKEN_KEY);
    sessionStorage.removeItem(REFRESH_TOKEN_KEY);
  },

  setTokens: (accessToken: string, refreshToken: string, rememberMe?: boolean): void => {
    if (rememberMe !== undefined) {
      tokenStorage.setRememberMe(rememberMe);
    }
    tokenStorage.setAccessToken(accessToken);
    tokenStorage.setRefreshToken(refreshToken);
  },
};

// API Error class
export class APIError extends Error {
  constructor(
    public status: number,
    public data: ApiError,
  ) {
    super(data.message || "An error occurred");
    this.name = "APIError";
  }
}

// Request options type
interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>;
}

// Build URL with query params
function buildUrl(
  endpoint: string,
  params?: Record<string, string | number | boolean | undefined>,
): string {
  const url = new URL(`${API_BASE_URL}${endpoint}`);

  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        url.searchParams.append(key, String(value));
      }
    });
  }

  return url.toString();
}

// Backend response wrapper type
interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
  };
}

// Main fetch wrapper
async function fetchApi<T>(
  endpoint: string,
  options: RequestOptions = {},
): Promise<T> {
  const { params, ...fetchOptions } = options;

  const url = buildUrl(endpoint, params);

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  // Add auth token if available
  const token = tokenStorage.getAccessToken();
  if (token) {
    (headers as Record<string, string>)["Authorization"] = `Bearer ${token}`;
  }

  let response: Response;
  try {
    response = await fetch(url, {
      ...fetchOptions,
      headers,
    });
  } catch (error) {
    // Network error - backend unreachable
    console.error("Network error:", error);
    throw new APIError(0, {
      message: "Unable to connect to server. Please check your connection.",
      code: "NETWORK_ERROR",
    });
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return {} as T;
  }

  // Handle non-JSON responses
  const contentType = response.headers.get("content-type");
  if (!contentType?.includes("application/json")) {
    if (!response.ok) {
      throw new APIError(response.status, {
        message: `HTTP error ${response.status}`,
      });
    }
    return {} as T;
  }

  let json: ApiResponse<T>;
  try {
    json = (await response.json()) as ApiResponse<T>;
  } catch {
    throw new APIError(response.status, {
      message: "Invalid response from server",
      code: "PARSE_ERROR",
    });
  }

  if (!response.ok) {
    throw new APIError(response.status, {
      message: json.error?.message || `HTTP error ${response.status}`,
      code: json.error?.code,
    });
  }

  // Backend wraps data in {success: true, data: {...}}
  // Return the unwrapped data
  if (json === null || json === undefined) {
    return {} as T;
  }
  return (json.data !== undefined ? json.data : json) as T;
}

// HTTP method helpers
export const api = {
  get: <T>(
    endpoint: string,
    params?: Record<string, string | number | boolean | undefined>,
  ) => fetchApi<T>(endpoint, { method: "GET", params }),

  post: <T>(endpoint: string, body?: unknown) =>
    fetchApi<T>(endpoint, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    }),

  put: <T>(endpoint: string, body?: unknown) =>
    fetchApi<T>(endpoint, {
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    }),

  patch: <T>(endpoint: string, body?: unknown) =>
    fetchApi<T>(endpoint, {
      method: "PATCH",
      body: body ? JSON.stringify(body) : undefined,
    }),

  delete: <T>(endpoint: string) => fetchApi<T>(endpoint, { method: "DELETE" }),

  // File upload method using FormData
  upload: async <T>(endpoint: string, formData: FormData): Promise<T> => {
    const url = buildUrl(endpoint);

    const headers: HeadersInit = {};
    // Don't set Content-Type - browser will set it with boundary for multipart/form-data

    const token = tokenStorage.getAccessToken();
    if (token) {
      (headers as Record<string, string>)["Authorization"] = `Bearer ${token}`;
    }

    let response: Response;
    try {
      response = await fetch(url, {
        method: "POST",
        headers,
        body: formData,
      });
    } catch (error) {
      console.error("Network error:", error);
      throw new APIError(0, {
        message: "Unable to connect to server. Please check your connection.",
        code: "NETWORK_ERROR",
      });
    }

    const contentType = response.headers.get("content-type");
    if (!contentType?.includes("application/json")) {
      if (!response.ok) {
        throw new APIError(response.status, {
          message: `HTTP error ${response.status}`,
        });
      }
      return {} as T;
    }

    let json: ApiResponse<T>;
    try {
      json = (await response.json()) as ApiResponse<T>;
    } catch {
      throw new APIError(response.status, {
        message: "Invalid response from server",
        code: "PARSE_ERROR",
      });
    }

    if (!response.ok) {
      throw new APIError(response.status, {
        message: json.error?.message || `HTTP error ${response.status}`,
        code: json.error?.code,
      });
    }

    return (json.data !== undefined ? json.data : json) as T;
  },
};

export default api;
