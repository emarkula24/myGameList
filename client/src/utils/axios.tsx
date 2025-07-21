import axios from "axios"

const axiosAuthorizationInstance = axios.create({
    baseURL: import.meta.env.VITE_BACKEND_URL,
    withCredentials: true,
})


const username = localStorage.getItem("username")
const userId = localStorage.getItem("userId")

axiosAuthorizationInstance.interceptors.request.use(request => {
    const accessToken = localStorage.getItem('tanstack.auth.accessToken');
    console.log(accessToken)
    if (accessToken) {
    request.headers['Authorization'] = `Bearer ${accessToken}`;
    }
    console.log(request)
    return request
}, error  => {  
    return Promise.reject(error)
})

axiosAuthorizationInstance.interceptors.response.use(
  response => response, // Directly return successful responses.
  async error => {
    const originalRequest = error.config;
    if (error.response.status === 403 && !originalRequest._retry) {
      originalRequest._retry = true; // Mark the request as retried to avoid infinite loops.
      try {
        console.log("trying to refresh accessToken...")
        const response = await axios.post('/user/refresh', {username: username, userId: userId}, {withCredentials: true});
        const { accessToken } = response.data;
        console.log(response)
        // Store the new access and refresh tokens.
        localStorage.setItem('tanstack.auth.accessToken', accessToken);
        // Update the authorization header with the new access token.
        axiosAuthorizationInstance.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
        return axiosAuthorizationInstance(originalRequest); // Retry the original request with the new access token.
      } catch (refreshError) {
        // Handle refresh token errors by clearing stored tokens and redirecting to the login page.
        console.error('Token refresh failed:', refreshError);
        console.log("logged out, no redirect")
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(error); // For all other errors, return the error as is.
  }
);

export default axiosAuthorizationInstance