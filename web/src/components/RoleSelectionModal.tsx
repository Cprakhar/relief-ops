"use client"

import { Role } from "@/types/types"
import { X, Users, Heart } from "lucide-react"

interface RoleSelectionModalProps {
  isOpen: boolean
  onClose: () => void
  onSelectRole: (role: Role) => void
  provider: "google" | "github"
}

export default function RoleSelectionModal({ isOpen, onClose, onSelectRole, provider }: RoleSelectionModalProps) {
  if (!isOpen) return null

  const providerName = provider === "google" ? "Google" : "GitHub"

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full p-6 relative animate-in fade-in zoom-in duration-200">
        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
        >
          <X className="w-5 h-5" />
        </button>

        {/* Header */}
        <div className="text-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
            Choose Your Role
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            Select how you&apos;d like to contribute with {providerName}
          </p>
        </div>

        {/* Role Options */}
        <div className="space-y-3">
          {/* Volunteer Role */}
          <button
            onClick={() => onSelectRole(Role.VOLUNTEER)}
            className="w-full p-4 bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 border-2 border-blue-200 dark:border-blue-700 rounded-xl hover:border-blue-400 dark:hover:border-blue-500 transition-all group"
          >
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 bg-blue-500 rounded-xl flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                <Heart className="w-6 h-6 text-white" />
              </div>
              <div className="text-left flex-1">
                <h3 className="font-semibold text-gray-900 dark:text-white mb-1">
                  Volunteer
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  Help coordinate relief efforts and report disasters in your area
                </p>
              </div>
            </div>
          </button>

          {/* User Role */}
          <button
            onClick={() => onSelectRole(Role.USER)}
            className="w-full p-4 bg-gradient-to-r from-green-50 to-teal-50 dark:from-green-900/20 dark:to-teal-900/20 border-2 border-green-200 dark:border-green-700 rounded-xl hover:border-green-400 dark:hover:border-green-500 transition-all group"
          >
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 bg-green-500 rounded-xl flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                <Users className="w-6 h-6 text-white" />
              </div>
              <div className="text-left flex-1">
                <h3 className="font-semibold text-gray-900 dark:text-white mb-1">
                  User
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  Access emergency resources and stay informed about relief operations
                </p>
              </div>
            </div>
          </button>
        </div>

        {/* Footer Note */}
        <p className="text-xs text-gray-500 dark:text-gray-400 text-center mt-6">
          You can change your role later in account settings
        </p>
      </div>
    </div>
  )
}
