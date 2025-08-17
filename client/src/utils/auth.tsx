import axios, { AxiosError } from "axios"
import type { RegisterResponse } from "../types/types"
import React from "react"

export class LoginFailedError extends Error { }
export class RegisterFailedError extends Error { }
export class UserExistsError extends Error { }
export class UserNotLoggedInError extends Error { }
export class InvalidRegisterError extends Error { }

export interface User {
  username: string
  userId: string
  accessToken: string
}
export interface AuthContext {
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  logout: (username?: string, userId?: string) => Promise<void>
  user: User | null
}

const AuthContext = React.createContext<AuthContext | null>(null)

const TOKEN_KEY = "tanstack.auth.accessToken"
const USERNAME_KEY = "tanstack.auth.username"
const USERID_KEY = "tanstack.auth.userId"

function getStoredUser(): User | null {
  const accessToken = localStorage.getItem(TOKEN_KEY)
  const username = localStorage.getItem(USERNAME_KEY)
  const userId = String(localStorage.getItem(USERID_KEY) ?? "")
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
  const [user, setUser] = React.useState<User | null>(() => getStoredUser())
  const isAuthenticated = !!user

  const logout = React.useCallback(async (username?: string, userId?: string) => {
    await logoutUser(username, userId)
    setUser(null)
    setStoredUser(null)
  }, [])

  const login = React.useCallback(async (username: string, password: string) => {
    const user = await postLogin(username, password)
    setStoredUser(user)
    setUser(user)
  }, [])
  const authContextValue = React.useMemo(() => ({
    isAuthenticated,
    user,
    login,
    logout,
  }), [isAuthenticated, user, login, logout])

  React.useEffect(() => {
    setUser(getStoredUser())
  }, [])

  return (
    <AuthContext.Provider value={authContextValue}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = React.use(AuthContext)
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
  try {
    const response = await axios.post<User>(`/user/login`, loginData)
    return response.data
  } catch (err) {
    if (err instanceof AxiosError) {
      const errStatus = err.response?.status
      if (errStatus === 401 || errStatus === 404 || errStatus === 500) {
        throw new LoginFailedError(`login failed for user ${username}`)
      }
    }
    throw err
  }
}

export const postRegister = async (email: string, password: string, username: string): Promise<RegisterResponse> => {
  const registerData = {
    "email": email,
    "password": password,
    "username": username,
  }
  try {
    const response = await axios.post<RegisterResponse>(`/user/register`, registerData)
    return response.data
  } catch (err) {
    if (err instanceof AxiosError) {
      const errStatus = err.response?.status
      if (errStatus === 400) {
        throw new InvalidRegisterError("unacceptable password")
      } 
      if (errStatus === 409) {
        throw new UserExistsError(`an account named ${username} already exists.`)
      }
      if (errStatus === 404 || errStatus === 500) {
        throw new RegisterFailedError(`register failed`)
      }
    }
    throw err
  }
}

export const logoutUser = async (username?: string, userId?: string): Promise<void> => {
  if (!username || !userId) {
    throw new Error("Missing username or userId for logout");
  }
  try {
    await axios.post(`/user/logout`, {
      userId: String(userId),
      username: username,
    })
  } catch (err) {
    if (err instanceof AxiosError) {
      const errStatus = err.response?.status
      if (errStatus === 400 || errStatus === 500) {
        throw new Error
      }
    }
    throw new Error
  }
}