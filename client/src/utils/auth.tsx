import axios from "axios"
import type { RegisterResponse } from "../types/types"
import React from "react"
import { sleep } from "./utils"
export class LoginFailedError extends Error { }
export class RegisterFailedError extends Error { }

export type User = {
  username: string
  userId: string
  accessToken: string
}
export interface AuthContext {
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
  user: User | null
}

const AuthContext = React.createContext<AuthContext | null>(null)

const TOKEN_KEY = "tanstack.auth.accessToken"
const USERNAME_KEY = "tanstack.auth.username"
const USERID_KEY = "tanstack.auth.userId"

function getStoredUser(): User | null {
  const accessToken = localStorage.getItem(TOKEN_KEY)
  const username = localStorage.getItem(USERNAME_KEY)
  const userId = localStorage.getItem(USERID_KEY)
  if (accessToken && username && userId) {
    return { username, userId, accessToken }
  }
  return null
}


function setStoredUser(user: User | null) {
  if (user) {
    localStorage.setItem(TOKEN_KEY, user.accessToken)
    localStorage.setItem(USERNAME_KEY, user.username)
    localStorage.setItem(USERID_KEY, user.userId)
  } else {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USERNAME_KEY)
    localStorage.removeItem(USERID_KEY)
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = React.useState<User | null>(getStoredUser())
  const isAuthenticated = !!user

  const logout = React.useCallback(async () => {
    await sleep(250)

    setStoredUser(null)
    setUser(null)
  }, [])

  const login = React.useCallback(async (username: string, password: string) => {
    const user = await postLogin(username, password)
    setStoredUser(user)
    setUser(user)
  }, [])

  React.useEffect(() => {
    setUser(getStoredUser())
  }, [])

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = React.useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}


const postLogin = async (username: string, password: string): Promise<User> => {
  const loginData = {
    "username": username,
    "password": password,
  }
  await new Promise((r) => setTimeout(r, 500))
  const response = await axios
    .post<User>(`/user/login`, loginData)
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
    .post<RegisterResponse>(`/user/register`, registerData)
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