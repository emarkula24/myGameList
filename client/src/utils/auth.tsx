
import axios from "axios"
import type { LoginResponse, RegisterResponse } from "../types/types"

const url = import.meta.env.VITE_BACKEND_URL

export class LoginFailedError extends Error { }
export class RegisterFailedError extends Error { }

export type Auth = {
  login: (username: string, password: string) => Promise<void>
  logout: () => void
  status: 'loggedOut' | 'loggedIn'
  username?: string
  userId?: string 
}
export const auth: Auth = {
  status: 'loggedOut',
  username: undefined,
  login: async (username: string, password: string) => {
    const user = await postLogin(username, password)
    localStorage.setItem('accessToken', user.accessToken)
    auth.userId = user.userId.toString()
    auth.status = 'loggedIn'
    auth.username = username
  },
  logout: () => {
    localStorage.removeItem('accessToken')
    auth.userId = undefined
    auth.status = 'loggedOut'
    auth.username = undefined
  },
}



const postLogin = async (username: string, password: string): Promise<LoginResponse> => {
  const loginData = {
    "username": username,
    "password": password,
  }
  await new Promise((r) => setTimeout(r, 500))
  const response = await axios
    .post<LoginResponse>(`${url}/user/login`, loginData)
    .then((r) => r.data)
    .catch((err) => {
      const errStatus = err.response?.status
      if (errStatus === 401 || errStatus === 404 || errStatus === 500) {
        throw new LoginFailedError(`login failed for user ${username}`)
      }
      throw err
    })

  return response

}

export const postRegister = async (email: string, password: string, username: string): Promise<RegisterResponse> => {
  const registerData = {
    "email": email,
    "password": password,
    "username": username,
  }
  await new Promise((r) => setTimeout(r, 500))
  const response = await axios
    .post<RegisterResponse>(`${url}/user/register`, registerData)
    .then((r) => r.data)
    .catch((err) => {
      const errStatus = err.response?.status
      if (errStatus === 401 || errStatus === 404 || errStatus === 500) {
        throw new RegisterFailedError(`login failed for user ${username}`)
      }
      throw err
    })
  return response

}