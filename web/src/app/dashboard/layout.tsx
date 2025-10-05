'use client'

import Navigation from "@/components/Navigation"
import { AuthContext } from "@/context/AuthContext"
import { ReactNode, useContext } from "react"

const DashboardLayout = ({ children }: { children: ReactNode }) => {
  const { user, logout } = useContext(AuthContext)

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">  
      <Navigation user={user} onLogout={logout} />
      <main className="pt-4">
        {children}
      </main>
    </div>
  )
}

export default DashboardLayout