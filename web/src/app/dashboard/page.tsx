'use client'

import { AuthContext } from "@/context/AuthContext"
import { Role } from "@/types/types"
import { AlertTriangle, MapPin, Users, TrendingUp, ChevronRight } from "lucide-react"
import Link from "next/link"
import { useContext } from "react"

const DashboardPage = () => {
  const { user, loading } = useContext(AuthContext)

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">
          {user ? `Welcome back, ${user.username}!` : 'Disaster Relief Dashboard'}
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          {user 
            ? 'Track disasters, view resources, and coordinate relief efforts'
            : 'View active disasters and available resources in your area'}
        </p>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Active Disasters</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">24</p>
            </div>
            <div className="w-12 h-12 bg-red-100 dark:bg-red-900/30 rounded-lg flex items-center justify-center">
              <AlertTriangle className="w-6 h-6 text-red-600 dark:text-red-400" />
            </div>
          </div>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-4">+3 in the last 24 hours</p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Available Resources</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">156</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
              <MapPin className="w-6 h-6 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-4">Hospitals, shelters, fire stations</p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Active Volunteers</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">1,247</p>
            </div>
            <div className="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
              <Users className="w-6 h-6 text-green-600 dark:text-green-400" />
            </div>
          </div>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-4">+89 joined this week</p>
        </div>
      </div>

      {/* Quick Actions */}
      {user && (
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">Quick Actions</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {(user.role === Role.ADMIN || user.role === Role.VOLUNTEER) && (
              <Link
                href="/report-disaster"
                className="group bg-gradient-to-br from-red-500 to-orange-600 hover:from-red-600 hover:to-orange-700 text-white rounded-xl p-6 shadow-lg hover:shadow-xl transition-all"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-semibold mb-1">Report Disaster</h3>
                    <p className="text-sm text-red-100">Submit a new disaster report</p>
                  </div>
                  <ChevronRight className="w-6 h-6 group-hover:translate-x-1 transition-transform" />
                </div>
              </Link>
            )}

            <Link
              href="/resources"
              className="group bg-gradient-to-br from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white rounded-xl p-6 shadow-lg hover:shadow-xl transition-all"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold mb-1">Find Resources</h3>
                  <p className="text-sm text-blue-100">Locate nearby emergency resources</p>
                </div>
                <ChevronRight className="w-6 h-6 group-hover:translate-x-1 transition-transform" />
              </div>
            </Link>

            {user.role === Role.ADMIN && (
              <Link
                href="/admin"
                className="group bg-gradient-to-br from-purple-500 to-pink-600 hover:from-purple-600 hover:to-pink-700 text-white rounded-xl p-6 shadow-lg hover:shadow-xl transition-all"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-semibold mb-1">Admin Panel</h3>
                    <p className="text-sm text-purple-100">Review and manage reports</p>
                  </div>
                  <ChevronRight className="w-6 h-6 group-hover:translate-x-1 transition-transform" />
                </div>
              </Link>
            )}
          </div>
        </div>
      )}

      {/* Call to Action for Non-Logged Users */}
      {!user && !loading && (
        <div className="bg-gradient-to-br from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 rounded-2xl border-2 border-dashed border-blue-300 dark:border-blue-700 p-8 text-center mb-8">
          <TrendingUp className="w-16 h-16 text-blue-600 dark:text-blue-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
            Join Our Relief Network
          </h2>
          <p className="text-gray-600 dark:text-gray-400 mb-6 max-w-2xl mx-auto">
            Sign up to report disasters, access exclusive resources, and help coordinate relief efforts in your community.
          </p>
          <div className="flex items-center justify-center gap-4">
            <Link
              href="/auth"
              className="px-6 py-3 bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white font-semibold rounded-lg shadow-lg shadow-blue-500/30 hover:shadow-blue-500/50 transition-all"
            >
              Sign Up Now
            </Link>
            <Link
              href="/auth"
              className="px-6 py-3 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 font-semibold rounded-lg border-2 border-gray-300 dark:border-gray-600 hover:border-blue-500 dark:hover:border-blue-400 transition-all"
            >
              Sign In
            </Link>
          </div>
        </div>
      )}

      {/* Recent Disasters */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Recent Disasters</h2>
          <Link
            href="/disasters"
            className="text-blue-600 dark:text-blue-400 hover:underline text-sm font-medium"
          >
            View All →
          </Link>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          {/* Placeholder for disaster list */}
          <div className="p-8 text-center text-gray-500 dark:text-gray-400">
            <AlertTriangle className="w-12 h-12 mx-auto mb-3 text-gray-400 dark:text-gray-600" />
            <p>Recent disaster reports will appear here</p>
            <p className="text-sm mt-1">Connect to your backend API to load live data</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DashboardPage