import { useEffect, useState } from 'react'
import { api } from '../lib/api.js'
import { AuthContext } from './context.js'

export default function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  useEffect(() => {
    const controller = new AbortController()
    const interceptor = api.interceptors.response.use(
      (response) => response,
      (error) => {
        if (
          error.response?.status === 401 &&
          !error.config?.url?.endsWith('/auth/login')
        )
          setUser(null)
        return Promise.reject(error)
      },
    )
    api
      .get('/auth/me', { signal: controller.signal })
      .then(({ data }) => {
        if (!controller.signal.aborted) setUser(data.data)
      })
      .catch(() => {})
      .finally(() => {
        if (!controller.signal.aborted) setLoading(false)
      })
    return () => {
      controller.abort()
      api.interceptors.response.eject(interceptor)
    }
  }, [])
  const login = async (email, password) => {
    await api.post('/auth/login', { email, password })
    // Confirm that this browser actually retained/sends the session cookie.
    const { data } = await api.get('/auth/me')
    setUser(data.data)
    return data.data
  }
  const logout = async () => {
    try {
      await api.post('/auth/logout')
    } catch (error) {
      if (error.response?.status !== 401) throw error
    }
    setUser(null)
  }
  return (
    <AuthContext.Provider value={{ user, loading, login, logout, setUser }}>
      {children}
    </AuthContext.Provider>
  )
}
