'use client'

import { getCurrentUser, logout as apiLogout } from "@/api/auth"
import { User } from "@/types/types"
import { useRouter } from "next/navigation"
import { createContext, ReactNode, useContext, useEffect, useState } from "react"

interface AuthContextType {
  user: User | null
  loading: boolean
  error: string | null
  setUser: (user: User | null) => void
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

export const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  error: null,
  setUser: () => {},
  logout: async () => {},
  refreshUser: async () => {}
})

// Custom hook to use the auth context
export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()

  // Fetch current user on mount
  useEffect(() => {
    const fetchUser = async () => {
      try {
        setLoading(true)
        const result = await getCurrentUser()
        setUser(result.data)
      } catch (err) {
        console.error('Failed to fetch user:', err)
        setUser(null)
      } finally {
        setLoading(false)
      }
    }

    fetchUser()
  }, [])

  const refreshUser = async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await getCurrentUser()
      setUser(result.data)
    } catch (err) {
      console.error('Failed to refresh user:', err)
      setError('Failed to refresh user')
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  const logout = async () => {
    try {
      setLoading(true)
      setError(null)
      await apiLogout()
      setUser(null)
      router.push('/auth')
    } catch (err) {
      console.error('Logout failed:', err)
      setError('Logout failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthContext.Provider value={{ user, loading, error, setUser, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  )
}