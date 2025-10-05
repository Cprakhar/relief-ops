"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/context/AuthContext"
import { getCurrentUser } from "@/api/auth"
import { Loader2 } from "lucide-react"

export default function OAuthSuccessPage() {
  const router = useRouter()
  const { setUser } = useAuth()

  useEffect(() => {
    const handleOAuthSuccess = async () => {
      try {
        // Token is already set in cookie by backend
        // Fetch the current user data
        const userData = await getCurrentUser()
        setUser(userData)
        
        // Redirect to dashboard
        router.push("/dashboard")
      } catch (error) {
        console.error("OAuth success error:", error)
        router.push("/auth?error=oauth_failed")
      }
    }

    handleOAuthSuccess()
  }, [router, setUser])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
      <div className="text-center">
        <Loader2 className="w-12 h-12 animate-spin text-blue-600 mx-auto mb-4" />
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
          Completing Sign In
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          Please wait while we set up your account...
        </p>
      </div>
    </div>
  )
}
