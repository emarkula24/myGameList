import axios, { AxiosError, type AxiosRequestConfig, type AxiosResponse } from "axios"

export class AuthorizationError extends Error { }
export class RefreshError extends Error { }

const axiosAuthorizationInstance = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL,
  withCredentials: true,
})

const username = localStorage.getItem("tanstack.auth.username")
const userId = localStorage.getItem("tanstack.auth.userId")

axiosAuthorizationInstance.interceptors.request.use(request => {
  const accessToken = localStorage.getItem('tanstack.auth.accessToken')

  if (accessToken) {
    request.headers['Authorization'] = `Bearer ${accessToken}`;
  }
  return request
}, (error: AxiosError) => {
  return Promise.reject(new AuthorizationError("failed to authoritize with current accesstoken" + String(error.message)))
})

interface RefreshReponse {
  accessToken: string;
}
interface CustomAxiosRequestConfig extends AxiosRequestConfig {
  _retry?: boolean;
}

axiosAuthorizationInstance.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => response, // Directly return successful responses.
  async (error: unknown): Promise<AxiosError> => {
    if (error instanceof AxiosError && error.response) {
    const originalRequest = error.config as CustomAxiosRequestConfig

    if (error.response.status === 403 && !originalRequest._retry) {
      originalRequest._retry = true; // Mark the request as retried to avoid infinite loops.

      try {
        console.log("trying to refresh accessToken...")
        const response = await axios.post<RefreshReponse>('/user/refresh', { username: username, userId: userId }, { withCredentials: true })
        const { accessToken } = response.data
        console.log(response)
        // Store the new access and refresh tokens.
        localStorage.setItem('tanstack.auth.accessToken', accessToken)
        // Update the authorization header with the new access token.
        axiosAuthorizationInstance.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
        
        return axiosAuthorizationInstance(originalRequest); // Retry the original request with the new access token.
      } catch (refreshError: unknown) {
        if (refreshError instanceof AxiosError) {
          console.error('Token refresh failed:', refreshError)
          console.log("logged out, no redirect")
          // Handle refresh token errors by clearing stored tokens and redirecting to the login page.
          return Promise.reject(new RefreshError("failed to fetch a new accessToken when refreshing" + String(refreshError.message)))
        }
        throw refreshError
      }
      }

    }
    return Promise.reject(new Error("failed fetch a new accessToken when refreshing" + String(error))) // For all other errors, return the error as is.
  }
)

export default axiosAuthorizationInstance