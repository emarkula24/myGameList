import axios from "axios"

const axiosAuthorizationInstance = axios.create({
    withCredentials: true,
})



axiosAuthorizationInstance.interceptors.request.use(request => {
    const accessToken = localStorage.getItem('accessToken');
    if (accessToken) {
    request.headers['Authorization'] = `Bearer ${accessToken}`;
    }
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
        const response = await axios.post('/user/refresh', {}, {withCredentials: true});
        const { accessToken } = response.data;
        // Store the new access and refresh tokens.
        localStorage.setItem('accessToken', accessToken);
        // Update the authorization header with the new access token.
        axiosAuthorizationInstance.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
        return axiosAuthorizationInstance(originalRequest); // Retry the original request with the new access token.
      } catch (refreshError) {
        // Handle refresh token errors by clearing stored tokens and redirecting to the login page.
        console.error('Token refresh failed:', refreshError);
        localStorage.removeItem('accessToken');
        console.log("logged out, no redirect")
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(error); // For all other errors, return the error as is.
  }
);

export default axiosAuthorizationInstance