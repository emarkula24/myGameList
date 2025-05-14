import axios from 'axios';


// Response interceptor for handling 401 errors
axios.interceptors.response.use(
  (response) => response, 
  async (error) => {
    if (error.response?.status === 401) {
      try {
        const res = await axios.post('/refresh')
        const { accessToken } = res.data
        
        // Store new token and retry failed request
        localStorage.setItem('token', accessToken)
        error.config.headers.Authorization = `Bearer ${accessToken}`
        return axios(error.config);  // Retry the original request
      } catch (refreshError) {
        // Logout or redirect if refresh fails
        console.error('Refresh token expired or invalid')
        return Promise.reject(refreshError);
        
      }
    }
    return Promise.reject(error);
  }
);

axios.interceptors.request.use((config) => {
  if (config.headers?.requiresAuth) {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    delete config.headers['requiresAuth'];
  }
  return config;
}, (error) => {
  return Promise.reject(error);
})