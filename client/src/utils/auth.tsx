
import axios from "axios"
import type { LoginResponse } from "../types/types"

const url = import.meta.env.VITE_BACKEND_URL

export class LoginFailedError extends Error { }

export type Auth = {
  login: (username: string, password: string) =>void
  logout: () => void
  status: 'loggedOut' | 'loggedIn'
  username?: string
}
export const auth: Auth = {
  status: 'loggedOut',
  username: undefined,
  login: async (username: string, password: string) => {
    const user = await postLogin(username, password)
    localStorage.setItem('accessToken', user.accessToken)
    localStorage.setItem('userId', user.userId.toString())
    auth.status = 'loggedIn'
    auth.username = username
  },
  logout: () => {
    localStorage.removeItem('accessToken')
    localStorage.removeItem('userId')
    auth.status = 'loggedOut'
    auth.username = undefined
  },
}



const postLogin = async (username: string, password: string ): Promise<LoginResponse> =>  {
  const loginData = {
    "username": username,
    "password": password,
  }
  await new Promise((r) => setTimeout(r, 500))
  const response = await axios
    .post<LoginResponse>(`${url}/user/login`, loginData)
    .then((r) => r.data)
    .catch((err) => {
      if (err.status === 401 || err.status === 404 || err.status === 500) {
        throw new LoginFailedError(`login failed for user ${username}`)
      }
      throw err
    })
    
  return response

}