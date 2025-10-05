'use client'

import Navigation from "@/components/Navigation"
import { AuthContext } from "@/context/AuthContext"
import { MapPin, Hospital, Home, Flame, Phone, Pill } from "lucide-react"
import { useContext } from "react"

const ResourcesPage = () => {
  const { user, logout } = useContext(AuthContext)

  const resourceTypes = [
    { name: "Hospitals", icon: Hospital, count: 45, color: "red" },
    { name: "Shelters", icon: Home, count: 23, color: "blue" },
    { name: "Fire Stations", icon: Flame, count: 18, color: "orange" },
    { name: "Police Stations", icon: Phone, count: 32, color: "purple" },
    { name: "Pharmacies", icon: Pill, count: 38, color: "green" },
  ]

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Navigation user={user} onLogout={logout} />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
              <MapPin className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-4xl font-bold text-gray-900 dark:text-white">Emergency Resources</h1>
              <p className="text-gray-600 dark:text-gray-400">Find nearby hospitals, shelters, and emergency services</p>
            </div>
          </div>
        </div>

        {/* Resource Types */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
          {resourceTypes.map((resource) => {
            const colorClasses = {
              red: "bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400",
              blue: "bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400",
              orange: "bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400",
              purple: "bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400",
              green: "bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400",
            }

            return (
              <div
                key={resource.name}
                className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6 hover:shadow-xl transition-shadow cursor-pointer"
              >
                <div className="flex items-center justify-between mb-4">
                  <div className={`w-12 h-12 rounded-lg flex items-center justify-center ${colorClasses[resource.color as keyof typeof colorClasses]}`}>
                    <resource.icon className="w-6 h-6" />
                  </div>
                  <span className="text-2xl font-bold text-gray-900 dark:text-white">{resource.count}</span>
                </div>
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{resource.name}</h3>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Click to view nearby locations</p>
              </div>
            )
          })}
        </div>

        {/* Map Placeholder */}
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <h2 className="text-xl font-bold text-gray-900 dark:text-white">Resource Map</h2>
          </div>
          
          <div className="aspect-video bg-gray-100 dark:bg-gray-900 flex items-center justify-center">
            <div className="text-center text-gray-500 dark:text-gray-400">
              <MapPin className="w-16 h-16 mx-auto mb-4 text-gray-400 dark:text-gray-600" />
              <p className="font-medium">Interactive Map</p>
              <p className="text-sm mt-1">Integrate with Google Maps or OpenStreetMap</p>
              <p className="text-xs mt-2 text-gray-400">Show resources based on user location</p>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export default ResourcesPage
