'use client'

import Navigation from "@/components/Navigation"
import { AuthContext } from "@/context/AuthContext"
import { Role } from "@/types/types"
import { Shield, ShieldAlert } from "lucide-react"
import { useRouter } from "next/navigation"
import { useContext, useEffect } from "react"

const AdminPage = () => {
  const { user, loading, logout } = useContext(AuthContext)
  const router = useRouter()

  useEffect(() => {
    if (!loading && (!user || user.role !== Role.ADMIN)) {
      router.push("/dashboard")
    }
  }, [user, loading, router])

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600 dark:text-gray-400">Loading admin panel...</p>
        </div>
      </div>
    )
  }

  if (!user || user.role !== Role.ADMIN) {
    return null
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Navigation user={user} onLogout={logout} />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-12 h-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl flex items-center justify-center shadow-lg">
              <Shield className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-4xl font-bold text-gray-900 dark:text-white">Admin Dashboard</h1>
              <p className="text-gray-600 dark:text-gray-400">Review and manage disaster reports</p>
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
            <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Pending Review</p>
            <p className="text-3xl font-bold text-orange-600 dark:text-orange-400 mt-2">8</p>
          </div>
          
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
            <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Approved</p>
            <p className="text-3xl font-bold text-green-600 dark:text-green-400 mt-2">142</p>
          </div>
          
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
            <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Rejected</p>
            <p className="text-3xl font-bold text-red-600 dark:text-red-400 mt-2">23</p>
          </div>
          
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
            <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Total Users</p>
            <p className="text-3xl font-bold text-blue-600 dark:text-blue-400 mt-2">1,247</p>
          </div>
        </div>

        {/* Pending Reviews */}
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <h2 className="text-xl font-bold text-gray-900 dark:text-white">Pending Reviews</h2>
          </div>
          
          <div className="p-8 text-center text-gray-500 dark:text-gray-400">
            <ShieldAlert className="w-12 h-12 mx-auto mb-3 text-gray-400 dark:text-gray-600" />
            <p className="font-medium">No pending disaster reports</p>
            <p className="text-sm mt-1">New reports will appear here for review</p>
          </div>
        </div>
      </main>
    </div>
  )
}

export default AdminPage
