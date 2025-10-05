'use client'

import Navigation from "@/components/Navigation"
import { AuthContext } from "@/context/AuthContext"
import { Role } from "@/types/types"
import { AlertTriangle, MapPin, Camera, Tag, Activity } from "lucide-react"
import { useRouter } from "next/navigation"
import { useContext, useEffect, useState } from "react"

const ReportDisasterPage = () => {
  const { user, loading, logout } = useContext(AuthContext)
  const router = useRouter()
  const [severity, setSeverity] = useState(5)

  useEffect(() => {
    if (!loading && !user) {
      router.push("/auth")
    }
  }, [user, loading, router])

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-red-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600 dark:text-gray-400">Loading...</p>
        </div>
      </div>
    )
  }

  if (!user) {
    return null
  }

  // Check if user is allowed to report
  if (user.role !== Role.ADMIN && user.role !== Role.VOLUNTEER) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <Navigation user={user} onLogout={logout} />
        <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-16 text-center">
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-8">
            <AlertTriangle className="w-16 h-16 text-orange-500 mx-auto mb-4" />
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">Access Restricted</h1>
            <p className="text-gray-600 dark:text-gray-400 mb-6">
              Only volunteers and administrators can report disasters.
            </p>
            <button
              onClick={() => router.push("/dashboard")}
              className="px-6 py-3 bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white font-semibold rounded-lg transition-all"
            >
              Back to Dashboard
            </button>
          </div>
        </div>
      </div>
    )
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // TODO: Implement disaster report submission
    console.log("Disaster report submitted")
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Navigation user={user} onLogout={logout} />
      
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-12 h-12 bg-gradient-to-br from-red-500 to-orange-600 rounded-xl flex items-center justify-center shadow-lg">
              <AlertTriangle className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-4xl font-bold text-gray-900 dark:text-white">Report a Disaster</h1>
              <p className="text-gray-600 dark:text-gray-400">Help coordinate relief efforts by reporting disasters</p>
            </div>
          </div>
        </div>

        {/* Info Banner */}
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-4 mb-8">
          <div className="flex items-start gap-3">
            <div className="w-5 h-5 bg-blue-500 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
              <span className="text-white text-xs font-bold">i</span>
            </div>
            <div className="flex-1">
              <h3 className="text-sm font-semibold text-blue-900 dark:text-blue-200">Important Information</h3>
              <p className="text-sm text-blue-800 dark:text-blue-300 mt-1">
                Your report will be reviewed by our admin team before being made public. Please provide accurate and detailed information.
              </p>
            </div>
          </div>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-8">
          <div className="space-y-6">
            {/* Title */}
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Disaster Title *
              </label>
              <input
                id="title"
                type="text"
                name="title"
                required
                placeholder="e.g., Earthquake in San Francisco"
                className="w-full px-4 py-3 bg-gray-50 dark:bg-gray-700/50 border border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all outline-none text-gray-900 dark:text-white placeholder-gray-400"
              />
            </div>

            {/* Description */}
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Description *
              </label>
              <textarea
                id="description"
                name="description"
                required
                rows={5}
                placeholder="Provide detailed information about the disaster..."
                className="w-full px-4 py-3 bg-gray-50 dark:bg-gray-700/50 border border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all outline-none text-gray-900 dark:text-white placeholder-gray-400 resize-none"
              />
            </div>

            {/* Location */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label htmlFor="latitude" className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  <MapPin className="w-4 h-4" />
                  Latitude *
                </label>
                <input
                  id="latitude"
                  type="number"
                  name="latitude"
                  step="any"
                  required
                  placeholder="37.7749"
                  className="w-full px-4 py-3 bg-gray-50 dark:bg-gray-700/50 border border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all outline-none text-gray-900 dark:text-white placeholder-gray-400"
                />
              </div>
              
              <div>
                <label htmlFor="longitude" className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  <MapPin className="w-4 h-4" />
                  Longitude *
                </label>
                <input
                  id="longitude"
                  type="number"
                  name="longitude"
                  step="any"
                  required
                  placeholder="-122.4194"
                  className="w-full px-4 py-3 bg-gray-50 dark:bg-gray-700/50 border border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all outline-none text-gray-900 dark:text-white placeholder-gray-400"
                />
              </div>
            </div>

            {/* Severity */}
            <div>
              <label htmlFor="severity" className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                <Activity className="w-4 h-4" />
                Severity: {severity}/10
              </label>
              <input
                id="severity"
                type="range"
                name="severity"
                min="1"
                max="10"
                value={severity}
                onChange={(e) => setSeverity(parseInt(e.target.value))}
                className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-red-500"
              />
              <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-1">
                <span>Minor</span>
                <span>Moderate</span>
                <span>Severe</span>
                <span>Critical</span>
              </div>
            </div>

            {/* Tags */}
            <div>
              <label htmlFor="tags" className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                <Tag className="w-4 h-4" />
                Tags (comma-separated)
              </label>
              <input
                id="tags"
                type="text"
                name="tags"
                placeholder="e.g., earthquake, urban, urgent"
                className="w-full px-4 py-3 bg-gray-50 dark:bg-gray-700/50 border border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all outline-none text-gray-900 dark:text-white placeholder-gray-400"
              />
            </div>

            {/* Photo Upload */}
            <div>
              <label htmlFor="photo" className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                <Camera className="w-4 h-4" />
                Upload Photo (optional)
              </label>
              <div className="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-xl p-6 text-center hover:border-red-500 dark:hover:border-red-400 transition-colors">
                <input
                  id="photo"
                  type="file"
                  name="photo"
                  accept="image/*"
                  className="hidden"
                />
                <label htmlFor="photo" className="cursor-pointer">
                  <Camera className="w-12 h-12 text-gray-400 dark:text-gray-600 mx-auto mb-2" />
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Click to upload or drag and drop
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                    PNG, JPG up to 10MB
                  </p>
                </label>
              </div>
            </div>

            {/* Submit Button */}
            <div className="flex items-center gap-4 pt-4">
              <button
                type="submit"
                className="flex-1 py-3 px-6 bg-gradient-to-r from-red-500 to-orange-600 hover:from-red-600 hover:to-orange-700 text-white font-semibold rounded-xl shadow-lg shadow-red-500/30 hover:shadow-red-500/50 transition-all"
              >
                Submit Report
              </button>
              <button
                type="button"
                onClick={() => router.push("/dashboard")}
                className="px-6 py-3 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 font-semibold rounded-xl hover:bg-gray-300 dark:hover:bg-gray-600 transition-all"
              >
                Cancel
              </button>
            </div>
          </div>
        </form>
      </main>
    </div>
  )
}

export default ReportDisasterPage
